// +build darwin dragonfly freebsd netbsd openbsd

// Package route provides basic functions for the manipulation of
// packet routing facilities on BSD variants.
//
// The package supports any version of Darwin, any version of
// DragonFly BSD, FreeBSD 7 and above, NetBSD 6 and above, and OpenBSD

package gateway

import (
	"net"

	"golang.org/x/net/route"
)

var defaultRoute = [4]byte{0, 0, 0, 0}

func (inte *Interface) discoverRoute() (*route.RouteMessage, error) {
	rib, err := route.FetchRIB(0, route.RIBTypeRoute, 0)
	if err != nil {
		return nil, err
	}
	messages, err := route.ParseRIB(route.RIBTypeRoute, rib)
	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		route_message := message.(*route.RouteMessage)
		addresses := route_message.Addrs

		var destination, gateway *route.Inet4Addr
		ok := false

		if destination, ok = addresses[0].(*route.Inet4Addr); !ok {
			continue
		}

		if gateway, ok = addresses[1].(*route.Inet4Addr); !ok {
			continue
		}

		if destination == nil || gateway == nil {
			continue
		}

		if destination.IP == defaultRoute {
			// fmt.Println(gateway.IP)
			return route_message, nil
		}

	}

	return nil, errNoGateway
}

func (inte *Interface) discoverGatewayOS() error {

	route_, err := inte.discoverRoute()
	if err != nil {
		return err
	}
	addresses := route_.Addrs
	gw, _ := addresses[1].(*route.Inet4Addr)
	routeMetric := route_.Sys()
	inte.Gw = net.IP(gw.IP[:])
	if len(routeMetric) > 0 {
		inte.Metric = int(routeMetric[0].SysType())
	}

	return nil
}

func (inte *Interface) discoverGatewayInterfaceOS() error {
	route_, err := inte.discoverRoute()
	if err != nil {
		return err
	}
	addresses := route_.Addrs
	gw, _ := addresses[1].(*route.Inet4Addr)
	routeMetric := route_.Sys()
	inte.Gw = net.IP(gw.IP[:])
	if len(routeMetric) > 0 {
		inte.Metric = int(routeMetric[0].SysType())
	}

	inte.Inte, err = net.InterfaceByIndex(route_.Index)

	return err
}
