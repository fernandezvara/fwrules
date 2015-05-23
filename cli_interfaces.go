package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/codegangsta/cli"
)

func fwrulesInterfaces(c *cli.Context) {

	// assertLinux()

	interfaces, err := net.Interfaces()
	assertExit("Cannot get network interfaces from the machine. Do you have permissions?", err, 2)
	for _, int := range interfaces {
		addrs, err := int.Addrs()
		assert(err)
		var ips []string
		var networks []string
		for _, addr := range addrs {
			ips = append(ips, addr.String())
			networks = append(networks, addr.Network())
		}
		fmt.Printf("%s (%s) (%s)\n", int.Name, strings.Join(ips, ","), strings.Join(networks, ","))
	}
}
