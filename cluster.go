package k8t

import (
	"errors"
	"os"

	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// representing running cluster for your tests
type Cluster struct {

	// kubeconfig from which is created REST config & client
	apiConfig *api.Config

	// configuration for k8s REST client
	restConfig *rest.Config

	// often we need k8s client, it's initialized by New()
	k8sClient *kubernetes.Clientset

	// When we test things, we often use separated namespace, so we don't mess up
	// other things This special namespace is used by all the functions automatically
	testNamespace string

	restMapper *restmapper.DeferredDiscoveryRESTMapper
}

// create cluster instance from KUBECONFIG
func NewFromEnvironment() (*Cluster, error) {
	return NewFromFile(os.Getenv("KUBECONFIG"))
}

// create cluster instance from given kubeconfig file
func NewFromFile(path string) (*Cluster, error) {
	kubeconfig, err := clientcmd.LoadFromFile(path)
	if err != nil {
		return nil, err
	}
	return New(kubeconfig)
}

func New(apiConfig *api.Config) (*Cluster, error) {

	restCfg, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (*api.Config, error) {
		return apiConfig, nil
	})
	if err != nil {
		return nil, err
	}

	k8sClient, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(restCfg)
	if err != nil {
		return nil, err
	}
	restMapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discoveryClient))

	testNamespace := getDefaultNamespace(apiConfig)
	if testNamespace == "" {
		return nil, errors.New("test namespace is empty")
	}

	c := &Cluster{
		apiConfig:     apiConfig,
		restConfig:    restCfg,
		k8sClient:     k8sClient,
		restMapper:    restMapper,
		testNamespace: getDefaultNamespace(apiConfig),
	}

	return c, nil
}

// returns low-level client rest config
func (c *Cluster) RESTConfig() *rest.Config {
	return c.restConfig
}

func (c *Cluster) APIConfig() *api.Config {
	return c.apiConfig
}

func (c *Cluster) TestNamespace() string {
	return c.testNamespace
}

// When user didn't specify test namespace, then the cluster will
// be using default namespace from kubeconfi default namespace from kubeconfigg
func getDefaultNamespace(apiConfig *api.Config) string {
	ctx, ok := apiConfig.Contexts[apiConfig.CurrentContext]
	if !ok {
		return ""
	}

	return ctx.Namespace
}
