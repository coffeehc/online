package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ReneKroon/ttlcache"
	"github.com/pkg/errors"
	"math/big"
	"net"
	"online/common/log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func ParseStringToPorts(ports string) []int {
	lports := []int{}

	if strings.HasPrefix(ports, "-") {
		ports = "1" + ports
	}

	if strings.HasSuffix(ports, "-") {
		ports += "65535"
	}

	for _, raw := range strings.Split(ports, ",") {
		raw = strings.TrimSpace(raw)
		if strings.Contains(raw, "-") {
			var (
				low  int64
				high int64
				err  error
			)
			portRange := strings.Split(raw, "-")

			low, err = strconv.ParseInt(portRange[0], 10, 32)
			if err != nil {
				continue
			}

			if portRange[1] != "" {
				high, err = strconv.ParseInt(portRange[1], 10, 32)
				if err != nil {
					continue
				}
			} else {
				continue
			}

			if low > high {
				continue
			}

			for i := low; i <= high; i++ {
				lports = append(lports, int(i))
			}
		} else {
			port, err := strconv.ParseInt(raw, 10, 32)
			if err != nil {
				continue
			}
			lports = append(lports, int(port))
		}
	}

	sort.Ints(lports)
	return lports
}

func SplitHostsToPrivateAndPublic(hosts ...string) (privs, pub []string) {
	for _, host := range hosts {
		if IsPrivateIP(net.ParseIP(FixForParseIP(host))) {
			privs = append(privs, host)
		} else {
			pub = append(pub, host)
		}
	}
	return
}

func ParseStringToHosts(raw string) []string {
	targets := []string{}
	for _, h := range strings.Split(raw, ",") {
		// 解析 IP
		if ret := net.ParseIP(FixForParseIP(h)); ret != nil {
			targets = append(targets, ret.String())
			continue
		}

		// 解析 CIDR 网段
		_ip, netBlock, err := net.ParseCIDR(h)
		if err != nil {
			if strings.Count(h, "-") == 1 {
				// 这里开始解析 1.1.1.1-3 的情况
				rets := strings.Split(h, "-")

				// 检查第一部分是不是 IP 地址
				var startIP net.IP
				if startIP = net.ParseIP(rets[0]); startIP == nil {
					targets = append(targets, h)
					continue
				}

				if strings.Count(rets[0], ".") == 3 {
					ipBlocks := strings.Split(rets[0], ".")
					startInt, err := strconv.ParseInt(ipBlocks[3], 10, 64)
					if err != nil {
						targets = append(targets, h)
						continue
					}

					endInt, err := strconv.ParseInt(rets[1], 10, 64)
					if err != nil {
						targets = append(targets, h)
						continue
					}

					if (endInt > 256) || endInt < startInt {
						targets = append(targets, h)
						continue
					}

					additiveRange := endInt - startInt
					low, err := IPv4ToUint32(startIP.To4())
					if err != nil {
						targets = append(targets, h)
						continue
					}

					for i := 0; i <= int(additiveRange); i++ {
						_ip := Uint32ToIPv4(uint32(i) + low)
						if _ip != nil {
							targets = append(targets, _ip.String())
						}
					}
				} else {
					targets = append(targets, h)
					continue
				}
			} else {
				targets = append(targets, h)
			}
			continue
		}

		// 如果是 IPv6 的网段，暂不处理
		if _ip.To4() == nil {
			targets = append(targets, h)
			continue
		}

		// 把 IPv4 专成 int
		low, err := IPv4ToUint32(netBlock.IP)
		if err != nil {
			targets = append(targets, h)
			continue
		}

		for i := low; true; i++ {
			_ip := Uint32ToIPv4(i)
			if netBlock.Contains(_ip) {
				targets = append(targets, _ip.String())
			} else {
				break
			}
		}
	}

	return StringArrayFilterEmpty(targets)
}

func IPv4ToUint32(ip net.IP) (uint32, error) {
	if len(ip) == 4 {
		return binary.BigEndian.Uint32(ip), nil
	} else {
		return 0, errors.Errorf("cannot convert for ip is not ipv4 ip byte len: %d", len(ip))
	}
}

