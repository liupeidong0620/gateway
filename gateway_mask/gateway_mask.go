package gm

type RouteInfo struct {
	Iface       string
	Destination string
	Gateway     string
	Flags       string
	RefCnt      string
	Use         string
	Metric      string
	Mask        string
	MTU         string
	Window      string
	IRTT        string
}

type GatewayMask struct {
	// iface ----> info
	Mask map[string][]RouteInfo
	// default gateway
	GateWay RouteInfo
}

func GatewayMasks() (GatewayMask, error) {
}
