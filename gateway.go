package gateway

import (
	"errors"
	"net"
	"runtime"
)

var (
	errNoGateway      = errors.New("no gateway found")
	errNotImplemented = errors.New("not implemented for OS: " + runtime.GOOS)
)

type Interface struct {
	Inte *net.Interface

	Gw net.IP

	Metric int
}

func DiscoverGateway() (*net.IP, error) {
	inte := &Interface{}

	err := inte.discoverGatewayOS()
	if err != nil {
		return nil, err
	}

	return &inte.Gw, nil
}

func DiscoverInterface() (*Interface, error) {
	inte := &Interface{}

	err := inte.discoverGatewayInterfaceOS()
	if err != nil {
		return nil, err
	}

	return inte, nil
}
