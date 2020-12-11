package inet

import (
	"net"
)

type InetInfo struct {
	Inter net.Interface

	GateInfo []net.IPNet
}

func InetInfos() ([]InetInfo, error) {
	inters, err := net.Interfaces()
	if err != nil || len(inters) <= 0 {
		return nil, err
	}

	inetInfos := make([]InteInfo, len(inters))

	for i := 0; i < len(inters); i++ {
		inetInfos[i].Inter = inters[i]
	}

	// get gateway info

	return inetInfos, err
}
