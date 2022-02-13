package utils

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/denisbrodbeck/machineid"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"online/common/log"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func IContains(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}

func Errorf(origin string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(origin, args...))
}

func Error(i interface{}) error {
	return errors.New(fmt.Sprint(i))
}

func StringLowerAndTrimSpace(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func TrimFileNameExt(raw string) string {
	e := filepath.Ext(raw)
	if e == "" {
		return raw
	}

	return strings.Trim(strings.TrimSuffix(raw, e), ". ")
}

func IntArrayContains(array []int, element int) bool {
	for _, s := range array {
		if element == s {
			return true
		}
	}
	return false
}

// BKDR Hash Function
func BKDRHash(str []byte) uint32 {
	var seed uint32 = 131 // 31 131 1313 13131 131313 etc..
	var hash uint32 = 0
	for i := 0; i < len(str); i++ {
		hash = hash*seed + uint32(str[i])
	}

	return hash
}

func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func TimeoutContext(d time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), d)
	return ctx
}

func FloatSecondDuration(f float64) time.Duration {
	return time.Duration(float64(time.Second) * f)
}

func StringAsFileParams(target interface{}) []byte {
	switch ret := target.(type) {
	case string:
		if GetFirstExistedPath(ret) != "" {
			raw, err := ioutil.ReadFile(ret)
			if err != nil {
				return []byte(ret)
			}
			return raw
		} else {
			return []byte(ret)
		}
	case []string:
		return []byte(strings.Join(ret, "\n"))
	case []byte:
		return ret
	case io.Reader:
		raw, err := ioutil.ReadAll(ret)
		if err != nil {
			return nil
		}
		return raw
	default:
		log.Errorf("cannot covnert %v to file content", spew.Sdump(target))
		return nil
	}
}

func TimeoutContextSeconds(d float64) context.Context {
	return TimeoutContext(FloatSecondDuration(d))
}

func interfaceToBytes(i interface{}) []byte {
	var bytes []byte

	switch ret := i.(type) {
	case string:
		bytes = []byte(ret)
	case []byte:
		bytes = ret
	case io.Reader:
		bytes, _ = ioutil.ReadAll(ret)
	default:
		bytes = []byte(fmt.Sprint(i))
	}

	return bytes
}

func EncodeToHex(i interface{}) string {
	raw := interfaceToBytes(i)
	return hex.EncodeToString(raw)
}

func DecodeHex(i string) ([]byte, error) {
	return hex.DecodeString(i)
}

func GetMachineCode() (string, error) {
	switch runtime.GOOS {
	case "linux":
		raw, err := exec.Command("cat", "/sys/class/dmi/id/product_uuid").CombinedOutput()
		if err != nil {
			return "", Errorf("fetch system machine code failed: %s", err)
		}
		return EncodeToHex(raw), nil
	}

	id, _ := machineid.ID()
	if id == "" {
		return "", Errorf("fetch machine code failed")
	}
	//if err != nil {
	//	m := fmt.Sprintf("get machine id failed: %s", err)
	//	return "", errors.New(m)
	//}
	//
	//if id == "" {
	//	return "", Errorf("empty machine-id...")
	//}
	return id, nil
}

func FixJsonRawBytes(rawBytes []byte) []byte {
	rawBytes = []byte(EscapeInvalidUTF8Byte(rawBytes))
	rawBytes = bytes.ReplaceAll(rawBytes, []byte("\\u0000"), []byte(" "))
	return rawBytes
}

func Jsonify(i interface{}) []byte {
	raw, err := json.Marshal(i)
	if err != nil {
		return []byte("{}")
	}
	return raw
}

func NewDefaultHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: nil,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			MaxConnsPerHost:    50,
			DisableCompression: true,
			DisableKeepAlives:  true,
		},
		Timeout: 15 * time.Second,
	}
}

func NewDefaultTLSClient(conn net.Conn) *tls.Conn {
	return tls.Client(conn, NewDefaultTLSConfig())
}

func NewDefaultTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
	}
}

func FixHTTPRequestForHTTPDo(r *http.Request) (*http.Request, error) {
	return FixHTTPRequestForHTTPDoWithHttps(r, false)
}

func FixHTTPRequestForHTTPDoWithHttps(r *http.Request, isHttps bool) (*http.Request, error) {
	var bodyRaw []byte
	var err error
	if r.Body != nil {
		bodyRaw, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, Errorf("read body failed: %s", err)
		}
	}

	if r.URL.Scheme == "" {
		if isHttps {
			r.URL.Scheme = "https"
		} else {
			r.URL.Scheme = "http"
		}
	}

	req, err := http.NewRequest(r.Method, r.URL.String(), bytes.NewBuffer(bodyRaw))
	if err != nil {
		return nil, Errorf("build http.Request[%v] failed: %v", r.URL.String(), err)
	}

	for key, values := range r.Header {
		if len(values) == 1 {
			req.Header.Set(key, values[0])
		} else {
			for _, v := range values {
				req.Header.Add(key, v)
			}
		}
	}

	return req, nil
}
