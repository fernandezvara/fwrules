package main

import (
	"fmt"
	"time"

	"github.com/codegangsta/cli"
)

func fwrulesService(c *cli.Context) {
	config := readConfig(c.GlobalString("config"))
	client := NewClient(config)

	// registers itself as firewall rules service
	err := client.ServiceRegister()
	assert(err)

	// Monitor Firewall Members
	go func() {
		for {
			services, err := client.WatchServiceMembers()
			assert(err)
			for _, service := range services {
				fmt.Println("---")
				fmt.Println(service.Address)
			}
		}
	}()

	// go func() {
	// 	for {
	//
	// 	}
	// }()

	for {
		time.Sleep(10 * time.Second)
		fmt.Println("Loop ...")
	}

	// fmt.Println(config)
	// fmt.Println("works")
}
