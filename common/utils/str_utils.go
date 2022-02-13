package utils

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gobwas/glob"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"online/common/log"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

import (
	"github.com/ReneKroon/ttlcache"
)

func ExtractStrContextByKeyword(raw []byte, res []string) []string {
	rawStrContent := string(raw)
	var details []string
	for _, keyword := range res {
		if index := strings.Index(rawStrContent, keyword); index > 0 {
			info := ""

			end := index + len(keyword) + 512

			if index <= 512 {
				info += rawStrContent[:index]
			} else {
				info += rawStrContent[index-512 : index+len(keyword)]
			}

			if end >= len(rawStrContent) {
				info += rawStrContent[index:]
			} else {
				info += rawStrContent[index:end]
			}

			details = RemoveRepeatStringSlice(append(details, EscapeInvalidUTF8Byte([]byte(info))))
		}
	}
	return details
}

// SliceContainsSlick
func StringSliceContainsAll(o []string, elements ...string) bool {
	for _, e := range elements {
		if !StringArrayContains(o, e) {
			return false
		}
	}
	return true
}

// 元素去重
func RemoveRepeatStringSlice(slc []string) []string {
	if len(slc) < 1024 {
		return RemoveRepeatStringSliceByLoop(slc)
	} else {
		return RemoveRepeatStringSliceByMap(slc)
	}
}

// 元素去重
func RemoveRepeatUintSlice(slc []uint) []uint {
	if len(slc) < 1024 {
		return RemoveRepeatUintSliceByLoop(slc)
	} else {
		return RemoveRepeatUintSliceByMap(slc)
	}
}

func RemoveRepeatStringSliceByLoop(slc []string) []string {
	result := []string{}
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false
				break
			}
		}
		if flag {
			result = append(result, slc[i])
		}
	}
	return result
}

func RemoveRepeatStringSliceByMap(slc []string) []string {
	result := []string{}
	tempMap := map[string]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

func RemoveRepeatUintSliceByLoop(slc []uint) []uint {
	result := []uint{}
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false
				break
			}
		}
		if flag {
			result = append(result, slc[i])
		}
	}
	return result
}

func RemoveRepeatUintSliceByMap(slc []uint) []uint {
	result := []uint{}
	tempMap := map[uint]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

func StringArrayContains(array []string, element string) bool {
	for _, s := range array {
		if element == s {
			return true
		}
	}
	return false
}

func StringHasPrefix(s string, prefix []string) bool {
	for _, x := range prefix {
		if strings.HasPrefix(strings.ToLower(s), strings.ToLower(x)) {
			return true
		}
	}
	return false
}

func StringGlobArrayContains(array []string, element string) bool {
	var globs []glob.Glob
	var strList []string
	for _, r := range array {
		rule, err := glob.Compile(r)
		if err != nil {
			strList = append(strList, r)
			continue
		}
		globs = append(globs, rule)
	}

	for _, s := range array {
		if element == s {
			return true
		}
	}

	for _, g := range globs {
		if g.Match(element) {
			return true
		}
	}
	return false
}

func StringArrayIndex(array []string, element string) int {
	for index, s := range array {
		if element == s {
			return index
		}
	}
	return -1
}

func StringOr(s ...string) string {
	for _, i := range s {
		if i != "" {
			return i
		}
	}
	return ""
}

func IntLargerZeroOr(s ...int) int {
	for _, i := range s {
		if i > 0 {
			return i
		}
	}
	return 0
}

func StringArrayFilterEmpty(array []string) []string {
	var ret []string
	for _, a := range array {
		if a == "" {
			continue
		}
		ret = append(ret, a)
	}
	return ret
}

func StringArrayMerge(t ...[]string) []string {
	m := map[string]interface{}{}
	for _, ta := range t {
		for _, i := range ta {
			m[i] = false
		}
	}

	var n []string
	for k, _ := range m {
		n = append(n, k)
	}
	return n
}

func StringSplitAndStrip(raw string, sep string) []string {
	var l = []string{}

	for _, v := range strings.Split(raw, sep) {
		s := strings.TrimSpace(v)
		if s != "" {
			l = append(l, s)
		}
	}

	return l
}

func BytesStripUntil(raw []byte, s byte) []byte {
	for end := len(raw); end > 0; end-- {
		c_index := end - 1
		if raw[c_index] == s {
			return raw[:c_index]
		}
	}
	return []byte{}
}

func ExtractRawHeadersFromRequest(req *http.Request) string {
	header := []string{}

	header = append(header, fmt.Sprintf("%s %s %s", req.Method, req.URL.Path, req.Proto))

	for headerName, headerValues := range req.Header {
		for _, value := range headerValues {
			header = append(header, fmt.Sprintf("%s: %s", headerName, value))
		}
	}

	header = append(header, "")
	return strings.Join(header, "\r\n")
}

func ExtractRawHeadersFromResponse(rsp *http.Response) string {
	header := []string{}

	header = append(header, fmt.Sprintf("%s %s", rsp.Proto, rsp.Status))

	for headerName, headerValues := range rsp.Header {
		for _, value := range headerValues {
			header = append(header, fmt.Sprintf("%s: %s", headerName, value))
		}
	}

	header = append(header, "")
	return strings.Join(header, "\r\n")
}

func StringToAsciiBytes(s string) []byte {
	t := make([]byte, utf8.RuneCountInString(s))
	i := 0
	for _, r := range s {
		t[i] = byte(r)
		i++
	}
	return t
}

func AsciiBytesToRegexpMatchedRunes(in []byte) []rune {
	result := make([]rune, len(in))
	for i, b := range in {
		result[i] = rune(b)
	}
	return result
}

func AsciiBytesToRegexpMatchedString(in []byte) string {
	return string(AsciiBytesToRegexpMatchedRunes(in))
}

func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	if i := strings.IndexByte(hostport, ']'); i != -1 {
		return strings.TrimPrefix(hostport[:i], "[")
	}
	return hostport[:colon]
}

