package utils

import (
	"context"
	"github.com/ReneKroon/ttlcache"
	"github.com/mdlayher/arp"
	"net"
	"online/common/log"
	"online/common/utils/arptable"
	"strings"
	"sync"
	"time"

	_ "online/common/utils/arptable"
)

func Arp(ifaceName string, target string) (net.HardwareAddr, error) {
	return ArpWithContext(TimeoutContext(5*time.Second), ifaceName, target)
}

func ArpWithTimeout(timeoutContext time.Duration, ifaceName string, target string) (net.HardwareAddr, error) {
	return ArpWithContext(TimeoutContext(timeoutContext), ifaceName, target)
}

var (
	TargetIsLoopback = Errorf("loopback")
)

func ArpWithContext(ctx context.Context, ifaceName string, target string) (net.HardwareAddr, error) {
	r, err := ArpIPAddressesWithContext(ctx, ifaceName, target)
	if err != nil {
		return nil, err
	}

	if r != nil {
		res, ok := r[target]
		if ok {
			return res, nil
		}
	}
	return nil, Errorf("empty result")
}

var (
	arpTableTTLCacheCreateOnce = new(sync.Once)
	arpTableTTLCache           *ttlcache.Cache
)

func ArpIPAddressesWithContext(ctx context.Context, ifaceName string, addrs string) (map[string]net.HardwareAddr, error) {
	//result, _ := arptable.SearchHardware(target)
	//if result != nil {
	//	return result, nil
	//}
	//
	//if IsLoopback(target) {
	//	return nil, TargetIsLoopback
	//}
	arpTableTTLCacheCreateOnce.Do(func() {
		if arpTableTTLCache == nil {
			arpTableTTLCache = ttlcache.NewCache()
			arpTableTTLCache.SetTTL(10 * time.Second)
		}
	})

	ddl, ok := ctx.Deadline()
	if !ok {
		ddl = time.Now().Add(5 * time.Second)
	}

	// 获取 iface，针对这个 iface 创建一个 arp 客户端
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, err
	}
	client, err := arp.Dial(iface)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	_ = client.SetDeadline(ddl)

	// 并发获取 arp 包
	results := new(sync.Map)
	wg := new(sync.WaitGroup)
	for _, target := range ParseStringToHosts(addrs) {
		target := target
		wg.Add(1)
		go func() {
			defer wg.Done()

			if res, ok := arpTableTTLCache.Get(target); ok {
				results.Store(target, res.(net.HardwareAddr))
				return
			}

			hwAddr, err := arptable.SearchHardware(target)
			if err != nil {
				log.Debugf("")
			}
			if hwAddr != nil {
				results.Store(target, hwAddr)
				arpTableTTLCache.Set(target, hwAddr)
				return
			}

			targetIp := net.ParseIP(target)
			if targetIp == nil {
				log.Debugf("invalid target: %s", targetIp)
				return
				//return nil,
			}

			hw, err := client.Resolve(targetIp)
			if err != nil {
				log.Debugf("resolve arp for %v failed: %s", targetIp.String(), err)
			}
			if hw != nil {
				results.Store(target, hw)
				arpTableTTLCache.Set(target, hw)
				return
			}
		}()
	}
	wg.Wait()
	//for {
	//	select {
	//	case <-time.Tick(1 * time.Second):
	//		hw, _ := client.Resolve(targetIp)
	//		if hw != nil {
	//			return hw, nil
	//		}
	//	case <-newCtx.Done():
	//		return nil, Errorf("cannot found hw for %s", targetIp)
	//	}
	//}
	finalResult := make(map[string]net.HardwareAddr)
	results.Range(func(key, value interface{}) bool {
		finalResult[key.(string)] = value.(net.HardwareAddr)
		return true
	})
	return finalResult, nil
}

var (
	ipLoopback = make(map[string]interface{})
)

func init() {
	addrs, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, i := range addrs {
		ret, _ := i.Addrs()
		for _, addr := range ret {
			ipNet, ok := addr.(*net.IPNet)
			if ok {
				ipLoopback[ipNet.IP.String()] = ipNet
			}
		}
	}
}

func IsLoopback(t string) bool {
	ipInstance := net.ParseIP(FixForParseIP(t))
	if ipInstance != nil {
		if ipInstance.IsLoopback() {
			return true
		}
	}

	if strings.HasPrefix(FixForParseIP(t), "127.") {
		return true
	} else {
		_, ok := ipLoopback[FixForParseIP(t)]
		return ok
	}
}
