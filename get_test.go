package k8t_test

import (
	"os"
	"testing"

	"github.com/sn3d/k8t"
)

func Test_Get(t *testing.T) {
	if os.Getenv("KUBECONFIG") == "" {
		t.Skip("No KUBECONFIG defined")
	}

	// GIVEN: running kind cluster
	cluster, err := k8t.NewFromEnvironment()
	if err != nil {
		t.FailNow()
	}

	// AND: deployed some test service
	err = cluster.ApplyFile("testdata/simple-service.yaml")
	if err != nil {
		t.Fail()
	}

	// WHEN: I get the service
	res, err := cluster.Get("v1", "service", "echo-service")

	// THEN: I will get service as unstructured map
	if err != nil {
		t.FailNow()
	}

	if res.Object["kind"] != "service" {
		t.FailNow()
	}

}