func portOnly(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return ""
	}
	if i := strings.Index(hostport, "]:"); i != -1 {
		return hostport[i+len("]:"):]
	}
	if strings.Contains(hostport, "]") {
		return ""
	}
	return hostport[colon+len(":"):]
}

var (
	dnsTtlCache = ttlcache.NewCache()
)

func GetFirstIPByDnsWithCache(domain string, timeout time.Duration, dnsServers ...string) string {
	data, ok := dnsTtlCache.Get(domain)
	if ok {
		if raw, _ := data.(string); raw != "" {
			return raw
		}
	}
	r := GetFirstIPFromHostWithTimeout(timeout, domain, dnsServers)
	dnsTtlCache.SetWithTTL(domain, r, 60*time.Second)
	return r
}

func InterfaceToBytes(i interface{}) []byte {
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

func InterfaceToString(i interface{}) string {
	return string(InterfaceToBytes(i))
}

func ParseStringUrlToWebsiteRootPath(s string) string {
	ins, _ := ParseStringUrlToUrlInstance(s)
	if ins == nil {
		return s
	}

	ins.Path = "/"
	ins.RawPath = "/"
	ins.RawQuery = ""
	return ins.String()
}

func ParseStringUrlToUrlInstance(s string) (*url.URL, error) {
	return url.Parse(s)
}

func ParseStringToHostPort(raw string) (host string, port int, err error) {
	if strings.Contains(raw, "://") {
		urlObject, _ := url.Parse(raw)
		if urlObject != nil {
			// 处理 URL
			portRaw := urlObject.Port()
			portInt64, err := strconv.ParseInt(portRaw, 10, 32)
			if err != nil || portInt64 <= 0 {
				switch urlObject.Scheme {
				case "http":
					port = 80
				case "https":
					port = 443
				}
			} else {
				port = int(portInt64)
			}

			host = urlObject.Hostname()
			err = nil
			return host, port, err
		}
	}

	host = stripPort(raw)
	portStr := portOnly(raw)
	if len(portStr) <= 0 {
		return "", 0, errors.Errorf("unknown port for [%s]", raw)
	}

	portStr = strings.TrimSpace(portStr)
	portInt64, err := strconv.ParseInt(portStr, 10, 32)
	if err != nil {
		return "", 0, errors.Errorf("%s parse port(%s) failed: %s", raw, portStr, err)
	}

	port = int(portInt64)
	err = nil
	return
}

func UrlJoin(origin string, paths ...string) (string, error) {
	u, err := url.Parse(origin)
	if err != nil {
		return "", errors.Errorf("origin:[%s] is not a valid url: %s", origin, err)
	}
	u.Path = path.Join(append([]string{u.Path}, paths...)...)
	s := u.String()
	return s, nil
}

func ParseLines(raw string) chan string {
	outC := make(chan string)
	go func() {
		defer close(outC)

		for _, l := range strings.Split(raw, "\n") {
			hl := strings.TrimSpace(l)
			if hl == "" {
				continue
			}
			outC <- hl
		}
	}()
	return outC
}

func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func ByteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

// 每个单词首字母大写
func InitialCapitalizationEachWords(str string) string {
	if len(str) < 1 {
		return ""
	}
	words := strings.Split(str, " ")
	result := []string{}
	for _, w := range words {
		w = strings.ToUpper(w)[:1] + w[1:]
		result = append(result, w)
	}
	return strings.Join(result, " ")
}

func SliceGroup(origin []string, groupSize int) [][]string {
	var (
		result [][]string
	)

	var count int
	var buffer []string
	for _, i := range origin {
		count++
		buffer = append(buffer, i)

		if count >= groupSize {
			count = 0

			result = append(result, buffer)
			buffer = nil
		}
	}

	if len(buffer) > 0 {
		result = append(result, buffer)
	}
	return result
}

func DomainToIP(domain string, timeout time.Duration) []string {
	if net.ParseIP(FixForParseIP(domain)) == nil {
		// 不是 IP 尝试 dns 查询下
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		rets, err := net.DefaultResolver.LookupIPAddr(ctx, domain)
		if err != nil {
			return nil
		}

		var addrs []string
		for _, i := range rets {
			addrs = append(addrs, i.String())
		}
		return addrs
	} else {
		return []string{domain}
	}
}

var (
	DefaultDNSClient = dns.Client{
		Timeout: 5 * time.Second,
	}
	DefaultDNSConn   = dns.Dial
	DefaultDNSServer = []string{
		"223.5.5.5",       // ali
		"119.29.29.29",    // tencent
		"180.76.76.76",    // baidu
		"114.114.114.114", // dianxin
		"1.1.1.1",         // cf
		"8.8.8.8",
	}
)

func ToNsServer(server string) string {
	// 如果 server 只是一个 IP 则需要把端口加上
	ip := net.ParseIP(server)
	if ip != nil {
		server = ip.String() + ":53"
		return server
	}

	// 这里肯定不是 IP/IP6
	// 所以我们检测是否包含端口，如果不包含端口，则添加端口
	if strings.Contains(server, ":") {
		return server
	}

	for strings.HasSuffix(server, ".") {
		server = server[:len(server)-1]
	}
	server += ":53"
	return server
}

func GetIPFromHostWithContextAndDNSServers(
	ctx context.Context, domain string, DNSServers []string, cb func(domain string) bool,
) error {
	ctx, _ = context.WithCancel(ctx)

	nativeDNS := func() bool {
		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
			}
		}()

		ips, err := net.DefaultResolver.LookupIPAddr(ctx, domain)
		if err != nil {
			log.Error("default resolver parse failed: %s", err)
			return false
		}
		for _, i := range ips {
			if cb(i.String()) {
				continue
			} else {
				return true
			}
		}
		if len(ips) <= 0 {
			return false
		}
		return true
	}

	if nativeDNS() {
		return nil
	}

	if cb == nil {
		return Errorf("callback cannot be empty")
	}

	if DNSServers == nil {
		DNSServers = DefaultDNSServer
	}

	var (
		errs []error
	)

	haveResult := NewBool(false)
	callback := func(domain string) bool {
		haveResult.Set()
		return cb(domain)
	}

	if !strings.HasSuffix(domain, ".") {
		domain = domain + "."
	}

