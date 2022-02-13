package utils

import (
	"online/common/log"
	"testing"
)

func TestLastLine(t *testing.T) {
	trueCase := map[string]string{
		`aasdfas
asdf
asdf
asdf
aaaa`: "aaaa",
		`aaaa`: "aaaa",
	}

	for k, v := range trueCase {
		if v != string(LastLine([]byte(k))) {
			t.FailNow()
		}
		log.Infof("%s 's last line is %s", k, v)
	}
}
