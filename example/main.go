package main

import (
	"fmt"

	"github.com/liupeidong0620/gateway"
)

func main() {
	gw, err := gateway.DiscoverGateway()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("default gateway: ", gw.String())

	inte, err := gateway.DiscoverInterface()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("eht name: ", inte.Inte.Name, ", default gateway:", inte.Gw.String())
}
