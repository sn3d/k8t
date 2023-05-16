package k8t

import (
	"testing"

	"k8s.io/client-go/tools/clientcmd"
)

func Test_getDefaultNamespace(t *testing.T) {
	// GIVEN: some kubeconfig where current context have
	// detault namespace set
	kubeconfig, err := clientcmd.LoadFromFile("testdata/namespace-kubeconfig.yaml")
	if err != nil {
		t.FailNow()
	}

	// WHEN: I call getDefaultNamespace()
	defaultNs := getDefaultNamespace(kubeconfig)

	// THEN: I should get the namespace
	if defaultNs != "team-b" {
		t.FailNow()
	}
}