Main:
	for _, qType := range []uint16{dns.TypeA, dns.TypeAAAA} {
		var typeRaw = "A/AAAA"
		switch qType {
		case dns.TypeAAAA:
			typeRaw = "AAAA"
		case dns.TypeA:
			typeRaw = "A"
		}
		for _, server := range DNSServers {
			server = ToNsServer(server)
			msg := &dns.Msg{}
			msg.SetQuestion(domain, qType)
			rsp, rtt, err := DefaultDNSClient.ExchangeContext(ctx, msg, server)
			if err != nil {
				errs = append(errs, Errorf("DNS(%v)[%v] Err: %v", server, typeRaw, err))
				continue
			}

			_ = rtt
			for _, ans := range rsp.Answer {
				switch record := ans.(type) {
				case *dns.A:
					if record.A.String() != "" {
						if callback(record.A.String()) {
							continue
						} else {
							break Main
						}
					}
				case *dns.AAAA:
					if record.AAAA.String() != "" {
						if callback(record.AAAA.String()) {
							continue
						} else {
							break Main
						}
					}
				}
			}
		}
	}

	if !haveResult.IsSet() {
		return Errorf("empty results: QueryErrors: \n"+
			"%v"+
			"\n", spew.Sdump(errs))
	}

	return nil
}

