package gateway

import (
	"errors"
	"net"
)

var (
	errNoGateway = errors.New("no gateway found")
)

type Interface struct {
	inte *net.Interface

	Gw net.IP

	Metric int
}

func DiscoverGateway() (*net.IP, error) {
	inte := &Interface{}

	err := inte.discoverGatewayOS()
	if err != nil {
		return nil, nil
	}

	return &inte.Gw, nil
}

func DiscoverInterface() (*Interface, error) {
	inte := &Interface{}

	err := inte.discoverGatewayInterfaceOS()
	if err != nil {
		return nil, nil
	}

	return inte, nil
}
