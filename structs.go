package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/codegangsta/cli"
)

// Machine is the minimal data for any defined 'Server' to administer its firewall
type Machine struct {
	Name       string            `json:"name"`
	Interfaces map[string]string `json:"interfaces,omitempty"`
	RuleSet    string            `json:"ruleset"`
}

func newMachine() *Machine {
	return &Machine{
		Interfaces: make(map[string]string),
	}
}

func (m *Machine) getInterfaces(configInterfaces []string) {
	interfaces, err := net.Interfaces()
	assertExit("Cannot get network interfaces from the machine. Do you have enough permissions?", err, 2)
	for _, inter := range interfaces {

		// continue only if we want that interface
		if stringInSlice(inter.Name, configInterfaces) == false {
			continue
		}

		addrs, err := inter.Addrs()
		assert(err)

		for _, addr := range addrs {
			if isIPV4(addr.String()) {
				m.Interfaces[inter.Name] = strings.Split(addr.String(), "/")[0]
			}
		}
	}
}

func (m *Machine) toByte() ([]byte, error) {
	return json.Marshal(m)
}

// RuleSet is the minimal configuration passed to the service.
// Since a ruleset can contain n configuration groups, its used to minimize
// the configuration on the client side
type RuleSet struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Rules       []*Rule `json:"rules"`
}

// Rule is the definition of a IPtables rule
type Rule struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	PortStart   uint8  `json:"port_start"`
	PortEnd     uint8  `json:"port_end"`
	From        string `json:"from"`
	Protocol    string `json:"protocol"`
}

func pathMachine(fwid, name string) string {
	return pathFor(fwid, "machines", name)
}

func pathRuleSet(fwid, name string) string {
	return pathFor(fwid, "rulesets", name)
}

func pathFor(fwid, t, n string) string {
	return fmt.Sprintf("fwrules/%s/%s/%s", fwid, t, n)
}

// service is the main struct with all information needed to maintain the
// configuration
type service struct {
	sync.Mutex
	cli        *cli.Context
	client     *Client
	config     *Config
	ruleset    *RuleSet
	machine    *Machine
	neighbours []string
}

func (s *service) serviceRegister() error {
	return s.client.ServiceRegister(fmt.Sprintf("fwrules-%s", s.config.FWID))
}

func (s *service) readAndWatchRuleSet() {
	var (
		exists bool
		err    error
	)

	for {
		err = s.client.Watch(pathRuleSet(s.config.FWID, s.config.RuleSet))
		if err == nil {
			log.Println("firewall configuration ruleset updated...")
		}
		exists, err = s.client.GetInterface(pathRuleSet(s.config.FWID, s.config.RuleSet), &s.ruleset)
		assertExit("Error marshalling ruleset data", err, 3)
		if exists == false {
			logMsg("RuleSet does not exists on Consul")
			os.Exit(3)
		}
		s.update()
	}
}

func (s *service) readAndWatchMachine() {
	var (
		exists bool
		err    error
	)

	for {
		err = s.client.Watch(pathRuleSet(s.config.FWID, s.config.RuleSet))
		if err == nil {
			log.Println("machine configuration updated...")
		}
		exists, err = s.client.GetInterface(pathMachine(s.config.FWID, hostname), &s.machine)
		assertExit("Error marshalling machine data", err, 3)
		if exists == false {
			logMsg("Machine does not exists on Consul")
			os.Exit(3)
		}
		s.update()
	}
}

func (s *service) machineRegister() {
	// registers machine on Consul
	var (
		b   []byte
		err error
	)
	s.machine = newMachine()
	s.machine.Name = hostname
	s.machine.RuleSet = s.config.RuleSet
	s.machine.getInterfaces(s.config.Interfaces)
	b, err = s.machine.toByte()
	assertExit("Error marshalling machine data", err, 3)
	assert(s.client.Set(pathMachine(s.config.FWID, hostname), b))
}

func (s *service) neighboursMonitor() {
	for {
		services, err := s.client.WatchServiceMembers(fmt.Sprintf("fwrules-%s", s.config.FWID))
		assert(err)
		s.Lock()
		s.neighbours = []string{}
		for _, service := range services {
			log.Println("Neighbour:", service.Address)
			s.neighbours = append(s.neighbours, service.Address)
		}
		s.Unlock()
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
