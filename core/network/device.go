package network

import "net"

func Interfaces() ([]net.Interface, error) {
	return net.Interfaces()
}
