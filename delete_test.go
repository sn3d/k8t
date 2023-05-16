package k8t_test

import (
	"os"
	"testing"

	"github.com/sn3d/k8t"
)

func Test_Delete(t *testing.T) {

	if os.Getenv("KUBECONFIG") == "" {
		t.Skip("No KUBECONFIG defined")
	}

	// GIVEN: running kind cluster
	cluster, err := k8t.NewFromEnvironment()
	if err != nil {
		t.FailNow()
	}

	// AND: exist service to-be-deleted in cluster
	err = cluster.ApplyFile("testdata/delete-test.yaml")
	if err != nil {
		t.FailNow()
	}

	err = cluster.WaitFor(k8t.ResourceExist("v1", "Service", "to-be-deleted"))
	if err != nil {
		t.FailNow()
	}

	// WHEN: we delete the service
	err = cluster.DeleteFile("testdata/delete-test.yaml")
	if err != nil {
		t.FailNow()
	}

	// THEN: the service should not exist
	err = cluster.WaitFor(k8t.ResourceNotExist("v1", "Service", "to-be-deleted"))
	if err != nil {
		t.FailNow()
	}
}