func Uint32ToIPv4(ip uint32) net.IP {
	ipAddr := make([]byte, 4)
	binary.BigEndian.PutUint32(ipAddr, ip)
	return ipAddr
}
func IPv4ToUint64(ip string) (int64, error) {
	if strings.Contains(ip, ":") == false && len(ip) < 16 {
		ret := big.NewInt(0)
		ret.SetBytes(net.ParseIP(ip).To4())
		return ret.Int64(), nil
	}
	return 0, errors.Errorf("not correct ip=%v", ip)
}

func ParseHostToAddrString(host string) string {
	ip := net.ParseIP(host)
	if ip == nil {
		return host
	}

	if ret := ip.To4(); ret == nil {
		return fmt.Sprintf("[%v]", ip.String())
	}

	return host
}

func IsIPv6(raw string) bool {
	raw = FixForParseIP(raw)
	if ip := net.ParseIP(raw); ip != nil {
		return ip.To4() == nil
	}
	return false
}

func IsIPv4(raw string) bool {
	raw = FixForParseIP(raw)
	if ip := net.ParseIP(raw); ip != nil {
		return ip.To4() != nil
	}
	return false
}

func IsHttp(raw string) bool {
	return strings.HasPrefix(strings.TrimSpace(raw), "http://") || strings.HasPrefix(strings.TrimSpace(raw), "https://")
}

func IsGzip(raw []byte) bool {
	return bytes.HasPrefix(raw, []byte{0x1f, 0x8b, 0x08})
}

func HostPort(host string, port interface{}) string {
	return fmt.Sprintf("%v:%v", ParseHostToAddrString(host), port)
}

func FixForParseIP(host string) string {
	// 如果传入了 [::] 给 net.ParseIP 则会失败...
	// 所以这里要特殊处理一下
	if strings.Count(host, ":") >= 2 {
		if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
			return host[1 : len(host)-1]
		}
	}
	return host
}

type HostsFilter struct {
	origin []string

	strActions []string
	// 这些 Action 如果返回值为 True 说明，在范围内，如果为 False 则不在范围内
	ipActions []func(ip net.IP) bool
}

func (f *HostsFilter) createAction(h string) {
	defaultAction := func(t string) bool {
		return h == t
	}
	// 针对单个 IP 进行处理
	if ret := net.ParseIP(h); ret != nil {
		f.ipActions = append(f.ipActions, func(ip net.IP) bool {
			return ip.String() == ret.String()
		})
		return
	}

	// 解析 CIDR 网段
	_, netBlock, err := net.ParseCIDR(h)
	if err != nil {
		// 如果输入的不是 CIDR 网段
		// 检查 1.1.1.1-3 的情况
		if strings.Count(h, "-") == 1 {
			// 这里开始解析 1.1.1.1-3 的情况
			rets := strings.Split(h, "-")

			// 检查第一部分是不是 IP 地址
			var startIP net.IP
			if startIP = net.ParseIP(rets[0]); startIP == nil {
				f.strActions = append(f.strActions, h)
				return
			}

			if strings.Count(rets[0], ".") == 3 {
				ipBlocks := strings.Split(rets[0], ".")
				startInt, err := strconv.ParseInt(ipBlocks[3], 10, 64)
				if err != nil {
					f.strActions = append(f.strActions, h)
					return
				}

				endInt, err := strconv.ParseInt(rets[1], 10, 64)
				if err != nil {
					f.strActions = append(f.strActions, h)
					return
				}

				if (endInt > 256) || endInt < startInt {
					f.strActions = append(f.strActions, h)
					return
				}

				additiveRange := endInt - startInt
				startIPInt, err := IPv4ToUint32(startIP.To4())
				if err != nil {
					f.strActions = append(f.strActions, h)
					return
				}

				f.ipActions = append(f.ipActions, func(ret net.IP) bool {
					i, err := IPv4ToUint32(ret.To4())
					if err != nil {
						return defaultAction(ret.String())
					}

					return i >= startIPInt && (startIPInt+uint32(additiveRange)) >= i
				})
				return
			} else {
				f.strActions = append(f.strActions, h)
				return
			}
		} else {
			f.strActions = append(f.strActions, h)
			return
		}
	}

	f.ipActions = append(f.ipActions, func(ip net.IP) bool {
		return netBlock.Contains(ip)
	})
	return
}

