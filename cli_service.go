package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/codegangsta/cli"
)

// service is the main struct with all information needed to maintain the
// configuration
type service struct {
	sync.Mutex
	cli           *cli.Context
	client        *Client
	config        *Config
	template      *Template
	machine       Machine
	groups        []string
	groupsStructs map[string]*Group
	neighbours    []string
}

func newService() (s service) {
	s.Lock()
	s.groupsStructs = make(map[string]*Group)
	s.Unlock()
	return
}

func (s *service) machineRegister() {
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

func (s *service) serviceRegister() error {
	return s.client.ServiceRegister()
}

func (s *service) readAndWatchTemplate() {
	var (
		exists bool
		err    error
	)

	for {
		err = s.client.Watch(pathTemplate(s.cli.String("template")))
		if err == nil {
			log.Println("firewall configuration template updated...")
		}
		exists, err = s.client.GetInterface(pathTemplate(s.cli.String("template")), &s.template)
		assertExit("Error marshalling template data", err, 3)
		if exists == false {
			logMsg("Template does not exists on Consul")
			os.Exit(3)
		} else {
			// update group membership
			s.Lock()
			s.groups = []string{}
			s.groupsStructs = make(map[string]*Group)
			s.groups = s.template.Groups
			log.Println("Updated groups: ", s.groups)
			s.Unlock()
			for _, group := range s.groups {
				s.readGroup(group)
				go s.watchGroup(group)
			}
		}
		s.update()
	}
}

func (s *service) readGroup(name string) {
	var err error
	s.Lock()
	if _, ok := s.groupsStructs[name]; ok == false {
		var gg Group
		s.groupsStructs[name] = &gg
	}
	_, err = s.client.GetInterface(pathGroup(name), s.groupsStructs[name])
	s.Unlock()
	assertExit("Error marshalling group data", err, 3)
}

func (s *service) watchGroup(name string) {
	var (
		err    error
		exists bool
	)
	for {
		err = s.client.Watch(pathGroup(name))
		if err == nil {
			log.Printf("firewall group '%s' updated...\n", name)
		}
		s.Lock()
		if _, ok := s.groupsStructs[name]; ok == false {
			var gg Group
			s.groupsStructs[name] = &gg
		}
		exists, err = s.client.GetInterface(pathGroup(name), s.groupsStructs[name])
		s.Unlock()
		assertExit("Error marshalling group data", err, 3)
		if exists == false {
			s.Lock()
			delete(s.groupsStructs, name)
			s.Unlock()
			logMsg(fmt.Sprintf("Group '%s' lost from Consul", name))
		}
		s.update()
	}
}

func (s *service) update() {
	s.Lock()
	fmt.Println("------------------------------------------------------------")
	for _, n := range s.neighbours {
		fmt.Printf("-A INPUT -s %s/32 -j ACCEPT\n", n)
	}
	fmt.Println("------------------------------------------------------------")
	log.Println("Call for update")
	log.Println(s)
	s.Unlock()
}

func (s *service) neighboursMonitor() {
	for {
		services, err := s.client.WatchServiceMembers()
		assert(err)
		s.Lock()
		s.neighbours = []string{}
		for _, service := range services {
			log.Println("Neighbour:", service.Address)
			s.neighbours = append(s.neighbours, service.Address)
		}
		s.Unlock()
	}
}

func fwrulesService(c *cli.Context) {

	if c.String("template") == "" {
		log.Println("Error: Service cannot start if not 'template' set.")
		os.Exit(3)
	}

	s := newService()

	s.cli = c
	s.config = readConfig(c.GlobalString("config"))
	s.client = NewClient(s.config)

	// registers itself as firewall rules service
	assert(s.serviceRegister())

	// registers the machine configuration
	s.machineRegister()

	// reads and monitors the configuration template from consul
	go s.readAndWatchTemplate()

	// monitors other machines in the same firewall cluster that must be reachable
	go s.neighboursMonitor()

	for {
		time.Sleep(10 * time.Second)
		log.Println("Loop ...")
	}

}