func GetFirstIPFromHostWithContextE(ctx context.Context, domain string, DNSServers []string) (string, error) {
	var Result string
	err := GetIPFromHostWithContextAndDNSServers(ctx, domain, DNSServers, func(domain string) bool {
		Result = domain
		return false
	})
	return Result, err
}

func GetIPsFromHostWithContextE(ctx context.Context, domain string, DNSServers []string) ([]string, error) {
	var results []string
	err := GetIPFromHostWithContextAndDNSServers(ctx, domain, DNSServers, func(domain string) bool {
		results = append(results, domain)
		return true
	})
	return results, err
}

func GetIPsFromHostWithTimeoutE(timeout time.Duration, domain string, dnsServers []string) ([]string, error) {
	return GetIPsFromHostWithContextE(TimeoutContext(timeout), domain, dnsServers)
}

func GetFirstIPFromHostWithTimeoutE(timeout time.Duration, domain string, dnsServres []string) (string, error) {
	return GetFirstIPFromHostWithContextE(TimeoutContext(timeout), domain, dnsServres)
}

func GetFirstIPFromHostWithTimeout(timeout time.Duration, domain string, dnsServres []string) string {
	s, _ := GetFirstIPFromHostWithTimeoutE(timeout, domain, dnsServres)
	return s
}

func GetIPsFromHostWithTimeout(timeout time.Duration, domain string, dnsServers []string) []string {
	r, _ := GetIPsFromHostWithTimeoutE(timeout, domain, dnsServers)
	return r
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func RandNumberStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = numberChar[rand.Intn(len(numberChar))]
	}
	return string(b)
}

const (
	passwordSepcialChars = ",.<>?;:[]{}~!@#$%^&*()_+-="
	AllSepcialChars      = ",./<>?;':\"[]{}`~!@#$%^&*()_+-=\\|"
	littleChar           = "abcdefghijklmnopqrstuvwxyz"
	bigChar              = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChar           = "1234567890"
)

var (
	passwordBase = passwordSepcialChars + littleChar + bigChar + numberChar
)

func IsStrongPassword(s string) bool {
	if len(s) <= 8 {
		return false
	}

	var haveSpecial, haveLittleChar, haveBigChar, haveNumber bool
	for _, c := range s {
		ch := string(c)
		if strings.Contains(passwordSepcialChars, ch) {
			haveSpecial = true
		}

		if strings.Contains(littleChar, ch) {
			haveLittleChar = true
		}

		if strings.Contains(bigChar, ch) {
			haveBigChar = true
		}

		if strings.Contains(numberChar, ch) {
			haveNumber = true
		}
	}

	return haveSpecial && haveLittleChar && haveBigChar && haveNumber
}

func RandSecret(n int) string {
	rand.Seed(time.Now().UnixNano())

	if n <= 8 {
		n = 12
	}

	for {
		b := make([]byte, n)
		for i := range b {
			b[i] = passwordBase[rand.Intn(len(passwordBase))]
		}

		result := IsStrongPassword(string(b))
		if result {
			return string(b)
		}
	}
}

