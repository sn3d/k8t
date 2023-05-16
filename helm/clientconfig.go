package helm

import (
	"github.com/sn3d/k8t"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// Custom implementation of clientcmd.ClientConfig. I cannot use DirectClientConfig
// because private fields. For that reason I've created this package scoped
// implementation.
//
// This implementation is currently needed by helm only
type customClientConfig struct {
	cluster *k8t.Cluster
}

func newCustomClientConfig(c *k8t.Cluster) clientcmd.ClientConfig {
	return &customClientConfig{
		cluster: c,
	}
}

func (cc *customClientConfig) RawConfig() (api.Config, error) {
	return *cc.cluster.APIConfig(), nil
}

// ClientConfig returns a complete client config
func (cc *customClientConfig) ClientConfig() (*rest.Config, error) {
	return cc.cluster.RESTConfig(), nil
}

func (cc *customClientConfig) Namespace() (string, bool, error) {
	return cc.cluster.TestNamespace(), false, nil
}

// ConfigAccess returns the rules for loading/persisting the config.
func (c *customClientConfig) ConfigAccess() clientcmd.ConfigAccess {
	return clientcmd.NewDefaultClientConfigLoadingRules()
}
