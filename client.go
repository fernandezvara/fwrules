package main

import (
	"log"
	"os"

	"github.com/fernandezvara/fwrules/consul"
)

// Client is the common interface for all operations on backends
type Client interface {
	Get(key string) ([]byte, bool, error)
	Set(key string, value []byte) error
	Delete(key string) error
	Watch(key string, waitIndex uint64, stopChan chan bool) (uint64, error)
	ServiceRegister() error
}

// NewClient returns an instance to the backend client based on the configuration
func NewClient(config *Config) Client {

	switch config.Backend {
	case "consul":
		client, err := consul.New(config.BackendNodes, config.Scheme, config.Cert, config.Key, config.CaCert)
		assertExit("Error connecting to backend", err, 2)
		return client
	}

	log.Println("Unknown backend")
	os.Exit(2)
	return nil
}