func ParseStringToUrls(targets ...string) []string {

	var urls []string
	for _, target := range targets {
		target = strings.TrimSpace(target)
		_t := strings.ToLower(target)
		if strings.HasPrefix(_t, "https://") || strings.HasPrefix(_t, "http://") {
			urls = append(urls, target)
			continue
		}

		rawHost, port, err := ParseStringToHostPort(target)
		if err != nil {
			urls = append(urls, fmt.Sprintf("https://%v", target))
			urls = append(urls, fmt.Sprintf("http://%v", target))
			continue
		}

		if port == 80 {
			urls = append(urls, fmt.Sprintf("http://%v", rawHost))
			continue
		}

		if port == 443 {
			urls = append(urls, fmt.Sprintf("https://%v", rawHost))
			continue
		}

		urls = append(urls, fmt.Sprintf("https://%v:%v", rawHost, port))
		urls = append(urls, fmt.Sprintf("http://%v:%v", rawHost, port))
	}

	return urls
}

type blockParser struct {
	scanner *bufio.Scanner
}

func NewBlockParser(reader io.Reader) *blockParser {
	s := bufio.NewScanner(reader)
	s.Split(bufio.ScanWords)
	return &blockParser{scanner: s}
}

func (b *blockParser) NextStringBlock() string {
	b.scanner.Scan()
	return b.scanner.Text()
}

func (b *blockParser) NextBytesBlock() []byte {
	b.scanner.Scan()
	return b.scanner.Bytes()
}

func (b *blockParser) Next() bool {
	return b.scanner.Scan()
}

func (b *blockParser) GetString() string {
	return b.scanner.Text()
}

func (b *blockParser) GetBytes() []byte {
	return b.scanner.Bytes()
}

func (b *blockParser) GetScanner() *bufio.Scanner {
	return b.scanner
}

func DumpHostFileWithTextAndFiles(raw string, divider string, files ...string) (string, error) {
	l := PrettifyListFromStringSplited(raw, divider)
	return DumpFileWithTextAndFiles(
		strings.Join(ParseStringToHosts(strings.Join(l, ",")), divider),
		divider, files...)
}

func DumpFileWithTextAndFiles(raw string, divider string, files ...string) (string, error) {
	// 构建 targets
	targets := strings.Join(ParseStringToLines(raw), divider)
	fp, err := ioutil.TempFile("", "tmpfile-*.txt")
	if err != nil {
		return "", err
	}
	fp.WriteString(targets + divider)
	defer func() {
		fp.Close()
	}()
	for _, f := range files {
		raw, _ := ioutil.ReadFile(f)
		if raw == nil {
			continue
		}
		targetsFromFile := strings.Join(ParseStringToLines(string(raw)), divider)
		targetsFromFile += divider
		fp.WriteString(targetsFromFile)
	}
	return fp.Name(), nil
}

func ParseStringToLines(raw string) []string {
	var lines []string

	scanner := bufio.NewScanner(bytes.NewBufferString(raw))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line == "" {
			continue
		} else {
			lines = append(lines, line)
		}
	}
	return lines
}

func ParseStringToCClassHosts(targets string) string {
	var target []string
	var cclassMap = new(sync.Map)
	for _, r := range ParseStringToHosts(targets) {
		if IsIPv4(r) {
			netStr, err := IPv4ToCClassNetwork(r)
			if err != nil {
				target = append(target, r)
				continue
			}
			cclassMap.Store(netStr, nil)
			continue
		}

		if IsIPv6(r) {
			target = append(target, r)
			continue
		}

		ip := GetFirstIPFromHostWithTimeout(5*time.Second, r, nil)
		if ip != "" && IsIPv4(ip) {
			netStr, err := IPv4ToCClassNetwork(ip)
			if err != nil {
				target = append(target, r)
				continue
			}
			cclassMap.Store(netStr, nil)
			continue
		} else {
			target = append(target, r)
		}
	}
	cclassMap.Range(func(key, value interface{}) bool {
		s := key.(string)
		target = append(target, s)
		return true
	})
	return strings.Join(target, ",")
}
