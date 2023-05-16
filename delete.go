package k8t

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

// options for DeleteWithOpts function
type DeleteOpts struct {
	// here you can pass your context. It it's not set, the default
	// context.Background() will be used
	Context context.Context

	// namespace where delete resource. If it's not set, the cluster's
	// test namespace will be used. This is ignored for cluster-wide resources
	Namespace string
}

// Delete resource of given YAML. The resource must be located in
// cluster's default test namespace or it must be cluster-wide resource. If
// you want to delete resource in another namespace, you should use more
// verbose DeleteWithOpts.
func (c *Cluster) Delete(yml string) error {
	return c.DeleteWithOpts(yml, DeleteOpts{})
}

// More verbose version of Delete() function that allow you pass the context
// or change the namespace etc.
func (c *Cluster) DeleteWithOpts(yml string, opts DeleteOpts) error {

	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}

	namespace := opts.Namespace
	if namespace == "" {
		namespace = c.testNamespace
	}

	// Decode YAML manifest into unstructured.Unstructured
	decUnstructured := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode([]byte(yml), nil, obj)
	if err != nil {
		return err
	}

	// Prepare a RESTMapper and find GVR
	dc, err := discovery.NewDiscoveryClientForConfig(c.restConfig)
	if err != nil {
		return err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	// Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(c.restConfig)
	if err != nil {
		return err
	}

	// Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources are placed into test Namespace
		dr = dyn.Resource(mapping.Resource).Namespace(c.testNamespace)
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	// Create or Update the object with server-side-apply
	err = dr.Delete(ctx, obj.GetName(), metav1.DeleteOptions{})

	if err != nil {
		return err
	} else {
		return nil
	}
}
