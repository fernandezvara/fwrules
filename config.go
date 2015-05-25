package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config contains the configuration to access the backend
type Config struct {
	Backend      string   `toml:"backend"`
	BackendNodes []string `toml:"backend_nodes"`
	Datacenter   string   `toml:"datacenter"`
	Scheme       string   `toml:"scheme"`
	Cert         string   `toml:"cert"`
	Key          string   `toml:"key"`
	CaCert       string   `toml:"ca_cert"`
	Interfaces   []string `toml:"interfaces"`
}

func readConfig(fileName string) *Config {
	var config Config
	_, err := toml.DecodeFile(fileName, &config)
	if err != nil {
		fmt.Println("Config: ", err)
		os.Exit(1)
	}
	return &config
}
