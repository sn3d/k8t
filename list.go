package k8t

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// Optional options for ListWithOpts function
type ListOpts struct {
	// here you can pass your context. It it's not set, the default
	// context.Background() will be used
	Context context.Context

	// namespace where is the resource. If it's not set, the cluster's
	// testNamespace will be used
	Namespace string
}

// More verbose version of List() function. Use this function if you want to
// pass own context, or specify the namespace
func (c *Cluster) List(apiVersion, kind, labelSelector string) (*unstructured.UnstructuredList, error) {
	return c.ListWithOpts(apiVersion, kind, labelSelector, ListOpts{})
}

// More verbose version of List() function. Use this function if you want to
// pass own context, or specify the namespace
func (c *Cluster) ListWithOpts(apiVersion, kind, labelSelector string, opts ListOpts) (*unstructured.UnstructuredList, error) {

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

	res, err := d.Resource(rm.Resource).Namespace(c.testNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}
