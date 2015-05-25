package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func fwrulesTest(c *cli.Context) {
	config := readConfig(c.GlobalString("config"))
	client := NewClient(config)

	err := client.Set("molina", []byte("pruebecilla"))
	assert(err)

	for {
		err := client.Watch("molina")

		datos, b, err := client.Get("molina")
		if err == nil {
			if b == true {
				fmt.Println("existe!")
				fmt.Println(string(datos))
			} else {
				fmt.Println("no existe")
			}
		}
	}
}
