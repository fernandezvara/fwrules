package main

// Config contains the configuration to access the backend
type Config struct {
	Backend      string
	BackendNodes []string
	Scheme       string
	Cert         string
	Key          string
	CaCert       string
}
