package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codegangsta/cli"
)

// service is the main struct with all
type service struct {
	cli           *cli.Context
	client        *Client
	config        *Config
	template      *Template
	machine       Machine
	groupsToWatch []string
	neighbours    []string
}

func (s *service) RegisterMachine() {
	// registers machine on Consul
	var (
		b   []byte
		err error
	)
	s.machine = newMachine()
	s.machine.Name = hostname
	s.machine.Template = s.cli.String("template")
	s.machine.getInterfaces(s.config.Interfaces)
	b, err = s.machine.toByte()
	assertExit("Error marshalling machine data", err, 3)
	s.client.Set(s.machine.kvPath(), b)
}

func (s *service) ServiceRegister() error {
	return s.client.ServiceRegister()
}

func (s *service) ReadTemplate() {
	var (
		exists bool
		err    error
	)

	exists, err = s.client.GetInterface(fmt.Sprintf("fwrules/template/%s", s.cli.String("template")), &s.template)
	assertExit("Error marshalling template data", err, 3)
	if exists == false {
		logMsg("Template does not exists on Consul")
	} else {
		s.groupsToWatch = s.template.Groups
	}
}

func fwrulesService(c *cli.Context) {

	var (
		s service
	)

	if c.String("template") == "" {
		log.Println("Error: Service cannot start if not 'template' set.")
		os.Exit(3)
	}

	s.cli = c
	s.config = readConfig(c.GlobalString("config"))
	s.client = NewClient(s.config)

	// registers itself as firewall rules service
	assert(s.ServiceRegister())
	// reads the template from consul
	s.ReadTemplate()

	// var groups []Group

	// Monitor configuration template
	go func() {
		for {
			err := s.client.Watch(fmt.Sprintf("fwrules/templates/%s", c.String("template")))
			if err == nil {
				fmt.Println("firewall configuration template changed...")
			}
		}
	}()

	// Monitor Firewall Members
	go func() {
		for {
			services, err := s.client.WatchServiceMembers()
			assert(err)
			for _, service := range services {
				fmt.Println("---")
				fmt.Println(service.Address)
			}
		}
	}()

	// Monitor group
	go func() {
		for {
			err := s.client.Watch(fmt.Sprintf("fwrules/groups/%s", c.String("group")))
			if err == nil {
				fmt.Println("firewall configuration group changed...")
			}
		}
	}()

	for {
		time.Sleep(10 * time.Second)
		fmt.Println("Loop ...")
	}

}
