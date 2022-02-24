package mitmproxy

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"online/common/utils"
)

func mirrorRequest(reqInConn net.Conn, reqOut net.Conn, cbs ...func(r *http.Request)) error {
	reqIn := bufio.NewReader(io.TeeReader(reqInConn, reqOut))
	request, err := http.ReadRequest(reqIn)
	if err != nil && err != io.EOF {
		return utils.Errorf("read request-in failed: %s", err)
	}
	if request == nil {
		return nil
	}
	if request.Body != nil {
		raw, _ := ioutil.ReadAll(request.Body)
		request.GetBody = func() (io.ReadCloser, error) {
			return ioutil.NopCloser(bytes.NewBuffer(raw)), nil
		}
	}

	for _, cb := range cbs {
		cb(request)
	}
	return nil
}

func mirrorResponse(request *http.Request, rspInConn net.Conn, rspOut net.Conn, cbs ...func(r *http.Response)) error {
	rspIn := bufio.NewReader(io.TeeReader(rspInConn, rspOut))
	response, err := http.ReadResponse(rspIn, request)
	if err != nil {
		return utils.Errorf("read request-in failed: %s", err)
	}
	if response.Body != nil {
		raw, _ := ioutil.ReadAll(response.Body)
		response.Body = ioutil.NopCloser(bytes.NewBuffer(raw))
	}

	for _, cb := range cbs {
		cb(response)
	}
	return nil
}
