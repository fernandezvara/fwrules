package consul

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/consul/api"
)

// Client struct maintains the consul client using the common interface
type Client struct {
	kv    *api.KV
	agent *api.Agent
}

// New returns a Client interface for use Consul as backend
func New(nodes []string, scheme, cert, key, caCert string) (*Client, error) {
	conf := api.DefaultConfig()

	conf.Scheme = scheme

	if len(nodes) > 0 {
		conf.Address = nodes[0]
	}

	tlsConfig := &tls.Config{}
	if cert != "" && key != "" {
		clientCert, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{clientCert}
		tlsConfig.BuildNameToCertificate()
	}
	if caCert != "" {
		ca, err := ioutil.ReadFile(caCert)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(ca)
		tlsConfig.RootCAs = caCertPool
	}
	conf.HttpClient.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}
	return &Client{kv: client.KV(), agent: client.Agent()}, nil
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
func (c *Client) Watch(key string, waitIndex uint64, stopChan chan bool) (uint64, error) {
	return 0, nil
}

// ServiceRegister register the 'fwrules' service
func (c *Client) ServiceRegister() error {
	var service api.AgentServiceRegistration

	service.ID = "fwrules"
	service.Name = "fwrules"
	service.Address = "127.0.0.1"

	return c.agent.ServiceRegister(&service)
}
