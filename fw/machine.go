package fw

// Machine is the minimal data for any defined 'Server' to administer its firewall
type Machine struct {
	Name       string               `json:"name"`
	Interfaces map[string]Interface `json:"interfaces,omitempty"`
	AgentIP    string               `json:"agentIP"`
	Groups     []string             `json:"groups"`
}

// Interface maintains the data for every interface defined on the configuration
type Interface struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}

// Group have the required configuration as rules for the machines in
type Group struct {
	Name  string          `json:"name"`
	Rules map[string]Rule `json:"rules"`
}

// Rule is the definition of a IPtables rule
type Rule struct {
	PortStart uint8  `json:"port_start"`
	PortEnd   uint8  `json:"port_end"`
	From      string `json:"from"`
	Protocol  string `json:"protocol"`
  
}
