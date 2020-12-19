package gateway

import (
	"net"

	"github.com/vishvananda/netlink"
)

func (inte *Interface) discoverRoute() (netlink.Route, error) {
	routes, err := netlink.RouteList(nil, netlink.FAMILY_V4)
	if err != nil {
		return netlink.Route{}, nil
	}

	for i := 0; i < len(routes); i++ {
		if routes[i].Dst == nil {
			return routes[i], nil
		}
	}

	return netlink.Route{}, errNoGateway
}

func (inte *Interface) discoverGatewayOS() error {

	route, err := inte.discoverRoute()
	if err != nil {
		return err
	}
	inte.Gw = route.Gw
	inte.Metric = route.Priority

	return nil
}

func (inte *Interface) discoverGatewayInterfaceOS() error {
	route, err := inte.discoverRoute()
	if err != nil {
		return err
	}
	inte.Gw = route.Gw
	inte.Metric = route.Priority

	inte.Inte, err = net.InterfaceByIndex(route.LinkIndex)

	return err
}
