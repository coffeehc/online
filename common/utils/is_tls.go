package utils

import (
	"crypto/tls"
	"net"
	"time"
)

func IsTLSService(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err == nil {
		defer conn.Close()
		conn := tls.Client(conn, &tls.Config{InsecureSkipVerify: true, MinVersion: tls.VersionSSL30})
		err = conn.Handshake()
		return err == nil
	}
	return false
}
