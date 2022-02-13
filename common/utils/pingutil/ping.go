package pingutil

import (
	"net"
	"online/common/log"
	"time"
)

import "github.com/tatsushid/go-fastping"

type PingResult struct {
	IP     string
	Ok     bool
	RTT    int64
	Reason string
}

func PingNative(ip string, timeout time.Duration) *PingResult {
	core := fastping.NewPinger()
	err := core.AddIP(ip)
	if err != nil {
		return &PingResult{
			IP:     ip,
			Ok:     false,
			RTT:    0,
			Reason: err.Error(),
		}
	}

	var result = &PingResult{IP: ip, Reason: "initialized"}

	core.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		if addr.String() == ip {
			result.Ok = true
			result.RTT = int64(rtt) / int64(time.Millisecond)
			result.Reason = ""
		}
	}
	core.OnIdle = func() {

	}

	errChan := make(chan error, 1)
	go func() {
		defer close(errChan)
		err := core.Run()
		if err != nil {
			log.Error(err.Error())
			return
		}
	}()

	select {
	case err, _ := <-errChan:
		if err != nil {
			log.Errorf("ping native mode failed: %s", err)
			return &PingResult{
				IP:     ip,
				Ok:     false,
				RTT:    0,
				Reason: err.Error(),
			}
		}
	case <-time.After(timeout):
		log.Info("timeout ping for %v", ip)
		core.Stop()
	}

	return result
}
