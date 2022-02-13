package utils

import (
	"bufio"
	"bytes"
	"context"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net"
	"time"
)

func CopyReader(r io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
	var buf = *bytes.NewBufferString("")
	if r == nil {
		return ioutil.NopCloser(bytes.NewBuffer(nil)), ioutil.NopCloser(bytes.NewBuffer(nil)), Errorf("empty input reader")
	}

	if _, err := buf.ReadFrom(r); err != nil {
		return ioutil.NopCloser(bytes.NewBuffer(nil)), r, err
	}

	if err := r.Close(); err != nil {
		return ioutil.NopCloser(bytes.NewBuffer(nil)), r, err
	}

	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func ReaderToReaderCloser(body io.Reader) io.ReadCloser {
	if body == nil {
		return nil
	}

	rc, ok := body.(io.ReadCloser)
	if !ok {
		rc = ioutil.NopCloser(body)
	}
	return rc
}

func ReadWithChunkLen(raw []byte, length int) chan []byte {
	outC := make(chan []byte)
	go func() {
		defer close(outC)

		scanner := bufio.NewScanner(bytes.NewBuffer(raw))
		scanner.Split(bufio.ScanBytes)

		buffer := []byte{}
		n := 0
		for scanner.Scan() {
			buff := scanner.Bytes()
			buffSize := len(buff)
			n += buffSize
			buffer = append(buffer, buff...)

			if n >= length {
				outC <- buffer
				n = 0
				buffer = []byte{}
			}
		}

		if len(buffer) > 0 {
			outC <- buffer
		}
	}()
	return outC
}

func ReadWithContext(ctx context.Context, reader io.Reader) []byte {
	outc := make(chan []byte)
	go func() {
		defer close(outc)

		scanner := bufio.NewScanner(reader)
		scanner.Split(bufio.ScanBytes)
		for scanner.Scan() {

			if ctx.Err() != nil {
				return
			}

			outc <- scanner.Bytes()
		}
	}()

	var raw []byte
	for {
		select {
		case data, ok := <-outc:
			if !ok {
				return raw
			}
			raw = append(raw, data...)
		case <-ctx.Done():
			return raw
		}
	}
}

func ReadConnWithTimeout(r net.Conn, timeout time.Duration) ([]byte, error) {
	err := r.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return nil, errors.Errorf("set read timeout failed: %s", err)
	}

	raw, err := ioutil.ReadAll(ioutil.NopCloser(r))
	if len(raw) > 0 {
		return raw, nil
	}

	return nil, errors.Errorf("read empty: %s", err)
}

func WriteConnWithTimeout(w net.Conn, timeout time.Duration, data []byte) error {
	err := w.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		return errors.Errorf("write failed: %s", err)
	}

	_, err = w.Write(data)
	if err != nil {
		return errors.Errorf("write failed: %s", err)
	}

	return nil
}

func ConnExpect(c net.Conn, timeout time.Duration, callback func([]byte) bool) (bool, error) {
	err := c.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return false, errors.Errorf("set timeout for reading conn failed: %s", err)
	}

	scanner := bufio.NewScanner(c)
	scanner.Split(bufio.ScanBytes)

	var buf []byte
	for scanner.Scan() {
		buf = append(buf, scanner.Bytes()...)
		if callback(buf) {
			return true, nil
		}
	}
	return false, nil
}
