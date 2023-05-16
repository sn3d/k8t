package k8t_test

import (
	"os"
	"testing"

	"github.com/sn3d/k8t"
)

func Test_Apply(t *testing.T) {
	if os.Getenv("KUBECONFIG") == "" {
		t.Skip("No KUBECONFIG defined")
	}

	// GIVEN: running kind cluster
	cluster, err := k8t.NewFromEnvironment()
	if err != nil {
		t.FailNow()
	}

	// WHEN: I apply simple pod manifest
	err = cluster.ApplyFile("testdata/simple-pod.yaml")

	// THEN: apply must end with no error
	if err != nil {
		t.FailNow()
	}

}
