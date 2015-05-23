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
	client *api.KV
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
	return &Client{client.KV()}, nil
}

// Watch waits until the refered key changes
func Watch(key string, waitIndex uint64, stopChan chan bool) (uint64, error) {
	return 0, nil
}

// Get pulls the value for the desired key from the backend
func Get(key string) ([]byte, error) {
	return []byte{}, nil
}
