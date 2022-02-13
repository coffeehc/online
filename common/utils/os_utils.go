package utils

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"math/rand"
	"net"
	"os"
	"os/user"
	"online/common/log"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func IsUDPPortAvailable(p int) bool {
	return IsPortAvailableWithUDP("0.0.0.0", p)
}

func IsTCPPortAvailable(p int) bool {
	return IsPortAvailable("0.0.0.0", p)
}

func GetRandomAvailableTCPPort() int {
RESET:
	randPort := 55000 + rand.Intn(10000)
	if !IsTCPPortOpen("127.0.0.1", randPort) {
		return randPort
	} else {
		goto RESET
	}
}

func GetRandomAvailableUDPPort() int {
RESET:
	randPort := 55000 + rand.Intn(10000)
	if IsUDPPortAvailable(randPort) {
		return randPort
	} else {
		goto RESET
	}
}

func IsUDPPortAvailableWithLoopback(p int) bool {
	return IsPortAvailableWithUDP("127.0.0.1", p)
}

func IsTCPPortAvailableWithLoopback(p int) bool {
	return IsPortAvailable("127.0.0.1", p)
}

func IsPortAvailable(host string, p int) bool {
	lis, err := net.Listen("tcp", HostPort(host, p))
	if err != nil {
		return false
	}
	_ = lis.Close()
	return true
}

func IsTCPPortOpen(host string, p int) bool {
	dialer := net.Dialer{}
	dialer.Timeout = 10 * time.Second
	conn, err := dialer.Dial("tcp", HostPort(host, p))
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func IsPortAvailableWithUDP(host string, p int) bool {
	addr := fmt.Sprintf("%s:%v", host, p)
	lis, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Infof("%s is unavailable: %s", addr, err)
		return false
	}
	defer func() {
		_ = lis.Close()
	}()
	return true
}

func GetRandomLocalAddr() string {
	rand.Seed(time.Now().UnixNano())

	randPort := rand.Intn(5000)
	return fmt.Sprintf("%s:%v", "127.0.0.1", randPort+60000)
}

func GetSystemNameServerList() ([]string, error) {
	client, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return nil, errors.Errorf("get system nameserver list failed: %s", err)
	}
	return client.Servers, nil
}

func GetHomeDir() (string, error) {
	h, _ := os.UserHomeDir()
	if h != "" {
		return h, nil
	}

	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		usr, err := user.Current()
		if err != nil {
			return "", errors.Errorf("get os use failed: %s", err)
		} else {
			homeDir = usr.HomeDir
		}
	}
	return homeDir, nil
}

func GetHomeDirDefault(d string) string {
	home, err := GetHomeDir()
	if err != nil {
		return d
	}
	return home
}

func InDebugMode() bool {
	return os.Getenv("DEBUG") != "" || os.Getenv("PALMDEBUG") != ""
}

func Debug(f func()) {
	if InDebugMode() {
		f()
	}
}
