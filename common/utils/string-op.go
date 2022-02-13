package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func PrettifyListFromStringSplited(Raw string, sep string) (targets []string) {
	targetsRaw := strings.Split(Raw, sep)
	for _, tRaw := range targetsRaw {
		r := strings.TrimSpace(tRaw)
		if len(r) > 0 {
			targets = append(targets, r)
		}
	}
	return
}

func ToLowerAndStrip(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func StringSliceContain(s []string, raw string) bool {
	for _, k := range s {
		if k == raw {
			return true
		}
	}
	return false
}

func StringContainsAnyOfSubString(s string, subs []string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func IStringContainsAnyOfSubString(s string, subs []string) bool {
	for _, sub := range subs {
		if IContains(s, sub) {
			return true
		}
	}
	return false
}

func ConvertToStringSlice(raw ...interface{}) (r []string) {
	for _, e := range raw {
		r = append(r, fmt.Sprintf("%v", e))
	}
	return
}

func ChanStringToSlice(c chan string) (result []string) {
	for l := range c {
		result = append(result, l)
	}
	return
}

var (
	cStyleCharPRegexp, _ = regexp.Compile(`\\((x[0-9abcdef]{2})|([0-9]{1,3}))`)
)

func ParseCStyleBinaryRawToBytes(raw []byte) []byte {
	// like "\\x12" => "\x12"
	return cStyleCharPRegexp.ReplaceAllFunc(raw, func(i []byte) []byte {
		if bytes.HasPrefix(i, []byte("\\x")) {
			if len(i) == 4 {
				rawChar := string(i[2:])
				charInt, err := strconv.ParseInt("0x"+string(rawChar), 0, 16)
				if err != nil {
					return i
				}
				return []byte{byte(charInt)}
			} else {
				return i
			}
		} else if bytes.HasPrefix(raw, []byte("\\")) {
			if len(i) > 1 && len(i) <= 4 {
				rawChar := string(i[1:])
				charInt, err := strconv.ParseInt(string(rawChar), 10, 8)
				if err != nil {
					return i
				}
				return []byte{byte(charInt)}
			} else {
				return i
			}
		}
		return i
	})
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func EscapeInvalidUTF8Byte(s []byte) string {
	// 这个操作返回的结果和原始字符串是非等价的
	ret := make([]rune, 0, len(s)+20)
	start := 0
	for {
		r, size := utf8.DecodeRune(s[start:])
		if r == utf8.RuneError {
			// 说明是空的
			if size == 0 {
				break
			} else {
				// 不是 rune
				ret = append(ret, []rune(fmt.Sprintf("\\x%02x", s[start]))...)
			}
		} else {
			// 不是换行之类的控制字符
			if unicode.IsControl(r) && !unicode.IsSpace(r) {
				ret = append(ret, []rune(fmt.Sprintf("\\x%02x", r))...)
			} else {
				// 正常字符
				ret = append(ret, r)
			}
		}
		start += size
	}
	return string(ret)
}

func GBKSafeString(s []byte) (string, error) {
	if utf8.Valid(s) {
		return string(s), nil
	}

	raw, err := GbkToUtf8(s)
	if err != nil {
		return "", errors.Errorf("failed to parse gbk: %s", err)
	}

	if utf8.Valid(raw) {
		return string(raw), nil
	}

	return "", errors.Errorf("invalid utf8: %#v", raw)

}

func LastLine(s []byte) []byte {
	s = bytes.TrimSpace(s)
	scanner := bufio.NewScanner(bytes.NewReader(s))
	scanner.Split(bufio.ScanLines)

	var lastLine = s
	for scanner.Scan() {
		lastLine = scanner.Bytes()
	}

	return lastLine
}

func RemoveUnprintableChars(raw string) string {
	scanner := bufio.NewScanner(bytes.NewBufferString(raw))
	scanner.Split(bufio.ScanBytes)

	var r []byte
	for scanner.Scan() {
		c := scanner.Bytes()[0]

		if c <= 0x7e && c >= 0x20 {
			r = append(r, c)
		}
	}

	return string(r)
}

func RemoveRepeatedWithStringSlice(slice []string) []string {
	r := map[string]interface{}{}
	for _, s := range slice {
		r[s] = 1
	}

	var r2 []string
	for k, _ := range r {
		r2 = append(r2, k)
	}
	return r2
}

var (
	titleRegexp = regexp.MustCompile(`(?is)\<title\>(.*?)\</?title\>`)
)

func ExtractTitleFromHTMLTitle(s string, defaultValue string) string {
	var title string
	l := titleRegexp.FindString(s)
	if len(l) > 15 {
		title = EscapeInvalidUTF8Byte([]byte(l))[7 : len(l)-8]
	}
	titleRunes := []rune(title)
	if len(titleRunes) > 20 {
		title = string(titleRunes[0:17]) + "..."
	}

	if title == "" {
		return defaultValue
	}

	return title
}
