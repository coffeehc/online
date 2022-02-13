package netutil

import (
	"context"
	"github.com/dlclark/regexp2"
	"github.com/google/gopacket/routing"
	"github.com/pkg/errors"
	"net"
	"os/exec"
	"online/common/log"
	"online/common/utils"
	"online/common/utils/netutil/netroute"
	"online/common/utils/netutil/routewrapper"
	"runtime"
	"time"
)

func IsPrivateIPString(target string) bool {
	return utils.IsPrivateIP(net.ParseIP(utils.FixForParseIP(target)))
}

func RouteAndArp(target string) (net.HardwareAddr, error) {
	return RouteAndArpWithTimeout(5*time.Second, target)
}

func RouteAndArpWithTimeout(t time.Duration, target string) (net.HardwareAddr, error) {
	iface, targetIP, _, err := Route(t, target)
	if err != nil {
		return nil, err
	}

	if targetIP.String() == utils.FixForParseIP(target) {
		return iface.HardwareAddr, nil
	}

	return utils.ArpWithTimeout(t, iface.Name, targetIP.String())
}

var (
	DarwinGetawayExtractorRe   = regexp2.MustCompile(`gateway: ([\[\]0-9a-fA-TaskFunc:\.]+)`, regexp2.IgnoreCase|regexp2.Multiline)
	DarwinInterfaceExtractorRe = regexp2.MustCompile(`interface: ([^\s]+)`, regexp2.IgnoreCase|regexp2.Multiline)
)

