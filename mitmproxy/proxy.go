package mitmproxy

import (
	"bufio"
	"context"
	"crypto/tls"
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
	reader := bufio.NewReader(conn)
	httpRequest, err := http.ReadRequest(reader)
	if err != nil {
		log.Errorf("read request from[%v] failed: %s", conn.RemoteAddr(), err)
		return
	}

	packet, err := utils.HttpDumpWithBody(httpRequest, true)
	if err != nil {
		log.Errorf("dump packet error: %s", err)
		return
	}

	// 如果是 CONNECT 给人家回一个 Established，这种一般是开隧道用的
	if httpRequest.Method == "CONNECT" {
		_, err = conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
		if err != nil {
			log.Errorf("write CONNECT response to %v failed: %s", conn.RemoteAddr(), err)
			return
		}
		println(string(packet))
		log.Info("CONNECT response is sent back")
		m.connect(httpRequest, conn)
		return
	} else {
		// 如果是带了一个 Proxy-Connection 的话，则说明，这也是一个代理的请求类型。
		switch strings.ToLower(httpRequest.Header.Get("Proxy-Connection")) {
		case "":
			println(string(packet))
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
			httpRequest.RequestURI = httpRequest.URL.RequestURI()
			rawPacket, err := utils.HttpDumpWithBody(httpRequest, true)
			if err != nil {
				log.Errorf("fetch origin request failed: %s", err)
				return
			}
			_, err = newConn.Write(rawPacket)
			if err != nil {
				log.Errorf("write packet to %v failed: %s", err)
				return
			}
			rspReader := bufio.NewReader(newConn)
			rsp, err := http.ReadResponse(rspReader, httpRequest)
			if err != nil {
				log.Errorf("read response failed: %s", err)
				return
			}
			raw, err := utils.HttpDumpWithBody(rsp, true)
			if err != nil {
				log.Errorf("dump request failed: %s", err)
				return
			}
			_, err = conn.Write(raw)
			if err != nil {
				log.Errorf("write feedback to [%v] failed: %s", conn.RemoteAddr(), err)
				return
			}

			if !keepAlive {
				return
			}

			for {
				log.Infof("start read request from: %s", conn.RemoteAddr())
				request, err := http.ReadRequest(reader)
				if err != nil {
					log.Errorf("read request failed: %s", err)
					return
				}

				raw, err := utils.HttpDumpWithBody(request, true)
				if err != nil {
					return
				}
				newConn.Write(raw)
				response, err := http.ReadResponse(rspReader, request)
				if err != nil {
					log.Errorf("read response failed: %s", err)
					return
				}

				raw, err = utils.HttpDumpWithBody(response, true)
				if err != nil {
					return
				}

				conn.Write(raw)
			}
		}
	}
}

func (m *MITMProxy) connect(httpRequest *http.Request, conn net.Conn) {
	host := httpRequest.URL.Host
	log.Infof("%v CONNECTed %v start to peek first byte to identify https/tls", conn.RemoteAddr(), host)

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
	log.Infof("fetched first byte[0x%2x]", raw[0])

	isHttps := utils.NewAtomicBool()
	var originHttpConn net.Conn
	var sni string
	switch raw[0] {
	case 0x16:
		// HTTPS 升级，这是核心步骤
		log.Infof("upgrade %v to tls(https)", conn.RemoteAddr())
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

	if isHttps.IsSet() {
		log.Infof("MITM Hijacked HTTPS SNI: %v", sni)
	}
	_ = originHttpConn

	inCome := bufio.NewReader(originHttpConn)
	newConn, err := m.newConnFor(host, isHttps.IsSet())
	if err != nil {
		log.Errorf("create new conn to %v failed: %s", host, err)
		return
	}
	outConn := bufio.NewReader(newConn)
	first := true
	for {
		req, err := http.ReadRequest(inCome)
		if err != nil {
			log.Errorf("read request from client [%v] failed: %s", conn.RemoteAddr(), err)
			return
		}
		packet, err := utils.HttpDumpWithBody(req, true)
		if err != nil {
			return
		}

		newConn.Write(packet)
		rsp, err := http.ReadResponse(outConn, req)
		if err != nil {
			return
		}
		responsePacket, err := utils.HttpDumpWithBody(rsp, true)
		if err != nil {
			return
		}
		originHttpConn.Write(responsePacket)
		if first {
			first = false
			if connKeepalive {
				log.Errorf("close by connection for %v => %v", conn.RemoteAddr(), newConn.RemoteAddr())
				return
			}
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
