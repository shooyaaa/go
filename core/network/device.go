package network

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/shooyaaa/log"
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
			if !ip.IP.IsPrivate() {
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
							for i := 1; i < 255; i++ {
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

type HostStatus struct {
	Name    string
	Ip      net.IP
	dispear int
}

var hosts map[string]HostStatus

func (hs *HostStatus) Occur() {
	hs.dispear = -1
}

func (hs *HostStatus) DispearTime() int {
	return hs.dispear
}

func (hs *HostStatus) Dispear() {
	hs.dispear += 1
}

func (hs *HostStatus) Down() bool {
	return hs.dispear > 5
}

func MonitorHosts() ([]HostStatus, []HostStatus) {
	firstRun := false
	if hosts == nil {
		hosts = make(map[string]HostStatus)
		firstRun = true
	}
	ifc, _ := MainInterface()
	ips := _scanHosts(ifc.Name)
	upHost := []HostStatus{}
	downHost := []HostStatus{}
	for _, ip := range ips {
		hs, ok := hosts[ip.String()]
		if ok {
			hs.Occur()
		} else {
			names, _ := net.LookupAddr(ip.String())
			name := ip.String()
			if len(names) > 0 {
				name = names[0]
			}
			hosts[ip.String()] = HostStatus{Ip: ip, Name: name, dispear: -1}
			if !firstRun {
				log.DebugF("host %s is up", name)
				upHost = append(upHost, hosts[ip.String()])
			}
		}
	}

	for _, host := range hosts {
		if host.DispearTime() > -1 {
			log.DebugF("host %s disappear %d times ", host.Name, host.DispearTime())
		}
		host.Dispear()
		if host.Down() {
			log.DebugF("host %s is down", host.Name)
			downHost = append(downHost, host)
		}
	}
	return upHost, downHost
}
func _scanHosts(ifcName string) []net.IP {
	hosts, _ := AllHostsOfInterface(ifcName)
	var wg sync.WaitGroup
	lives := []net.IP{}
	for _, h := range hosts {
		wg.Add(1)
		go func(host net.IP) {
			defer wg.Done()
			if Ping([4]byte{host[0], host[1], host[2], host[3]}, time.Second*20) {
				lives = append(lives, host)
				//log.DebugF("host %s is up\n", host)
			} else {
				//log.DebugF("host %s is down\n", host)
			}
		}(h)
	}
	wg.Wait()
	return lives
}
