package helm_test

import (
	"os"
	"testing"

	"github.com/sn3d/k8t"
	"github.com/sn3d/k8t/helm"
	testdata "github.com/sn3d/tdata"
)

func Test_HelmInstall(t *testing.T) {
	testdata.InitTestdata()
	if os.Getenv("KUBECONFIG") == "" {
		t.Skip("No KUBECONFIG defined")
	}

	// GIVEN: running cluster
	cluster, err := k8t.NewFromEnvironment()
	if err != nil {
		t.FailNow()
	}

	//  WHEN: I install the 'demo' helm chart with some values
	vals := helm.Value{
		"deployment": helm.Value{
			"replicaCount": 3,
		},
	}

	err = helm.Install(cluster, "testdata/demo", vals)
	if err != nil {
		t.FailNow()
	}

	//  THEN: The helm is installed and 3 pods for 'demo' are available
}
