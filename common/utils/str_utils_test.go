package utils

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"online/common/log"
	"testing"
	"time"
)

func TestRemoveUnprintableChars(t *testing.T) {
	cases := map[string]string{
		"\x00W\xffO\x00R\x00K": "WORK",
	}
	for input, output := range cases {
		if result := RemoveUnprintableChars(input); result == output {
			continue
		} else {
			t.Logf("expect %#v got %#v", output, result)
			t.FailNow()
		}
	}
}

func TestParseStringToHostPort(t *testing.T) {
	type Result struct {
		Host string
		Port int
	}
	cases := map[string]Result{
		"http://baidu.com":     {Host: "baidu.com", Port: 80},
		"https://baidu.com":    {Host: "baidu.com", Port: 443},
		"https://baidu.com:88": {Host: "baidu.com", Port: 88},
		"http://baidu.com:88":  {Host: "baidu.com", Port: 88},
		"1.2.3.4:1":            {Host: "1.2.3.4", Port: 1},
		"baidu.com:1":          {Host: "baidu.com", Port: 1},
	}

	falseCases := []string{
		"baidu.com", "1.2.3.5", "[1:123:123:123]",
	}

	for raw, result := range cases {
		host, port, err := ParseStringToHostPort(raw)
		if err != nil {
			t.Errorf("parse %s failed: %s", raw, err)
			t.FailNow()
		}

		if result.Host == host && result.Port == port {
			continue
		} else {
			t.Errorf("parse result failed: %s expect: %s:%v actually: %s:%v", raw, result.Host, result.Port,
				host, port)
			t.FailNow()
		}
	}

	for _, c := range falseCases {
		_, _, err := ParseStringToHostPort(c)
		if err != nil {

		} else {
			t.Errorf("%s should failed now", c)
			t.FailNow()
		}
	}
}

func TestSliceGroup(t *testing.T) {
	s := SliceGroup([]string{
		"1", "1", "1",
		"1", "1", "1",
		"1", "1", "1",
		"1", "1", "1",
		"1", "1", "1",
		"1", "1", "1",
		"1", "1", "1",
	}, 3)
	log.Info(spew.Sdump(s))
	assert.True(t, len(s) == 7, "%v", spew.Sdump(s))
}

func TestGetFirstIPFromHostWithTimeout(t *testing.T) {
	ip := GetFirstIPFromHostWithTimeout(5 * time.Second, "baidu.com", nil)
	spew.Dump(ip)
}

func TestParseStringToCClassHosts(t *testing.T) {
	spew.Dump(ParseStringToCClassHosts("192.168.1.2,baidu.com,192.168.1.22,www.uestc.edu.cn"))
}