package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func kvSet(c *cli.Context) {
	config := readConfig(c.GlobalString("config"))

	client := NewClient(config)

	if len(c.Args()) != 2 {
		fmt.Println("Invalid number of arguments.")
		fmt.Println("<key> <value>")
		os.Exit(2)
	}

	client.Set(c.Args()[0], []byte(c.Args()[1]))
}

func kvGet(c *cli.Context) {
	config := readConfig(c.GlobalString("config"))

	client := NewClient(config)

	if len(c.Args()) != 1 {
		fmt.Println("Invalid number of arguments.")
		fmt.Println("<key>")
		os.Exit(2)
	}

	value, found, err := client.Get(c.Args()[0])
	assert(err)
	if found == true {
		fmt.Println(string(value))
	} else {
		fmt.Println("Not found")
	}
}

func kvDelete(c *cli.Context) {
	config := readConfig(c.GlobalString("config"))

	client := NewClient(config)

	if len(c.Args()) != 1 {
		fmt.Println("Invalid number of arguments.")
		fmt.Println("<key>")
		os.Exit(2)
	}

	err := client.Delete(c.Args()[0])
	assert(err)
}
