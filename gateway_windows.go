package gateway

import (
	"fmt"
	"net"
	"syscall"
	"unsafe"
)

var defaultRoute = [4]byte{0, 0, 0, 0}

type routeTable struct {
	iphlpapi          *syscall.LazyDLL
	getIpForwardTable *syscall.LazyProc
}

type RouteRow struct {
	ForwardDest      [4]byte //目标网络
	ForwardMask      [4]byte //掩码
	ForwardPolicy    uint32  //ForwardPolicy:0x0
	ForwardNextHop   [4]byte //网关
	ForwardIfIndex   uint32  // 网卡索引 id
	ForwardType      uint32  //3 本地接口  4 远端接口
	ForwardProto     uint32  //3静态路由 2本地接口 5EGP网关
	ForwardAge       uint32  //存在时间 秒
	ForwardNextHopAS uint32  //下一跳自治域号码 0
	ForwardMetric1   uint32  //度量衡(跃点数)，根据 ForwardProto 不同意义不同。
	ForwardMetric2   uint32
	ForwardMetric3   uint32
	ForwardMetric4   uint32
	ForwardMetric5   uint32
}

func (rr *RouteRow) GetForwardDest() net.IP {
	return net.IP(rr.ForwardDest[:])
}

func (rr *RouteRow) GetForwardMask() net.IP {
	return net.IP(rr.ForwardMask[:])
}

func (rr *RouteRow) GetForwardNextHop() net.IP {
	return net.IP(rr.ForwardNextHop[:])
}

func (rt *routeTable) getRoutes() ([]RouteRow, error) {
	var err error
	buf := make([]byte, 4+unsafe.Sizeof(RouteRow{}))
	buf_len := uint32(len(buf))

	// get route buf size
	rt.getIpForwardTable.Call(uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&buf_len)), 0)

	var r1 uintptr
	for i := 0; i < 5; i++ {
		buf = make([]byte, buf_len)
		r1, _, err = rt.getIpForwardTable.Call(uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(&buf_len)), 0)
		if r1 == 122 {
			continue
		}
		break
	}

	if r1 != 0 {
		return nil, fmt.Errorf("Failed to get the routing table, return value：%v, errMsg: %v", r1, err)
	}

	num := *(*uint32)(unsafe.Pointer(&buf[0]))
	routes := make([]RouteRow, num)
	sr := uintptr(unsafe.Pointer(&buf[0])) + unsafe.Sizeof(num)
	rowSize := unsafe.Sizeof(RouteRow{})

	// 安全检查
	if len(buf) < int((unsafe.Sizeof(num) + rowSize*uintptr(num))) {
		return nil, fmt.Errorf("System error: GetIpForwardTable returns the number is too long, beyond the buffer。")
	}

	for i := uint32(0); i < num; i++ {
		routes[i] = *((*RouteRow)(unsafe.Pointer(sr + (rowSize * uintptr(i)))))
	}

	return routes, nil
}

func (inte *Interface) discoverRoute() (RouteRow, error) {

	iphlpapi := syscall.NewLazyDLL("iphlpapi.dll")
	getIpForwardTable := iphlpapi.NewProc("GetIpForwardTable")

	routeTable := &routeTable{
		iphlpapi:          iphlpapi,
		getIpForwardTable: getIpForwardTable,
	}

	routes, err := routeTable.getRoutes()
	if err != nil {
		return RouteRow{}, errNoGateway
	}

	for i := 0; i < len(routes); i++ {
		if routes[i].ForwardDest == defaultRoute {
			return routes[i], nil
		}
	}

	return RouteRow{}, errNoGateway
}

func (inte *Interface) discoverGatewayOS() error {

	route, err := inte.discoverRoute()
	if err != nil {
		return err
	}
	inte.Gw = route.GetForwardNextHop()
	inte.Metric = int(route.ForwardMetric1)

	return nil
}

func (inte *Interface) discoverGatewayInterfaceOS() error {
	route, err := inte.discoverRoute()
	if err != nil {
		return err
	}
	inte.Gw = route.GetForwardNextHop()
	inte.Metric = int(route.ForwardMetric1)

	inte.inte, err = net.InterfaceByIndex(int(route.ForwardIfIndex))

	return err
}