func (f *HostsFilter) Add(block ...string) {
	for _, b := range block {
		for _, sub := range strings.Split(b, ",") {
			sub = strings.TrimSpace(sub)
			f.createAction(sub)
		}
	}
}

func (f *HostsFilter) Contains(target string) bool {
	// 如果解析出 IP 优先判断 IP
	if len(f.ipActions) > 0 {
		ret := net.ParseIP(target)
		if ret != nil {
			for _, e := range f.ipActions {
				if e(ret) {
					return true
				}
			}
		}
	}

	for _, b := range f.strActions {
		if b == target {
			return true
		}
	}
	return false
}

func NewHostsFilter(excludeHosts ...string) *HostsFilter {
	f := &HostsFilter{
		origin: excludeHosts,
	}
	f.Add(excludeHosts...)
	return f
}

type PortsFilter struct {
	origin []string

	singlePort []int
	actions    []func(i int) bool
}

func (f *PortsFilter) createAction(ports string) {
	for _, raw := range strings.Split(ports, ",") {
		if strings.HasPrefix(raw, "-") {
			raw = "1" + raw
		}

		if strings.HasSuffix(raw, "-") {
			raw += "65535"
		}

		raw = strings.TrimSpace(raw)
		if strings.Contains(raw, "-") {
			var (
				low  int64
				high int64
				err  error
			)
			portRange := strings.Split(raw, "-")

			low, err = strconv.ParseInt(portRange[0], 10, 32)
			if err != nil {
				continue
			}

			if portRange[1] != "" {
				high, err = strconv.ParseInt(portRange[1], 10, 32)
				if err != nil {
					continue
				}
			} else {
				continue
			}

			if low > high {
				continue
			}

			f.actions = append(f.actions, func(i int) bool {
				ret := int64(i)
				return ret >= low && ret <= high
			})
			return
		} else {
			port, err := strconv.ParseInt(raw, 10, 32)
			if err != nil {
				continue
			}

			f.singlePort = append(f.singlePort, int(port))
		}
	}
}

func (f *PortsFilter) Add(block ...string) {
	for _, b := range block {
		for _, sub := range strings.Split(b, ",") {
			sub = strings.TrimSpace(sub)
			f.createAction(sub)
		}
	}
}

func (f *PortsFilter) Contains(port int) bool {
	if len(f.actions) > 0 {
		for _, h := range f.actions {
			if h(port) {
				return true
			}
		}
	}

	for _, i := range f.singlePort {
		if i == port {
			return true
		}
	}
	return false
}

func NewPortsFilter(blocks ...string) *PortsFilter {
	p := &PortsFilter{
		origin: blocks,
	}
	p.Add(blocks...)
	return p
}

type HostPortClassifier struct {
	idMap *sync.Map
	cache *ttlcache.Cache
}

type hostPortIdentifier struct {
	hF *HostsFilter
	pF *PortsFilter
}

func NewHostPortClassifier() *HostPortClassifier {
	cl := &HostPortClassifier{
		idMap: new(sync.Map),
		cache: ttlcache.NewCache(),
	}
	cl.cache.SetExpirationCallback(func(key string, value interface{}) {
		cl.idMap.Delete(key)
	})
	return cl
}

func (h *HostPortClassifier) AddHostPort(tag string, hosts []string, ports []string, ttl time.Duration) error {
	_, ok := h.cache.Get(tag)
	if ok {
		return errors.Errorf("register host port filter failed: %v", ok)
	}

	hf := NewHostsFilter(hosts...)
	pf := NewPortsFilter(ports...)

	f := &hostPortIdentifier{
		hF: hf,
		pF: pf,
	}

	h.cache.SetWithTTL(tag, f, ttl)
	h.idMap.Store(tag, f)

	return nil
}

func (h *HostPortClassifier) FilterTagByHostPort(host string, port int) []string {
	var r []string
	h.idMap.Range(func(key, value interface{}) bool {
		i, ok := value.(*hostPortIdentifier)
		if !ok {
			log.Errorf("key: %v 's host port filter BUG", key)
			return true
		}

		if i.hF.Contains(host) && i.pF.Contains(port) {
			r = append(r, fmt.Sprint(key))
		}
		return true
	})

	return r
}
