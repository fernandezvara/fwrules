package main

import (
	"errors"

	"github.com/fernandezvara/fwrules/consul"
)

// Client is the common interface for all operations on backends
type Client interface {
	Watch(key string, waitIndex uint64, stopChan chan bool) (uint64, error)
	Get(key string) ([]byte, error)
}

// NewClient returns an instance to the backend client based on the configuration
func NewClient(config Config) (Client, error) {

	switch config.Backend {
	case "consul":
		return consul.New(config.BackendNodes, config.Scheme, config.Cert, config.Key, config.CaCert)
	}

	return nil, errors.New("Unknown backend")
}