func Route(timeout time.Duration, target string) (iface *net.Interface, gateway, preferredSrc net.IP, err error) {
	addr := utils.GetFirstIPByDnsWithCache(target, timeout)
	if addr == "" {
		err = errors.Errorf("cannot found domain[%s]'s ip address", target)
		return nil, nil, nil, err
	}

	ip := net.ParseIP(utils.FixForParseIP(addr))
	if ip == nil {
		err = errors.Errorf("ip: %v is invalid", ip)
		return nil, nil, nil, err
	}

	route, err := netroute.New()
	if err == nil {
		return route.Route(ip)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	switch runtime.GOOS {
	case "linux":
		fallback := func() (*net.Interface, net.IP, net.IP, error) {
			log.Infof("using gopacket finding route to: %s", ip)
			router, err := routing.New()
			if err != nil {
				err = errors.Errorf("get route failed: %s", err)
				return nil, nil, nil, err
			}

			return router.Route(ip)
		}

		log.Infof("start to find iproute2 utils...")
		//ipUtil, err := exec.LookPath("ip")
		//if err != nil {
		//	log.Infof("start to find iproute2 utils... failed: %s", err)
		//	return fallback()
		//}

		cmd := exec.CommandContext(ctx, "ip", "route", "get", target)
		raw, err := cmd.CombinedOutput()
		if err != nil {
			log.Infof("exec iproute2 utils... failed: %s", err)
			return fallback()
		}

		result := Grok(string(raw), `(local +)?(%{IPORHOST:target} +)?( +via +)?%{IPORHOST:gateway} +dev +%{WORD:iface} +src +%{IP:ifaceIp}`)
		routeTarget := result.Get("target")
		_ = routeTarget
		//if routeTarget != target {
		//
		//	return fallback()
		//}

		gatewayIp := result.Get("gateway")
		ifaceName := result.Get("iface")
		ifaceIp := result.Get("ifaceIp")

		log.Infof("iproute2 found iface: %v ifaceIp: %s gIp: %s", ifaceName, ifaceIp, gatewayIp)
		iface, err = net.InterfaceByName(ifaceName)
		if err != nil {
			log.Infof("open net.InterfaceByName: %s failed: %s", iface, err)
			log.Infof("iproute failed: %s", string(raw))
			return fallback()
		}

		iface, gIp, sIp := iface, net.ParseIP(gatewayIp), net.ParseIP(ifaceIp)
		if gIp == nil || sIp == nil {
			return fallback()
		}
		return iface, gIp, sIp, nil
	case "openbsd", "darwin":
		cmd := exec.CommandContext(ctx, "/sbin/route", "-n", "get", ip.String())
		result, err := cmd.CombinedOutput()
		if err != nil {
			err = errors.Errorf("[route -n get %v] failed: %s", ip.String(), err)
			return nil, nil, nil, err
		}

		resultStr := string(result)
		match, err := DarwinGetawayExtractorRe.FindStringMatch(resultStr)
		if err != nil {
			return nil, nil, nil, errors.Errorf("find match failed: %s", err)
		}

		var (
			targetGateway net.IP
			iface         *net.Interface
			srcIp         net.IP
		)
		if match != nil {
			if getawayIp := match.GroupByNumber(1); getawayIp != nil {
				targetGateway = net.ParseIP(utils.FixForParseIP(getawayIp.String()))
			}
		}

		if targetGateway == nil {
			targetGateway = net.ParseIP(utils.FixForParseIP(target))
		}

		if targetGateway == nil {
			return nil, nil, nil, errors.Errorf("getaway [%s] is invalid/empty")
		}

		match, err = DarwinInterfaceExtractorRe.FindStringMatch(resultStr)
		if err != nil {
			return nil, nil, nil, errors.Errorf("find interface failed: %s", err)
		}
		if match == nil {
			return nil, nil, nil, errors.New("no match found for interface")
		}

		if ifaceName := match.GroupByNumber(1); ifaceName != nil {
			iface, err = net.InterfaceByName(ifaceName.String())
			if err != nil {
				return nil, nil, nil, errors.Errorf("get iface failed: %s", err)
			}

			addrs, err := iface.Addrs()
			if err != nil {
				return nil, nil, nil, errors.Errorf("iface: %v cannot get address: %s", iface.Name, err)
			}
			for _, addr := range addrs {
				raw := utils.FixForParseIP(addr.String())
				srcIpAddress, _, err := net.ParseCIDR(raw)
				if err != nil {
					continue
				}
				if utils.IsIPv6(srcIpAddress.String()) == utils.IsIPv6(targetGateway.String()) {
					srcIp = srcIpAddress
				}
			}
		} else {
			return nil, nil, nil, errors.New("cannot found interface ip")
		}

		return iface, targetGateway, srcIp, err
	default:
		var handleRoute = func(rs []routewrapper.Route) (*net.Interface, net.IP, net.IP, error) {
			for _, route := range rs {
				var srcIp net.IP
				if route.Destination.Contains(net.ParseIP(utils.FixForParseIP(target))) {
					// 获取 IP 地址
					addrs, err := route.Interface.Addrs()
					if err != nil {
						return nil, nil, nil, errors.Errorf("iface: %v cannot get address: %s", iface.Name, err)
					}
					for _, addr := range addrs {
						raw := utils.FixForParseIP(addr.String())
						srcIpAddress, _, err := net.ParseCIDR(raw)
						if err != nil {
							continue
						}
						if utils.IsIPv6(srcIpAddress.String()) == utils.IsIPv6(route.Gateway.String()) {
							srcIp = srcIpAddress
						}
					}
					return route.Interface, route.Gateway, srcIp, nil
				}
			}
			return nil, nil, nil, utils.Errorf("handle route failed: %s", err)
		}

		r, err := routewrapper.NewRouteWrapper()
		if err != nil {
			return nil, nil, nil, utils.Errorf("windows err: %v", err)
		}
		routes, err := r.Routes()
		if err != nil {
			return nil, nil, nil, utils.Errorf("fetch routes failed: %s", err)
		}

		ifaceIns, gatewayIns, ipLocal, err := handleRoute(routes)
		if err != nil {
			routes, err = r.DefaultRoutes()
			if err != nil {
				return nil, nil, nil, utils.Errorf("get default routes failed: %s", err)
			}
			return handleRoute(routes)
		}
		return ifaceIns, gatewayIns, ipLocal, nil
	}
}
