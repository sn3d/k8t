package k8t_test

import (
	"os"
	"testing"

	"github.com/sn3d/k8t"
	testdata "github.com/sn3d/tdata"
)

func Test_List(t *testing.T) {
	testdata.InitTestdata()
	if os.Getenv("KUBECONFIG") == "" {
		t.Skip("No KUBECONFIG defined")
	}

	// GIVEN: running kind cluster
	cluster, err := k8t.NewFromEnvironment()
	if err != nil {
		t.FailNow()
	}

	// AND: deployed multiple pods with label app: list-test
	err = cluster.Apply(testdata.ReadStr("list-test.yaml"))
	if err != nil {
		t.Fail()
	}

	// WHEN: I list all pods with labelselector 'app=list-test'
	res, err := cluster.List("v1", "Pod", "app=list-test")
	if err != nil {
		t.FailNow()
	}

	// THEN: I'll get list of 3 pods
	if len(res.Items) != 3 {
		t.FailNow()
	}
}
