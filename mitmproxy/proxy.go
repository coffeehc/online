package mitmproxy

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"online/common/log"
	"online/common/utils"
	"strings"
	"time"
)

type MITMProxy struct {
	config *Config

	dialer *net.Dialer
}

func NewMITMProxy(opt ...Option) (*MITMProxy, error) {
	config, err := NewConfig(opt...)
	if err != nil {
		return nil, utils.Errorf("generate config failed: %s", err)
	}
	proxy := &MITMProxy{config: config}
	if config.Timeout <= 0 {
		config.Timeout = 15 * time.Second
	}
	proxy.dialer = &net.Dialer{
		Timeout: config.Timeout,
	}
	proxy.config.mitmConfig.SkipTLSVerify(true)
	return proxy, nil
}

func (m *MITMProxy) Run(ctx context.Context) error {
	addr := utils.HostPort(m.config.Host, m.config.Port)
	log.Infof("start to serve on proxy http://%v", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return utils.Errorf("create tcp listener [tcp://%v] failed: %v", addr, err)
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}
		log.Infof("accept from %v => %v", conn.RemoteAddr(), conn.LocalAddr())
		go m.serve(conn)
	}
}

func (m *MITMProxy) serve(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("panic from conn: %v", conn.RemoteAddr())
			return
		}
	}()

	defer conn.Close()

	// 透明模式和非透明模式处理方式不一样
	// 透明模式先不管了
	if m.config.TransparentMode {
		panic("transparent mode failed: not implemented")
	}

	// 下面是代理模式的实现过程
	// 代理首先会有一个 "CONNECT" 过来，所以需要针对这个进行处理，一般回复一个 HTTP/1.1 200 Established\r\n\r\n
	var firstRequestMirrorBytes bytes.Buffer
	httpRequest, err := http.ReadRequest(bufio.NewReader(io.TeeReader(conn, &firstRequestMirrorBytes)))
	if err != nil {
		log.Errorf("read request from[%v] failed: %s", conn.RemoteAddr(), err)
		return
	}

	if httpRequest.Method == "CONNECT" {
		// 如果是 CONNECT 给人家回一个 Established，这种一般是开隧道用的
		// 开隧道的情况一般用在 HTTPS 或者其他隧道上，兼容 HTTP 隧道
		_, err = conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
		if err != nil {
			log.Errorf("write CONNECT response to %v failed: %s", conn.RemoteAddr(), err)
			return
		}
		m.connect(httpRequest, conn)
		return
	} else {
		// 如果是带了一个 Proxy-Connection 的话，则说明，这也是一个代理的请求类型。
		switch strings.ToLower(httpRequest.Header.Get("Proxy-Connection")) {
		case "":
			log.Errorf("not a proxy connection from %v", conn.RemoteAddr())
			return
		default:
			var keepAlive = false
			if httpRequest.Header.Get("Proxy-Connection") == "keep-alive" {
				keepAlive = true
			}
			httpRequest.Header.Del("Proxy-Connection")
			host := httpRequest.URL.Host
			newConn, err := m.newConnFor(host, false)
			if err != nil {
				log.Errorf("create new conn to %v failed: %s", host, err)
				return
			}
			defer newConn.Close()
			httpRequest.RequestURI = httpRequest.URL.RequestURI()

			headers := firstRequestMirrorBytes.Bytes()
			headers = bytes.ReplaceAll(
				headers,
				[]byte(fmt.Sprintf("Proxy-Connection: %v\r\n", httpRequest.Header.Get("Proxy-Connection"))),
				nil,
			)

			var body []byte
			if httpRequest.Body != nil {
				body, _ = ioutil.ReadAll(httpRequest.Body)
			}

			newConn.Write(headers)
			newConn.Write(body)
			err = mirrorResponse(httpRequest, newConn, conn)
			if err != nil {
				log.Errorf("mirror response failed: %s", err)
				return
			}

			if !keepAlive {
				return
			}
			m.handleHTTP(conn, newConn, keepAlive)
		}
	}
}

func (m *MITMProxy) connect(httpRequest *http.Request, conn net.Conn) {
	host := httpRequest.URL.Host
	//log.Infof("%v CONNECTed %v start to peek first byte to identify https/tls", conn.RemoteAddr(), host)
	connKeepalive := false
	if httpRequest.Header.Get("Proxy-Connection") == "keep-alive" {
		connKeepalive = true
	}
	httpRequest.Header.Del("Proxy-Connection")
	_ = connKeepalive

	originConnPeekable := utils.NewPeekableNetConn(conn)
	raw, err := originConnPeekable.Peek(1)
	if err != nil {
		log.Errorf("peek [%v] failed: %s", conn.RemoteAddr(), err)
		return
	}
	isHttps := utils.NewAtomicBool()
	var originHttpConn net.Conn
	var sni string
	switch raw[0] {
	case 0x16:
		// HTTPS 升级，这是核心步骤
		//log.Infof("upgrade/hijacked %v to tls(https)", conn.RemoteAddr())
		tconn := tls.Server(originConnPeekable, m.config.mitmConfig.TLS())
		err := tconn.Handshake()
		if err != nil {
			log.Errorf("tls handshake failed: %s", err)
			println(string(m.config.Ca))
			println("----------------------")
			println(string(m.config.Key))
			return
		}
		originHttpConn = tconn
		sni = tconn.ConnectionState().ServerName
		isHttps.Set()
	default:
		// HTTP
		log.Infof("recognized %v as http", conn.RemoteAddr())
		originHttpConn = originConnPeekable
		isHttps.UnSet()
	}

	_ = sni

	newConn, err := m.newConnFor(host, isHttps.IsSet())
	if err != nil {
		log.Errorf("create new conn to %v failed: %s", host, err)
		return
	}
	defer newConn.Close()
	m.handleHTTP(originHttpConn, newConn, connKeepalive)
}

func (m *MITMProxy) handleHTTP(in net.Conn, out net.Conn, keepalive bool) {
	for {
		var req *http.Request
		var rsp *http.Response
		err := mirrorRequest(in, out, func(r *http.Request) {
			req = r
		})
		if err != nil {
			log.Errorf("mirror request failed: %s", err)
			return
		}

		err = mirrorResponse(req, out, in, func(r *http.Response) {
			rsp = r
		})
		if err != nil {
			log.Errorf("mirror response failed: %s", err)
			return
		}

		if !keepalive {
			log.Errorf("close by connection for %v <=> %v", in.RemoteAddr(), in.RemoteAddr())
			return
		}
	}
}

func (m *MITMProxy) newConnFor(target string, isTls bool) (net.Conn, error) {
	host, port, err := utils.ParseStringToHostPort(target)
	if err != nil {
		host = target
		port = 80
	}
	if !utils.IsIPv4(host) {
		host = utils.GetFirstIPByDnsWithCache(host, 10*time.Second)
		if host == "" {
			return nil, utils.Errorf("dns error for %v", host)
		}
	}

	if !isTls {
		conn, err := m.dialer.Dial("tcp", utils.HostPort(host, port))
		if err != nil {
			return nil, utils.Errorf("dial to new connection %v failed: %s", utils.HostPort(host, port))
		}
		return conn, nil
	} else {
		conn, err := tls.DialWithDialer(m.dialer, "tcp", utils.HostPort(host, port), &tls.Config{
			InsecureSkipVerify: true,
		})
		if err != nil {
			return nil, utils.Errorf("dial tls to conn %v failed: %s", utils.HostPort(host, port), err)
		}
		return conn, nil
	}
}
