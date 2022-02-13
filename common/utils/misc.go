package utils

import (
	"github.com/pkg/errors"
	"net"
	"time"
)

func GetTargetAddrLocalAddr(targetAddr string) (string, error) {
	dialer := net.Dialer{Timeout: 3 * time.Second}
	conn, err := dialer.Dial("tcp", targetAddr)
	if err != nil {
		return "", errors.Errorf("dial failed: %s", err)
	}
	localAddr := conn.LocalAddr()
	_ = conn.Close()

	host, _, err := ParseStringToHostPort(localAddr.String())
	if err != nil {
		return "", errors.Errorf("%v parse to host failed: %s", localAddr.String(), err)
	}
	return host, nil
}

func GetTargetAddrInterfaceName(targetAddr string) (string, error) {
	host, err := GetTargetAddrLocalAddr(targetAddr)
	if err != nil {
		return "", err
	}

	ifs, err := net.Interfaces()
	if err != nil {
		return "", errors.Errorf("get interfaces failed: %s", err)
	}

	for _, iface := range ifs {
		addrs, err := iface.Addrs()
		if err != nil {
			return "", errors.Errorf("fetch iface addr failed: %s", err)
		}

		for _, addr := range addrs {
			_, _net, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}

			if i := net.ParseIP(FixForParseIP(host)); i != nil {
				if _net.Contains(i) {
					return iface.Name, nil
				}
			}
		}
	}
	return "", errors.Errorf("cannot found local interface name for %s", targetAddr)
}
