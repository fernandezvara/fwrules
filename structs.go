package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

func pathFor(t, n string) string {
	return fmt.Sprintf("fwrules/%s/%s", t, n)
}

func pathMachine(name string) string {
	return pathFor("machines", name)
}

func pathGroup(name string) string {
	return pathFor("groups", name)
}

func pathTemplate(name string) string {
	return pathFor("templates", name)
}

// Machine is the minimal data for any defined 'Server' to administer its firewall
type Machine struct {
	Name       string            `json:"name"`
	Interfaces map[string]string `json:"interfaces,omitempty"`
	Template   string            `json:"template"`
	Groups     []string          `json:"groups"`
}

func isIPV4(ip string) bool {
	if strings.Contains(ip, ".") {
		return true
	}
	return false
}

func newMachine() (m Machine) {
	m.Interfaces = make(map[string]string)
	return
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

func (m *Machine) kvPath() string {
	return fmt.Sprintf("fwrules/machines/%s", m.Name)
}

// Template is the minimal configuration passed to the service.
// Since a template can contain n configuration groups, its used to minimize
// the configuration on the client side
type Template struct {
	Name   string   `json:"name"`
	Groups []string `json:"groups"`
}

func (t *Template) kvPath() string {
	return fmt.Sprintf("fwrules/templates/%s", t.Name)
}

// Group have the required configuration as rules for the machines in
type Group struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Rules       map[string]*Rule `json:"rules"`
}

func newGroup() (g Group) {
	g.Rules = make(map[string]*Rule)
	return
}

func (g *Group) kvPath() string {
	return fmt.Sprintf("fwrules/groups/%s", g.Name)
}

func (g *Group) toByte() ([]byte, error) {
	return json.Marshal(g)
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
