package k8t

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// Optional options for GetWithOpts function
type GetOpts struct {
	// here you can pass your context. It it's not set, the default
	// context.Background() will be used
	Context context.Context

	// namespace where is the resource. If it's not set, the cluster's
	// testNamespace will be used
	Namespace string
}

// Get unstructured data of given resource in test namespace. Resource is identified by it's name
// apiVersion e.g. ('networking.k8s.io/v1') and kind (e.g. 'NetworkPolicy')
func (c *Cluster) Get(apiVersion, kind, name string) (*unstructured.Unstructured, error) {
	return c.GetWithOpts(apiVersion, kind, name, GetOpts{})
}

// More verbose version of Get() function. Use this function if you want to
// pass own context, or specify the namespace
func (c *Cluster) GetWithOpts(apiVersion, kind, name string, opts GetOpts) (*unstructured.Unstructured, error) {

	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}

	namespace := opts.Namespace
	if namespace == "" {
		namespace = c.testNamespace
	}

	// we need to parse and convert apiVersion and Kind into
	// GroupVersionResource
	gvk := schema.FromAPIVersionAndKind(apiVersion, kind)

	gk := schema.GroupKind{
		Group: gvk.Group,
		Kind:  gvk.Kind,
	}

	rm, err := c.restMapper.RESTMapping(gk, gvk.Version)
	if err != nil {
		return nil, err
	}

	d, err := dynamic.NewForConfig(c.restConfig)
	if err != nil {
		return nil, err
	}

	res, err := d.Resource(rm.Resource).Namespace(c.testNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}
