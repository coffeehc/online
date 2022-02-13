package utils

import (
	"bufio"
	"bytes"
	"github.com/pkg/errors"
	"os"
)

func FileLineReader(file string) (chan []byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Errorf("failed to read file: %s", err)
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	outC := make(chan []byte)
	go func() {
		defer close(outC)

		for scanner.Scan() {
			outC <- bytes.TrimSpace(scanner.Bytes())
		}
	}()

	return outC, nil
}
