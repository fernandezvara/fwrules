package main

import (
	"log"
	"time"

	"github.com/codegangsta/cli"
)

func fwrulesService(c *cli.Context) {

	var s service

	s.cli = c
	s.config = readConfig(c.GlobalString("config"))
	s.client = NewClient(s.config)

	// registers itself as firewall rules service
	assert(s.serviceRegister())

	// registers the machine configuration
	s.machineRegister()

	// reads and monitors the machine configuration from consul
	go s.readAndWatchMachine()

	// reads and monitors the configuration template from consul
	go s.readAndWatchRuleSet()

	// monitors other machines in the same firewall cluster that must be reachable
	go s.neighboursMonitor()

	for {
		time.Sleep(1 * time.Minute)
		log.Println("Loop ...")
	}

}
