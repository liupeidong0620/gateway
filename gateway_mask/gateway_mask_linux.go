// +build linux
package gm

import (
	"bufio"
	"bytes"
	"fmt"
	//"io/ioutil"
	//"net"
	"io"
	"os"
	"strconv"
)

const (
	// See http://man7.org/linux/man-pages/man8/route.8.html
	file = "/proc/net/route"
)

const (
	destinationField = 1
	gatewayField     = 2
)

func (r *GatewayMask) parseRoute() error {
	/*  // cat /proc/net/route
	Iface	Destination	Gateway		Flags	RefCnt	Use	Metric	Mask		MTU	Window	IRTT
	ens33	00000000	0291A8C0	0003	0	0	0	00000000	0	0	0
	ens33	0091A8C0	00000000	0001	0	0	0	00FFFFFF	0	0	0
	*/
	var dataLine []byte
	var err error
	var headerFlag bool

	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("Can't access %s", file)
	}
	defer f.Close()

	buff := bufio.NewReader(f)

	for {
		dataLine, _, err = buff.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if headerFlag == false {
			headerFlag = true
			continue
		}

		fields := bytes.Fields(dataLine)
		if fields == nil || len(fields) < 11 {
			continue
		}
		/*
			for i := 0; i < len(fields); i++ {
				fmt.Println(string(fields[i]))
			}*/
		destinationHex := "0x" + string(fields[destinationField])
		destination, err := strconv.ParseInt(destinationHex, 0, 64)
		if err != nil {
			return fmt.Errorf(
				"parsing destination field hex '%s' in row '%s': %w",
				destinationHex,
				destinationField,
				err,
			)
		}
		gatewayHex, err := "0x" + string(fields[gatewayField])
		gateway, err := strconv.ParseInt(gatewayHex, 0, 64)
		if err != nil {
			return fmt.Errorf(
				"parsing gateway field hex '%s' in row '%s': %w",
				gatewayHex,
				gatewayField,
				err,
			)
		}

		if gateway != 0 || destination != 0 {
			continue
		}

	}

	return nil
}

func (r *GatewayMask) addMask(fields [][]byte) {
	if r.Mask == nil {
	}
}

func (r *GatewayMask) addDefaultGateway(fields [][]byte) {
}
