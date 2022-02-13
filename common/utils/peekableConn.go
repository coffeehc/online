package utils

import (
	"bytes"
	"io"
	"net"
)

type PeekableNetConn struct {
	net.Conn

	reader io.Reader
}

func NewPeekableNetConn(conn net.Conn) *PeekableNetConn {
	reader, writer := io.Pipe()
	go func() {
		_, _ = io.Copy(writer, conn)
		_ = writer.Close()
	}()

	pc := &PeekableNetConn{
		Conn:   conn,
		reader: reader,
	}
	return pc
}

func (p *PeekableNetConn) Peek(i int) ([]byte, error) {
	var buf = make([]byte, i)
	_, err := p.Conn.Read(buf)
	if err != nil {
		return nil, err
	}

	p.reader = io.MultiReader(bytes.NewBuffer(buf), p.reader)
	return buf, nil
}

func (p *PeekableNetConn) Read(b []byte) (int, error) {
	return p.reader.Read(b)
}
