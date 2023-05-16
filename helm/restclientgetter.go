package helm

import (
	"io/ioutil"

	"github.com/sn3d/k8t"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

// Custom implementation of helm's RESTClientGetter that is based
// on cluster's configuration. I need it because we're using
// own way ho to resolve configuration
type customRESTClientGetter struct {
	cfg             *rest.Config
	discoveryClient discovery.CachedDiscoveryInterface
	mapper          meta.RESTMapper
	cconf           clientcmd.ClientConfig
}

func newCustomRESTClientGetter(c *k8t.Cluster) *customRESTClientGetter {

	dir, err := ioutil.TempDir("", "gt-helm-*")
	if err != nil {
		panic(err)
	}

	// Prepare a RESTMapper and find GVR
	dc, err := disk.NewCachedDiscoveryClientForConfig(c.RESTConfig(), dir, dir, 0)
	if err != nil {
		panic(err)
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	return &customRESTClientGetter{
		cfg:             c.RESTConfig(),
		discoveryClient: dc,
		mapper:          mapper,
		cconf:           newCustomClientConfig(c),
	}
}

func (rcg *customRESTClientGetter) ToRESTConfig() (*rest.Config, error) {
	return rcg.cfg, nil
}

func (rcg *customRESTClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	return rcg.discoveryClient, nil
}

func (rcg *customRESTClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	return rcg.mapper, nil
}

func (rcg *customRESTClientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return rcg.cconf
}
