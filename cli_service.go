package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func fwrulesService(c *cli.Context) {
	config := readConfig(c.GlobalString("config"))
	client := NewClient(config)

	err := client.ServiceRegister()
	assert(err)
	fmt.Println(config)
	fmt.Println("works")
}
