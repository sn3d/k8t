package k8t

import (
	"context"
	"encoding/json"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

// options for GetWithOpts function
type ApplyOpts struct {
	// here you can pass your context. It it's not set, the default
	// context.Background() will be used
	Context context.Context

	// namespace where to apply resource. If it's not set, the cluster's
	// test namespace will be used. This is ignored for cluster-wide resources
	Namespace string
}

// Same as kubectl apply. Function apply any YAML string into given cluster,
// into default test namespace.
func (c *Cluster) Apply(yml string) error {
	return c.ApplyWithOpts(yml, ApplyOpts{})
}

// More verbose version of Apply() function where you can pass your own
// context or set namespace
func (c *Cluster) ApplyWithOpts(yml string, opts ApplyOpts) error {

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
		dr = dyn.Resource(mapping.Resource).Namespace(namespace)
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	// Marshal object into JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	// Create or Update the object with server-side-apply
	_, err = dr.Patch(ctx, obj.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{
		FieldManager: "k8t",
	})

	if err != nil {
		return err
	} else {
		return nil
	}
}
