package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/consul/api"
)

// Client struct maintains the consul client
type Client struct {
	kv        *api.KV
	agent     *api.Agent
	catalog   *api.Catalog
	waitIndex uint64
}

// NewClient returns a Client to operate in Consul
func NewClient(config *Config) *Client {
	//nodes []string, scheme, cert, key, caCert, datacenter string) (*Client, error) {

	conf := api.DefaultConfig()

	conf.Scheme = config.Scheme
	conf.Datacenter = config.Datacenter

	if len(config.BackendNodes) > 0 {
		conf.Address = config.BackendNodes[0]
	}

	tlsConfig := &tls.Config{}
	if config.Cert != "" && config.Key != "" {
		clientCert, err := tls.LoadX509KeyPair(config.Cert, config.Key)
		assertExit("Error connecting to Consul - Certificates problem", err, 2)
		tlsConfig.Certificates = []tls.Certificate{clientCert}
		tlsConfig.BuildNameToCertificate()
	}
	if config.CaCert != "" {
		ca, err := ioutil.ReadFile(config.CaCert)
		assertExit("Error connecting to Consul - Ca Certificate problem", err, 2)
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(ca)
		tlsConfig.RootCAs = caCertPool
	}
	conf.HttpClient.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client, err := api.NewClient(conf)
	assertExit("Error connecting to Consul", err, 2)
	return &Client{
		kv:      client.KV(),
		agent:   client.Agent(),
		catalog: client.Catalog(),
	}
}

// Get pulls the value for the desired key from the backend
func (c *Client) Get(key string) ([]byte, bool, error) {
	pair, _, err := c.kv.Get(key, nil)
	if err != nil {
		return []byte{}, false, err
	}
	if pair != nil {
		return pair.Value, true, nil
	}
	return []byte{}, false, nil
}

// GetInterface pulls the value for the desired key from the backend append
// fills the object passed as pointer
func (c *Client) GetInterface(key string, obj interface{}) (bool, error) {
	pair, _, err := c.kv.Get(key, nil)
	if err != nil {
		return false, err
	}
	if pair != nil {
		err = json.Unmarshal(pair.Value, obj)
		return true, err
	}
	return false, nil
}

// Set writes the desired value into key on the backend
func (c *Client) Set(key string, value []byte) error {
	p := &api.KVPair{Key: key, Value: value}
	_, err := c.kv.Put(p, nil)
	return err
}

// Delete the desired key
func (c *Client) Delete(key string) error {
	_, err := c.kv.Delete(key, nil)
	return err
}

// Watch waits until the refered key changes
func (c *Client) Watch(key string) error {
	opts := api.QueryOptions{WaitIndex: c.waitIndex}

	_, meta, err := c.kv.Get(key, &opts)
	if err == nil {
		c.waitIndex = meta.LastIndex
	}
	return err
}

// ServiceRegister register the 'fwrules' service
func (c *Client) ServiceRegister() error {
	var service api.AgentServiceRegistration

	service.ID = "fwrules"
	service.Name = "fwrules"
	service.Address = "127.0.0.1"

	return c.agent.ServiceRegister(&service)
}

// WatchServiceMembers watchs a service to get its changes
func (c *Client) WatchServiceMembers() ([]*api.CatalogService, error) {
	opts := api.QueryOptions{WaitIndex: c.waitIndex}

	services, meta, err := c.catalog.Service("fwrules", "", &opts)
	if err == nil {
		c.waitIndex = meta.LastIndex
	}
	return services, err
}
