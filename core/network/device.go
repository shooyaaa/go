package network

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/shooyaaa/log"
	"golang.org/x/exp/slices"
)

func Interfaces() ([]net.Interface, error) {
	return net.Interfaces()
}

func MainInterface() (*net.Interface, error) {
	ifc, _ := net.Interfaces()
	for _, ic := range ifc {
		addrs, _ := ic.Addrs()
		for _, addr := range addrs {
			ipNet, _ := addr.(*net.IPNet)
			ip := ipNet.IP
			if !ip.IsLoopback() && ip.IsPrivate() {
				return &ic, nil
			}

		}
	}
	return nil, errors.New("no main interface found")
}

func MainMacAddr() (*net.HardwareAddr, error) {
	ifc, err := MainInterface()
	if err != nil {
		return nil, err
	}
	return &ifc.HardwareAddr, nil
}

func AllHostsOfInterface(ifcName string) ([]net.IP, error) {
	ifc, err := net.InterfaceByName(ifcName)
	if err != nil {
		return nil, err
	}
	addrs, _ := ifc.Addrs()
	for _, addr := range addrs {
		ip, ok := addr.(*net.IPNet)
		if ok {
			if ip.IP.IsMulticast() {
				continue
			}
			if ip.IP.IsLoopback() {
				continue
			}
			ip4 := ip.IP.To4()
			hosts := []net.IP{ip4}
			if ip4 != nil {
				for index, bt := range ip.Mask {
					if bt != 0xff {
						total := len(hosts)
						for m := 0; m < total; m++ {
							parent := hosts[0]
							hosts = hosts[1:]
							for i := 0; i < 255; i++ {
								next := net.IP(make([]byte, 4))
								copy(next, parent)
								next[index] = byte(i)
								hosts = append(hosts, next)
							}
						}
					}
				}
				return hosts, nil
			}
		}
	}
	return nil, nil
}

var hosts []net.IP

func MonitorHosts(ifcName string) {
	if hosts == nil {
		hosts = make([]net.IP, 0)
	}
	ifc, _ := MainInterface()
	ips := _scanHosts(ifc.Name)
	for _, ip := range ips {
		idx := slices.IndexFunc(hosts, func(c net.IP) bool {
			return ip.String() == c.String()
		})

		if idx == -1 {
			log.DebugF("New host %v", ip)
		}
	}
}
func _scanHosts(ifcName string) []net.IP {
	hosts, _ := AllHostsOfInterface(ifcName)
	var wg sync.WaitGroup
	lives := []net.IP{}
	for _, h := range hosts {
		wg.Add(1)
		go func(host net.IP) {
			defer wg.Done()
			if Ping([4]byte{host[0], host[1], host[2], host[3]}, time.Second*3) {
				lives = append(lives, host)
				log.DebugF("host %s is up\n", host)
			} else {
				//log.DebugF("host %s is down\n", host)
			}
		}(h)
	}
	wg.Wait()
	return lives
}
