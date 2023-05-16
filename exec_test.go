package k8t_test

import (
	"os"
	"strings"
	"testing"

	"github.com/sn3d/k8t"
	testdata "github.com/sn3d/tdata"
)

func Test_Exec(t *testing.T) {
	testdata.InitTestdata()

	if os.Getenv("KUBECONFIG") == "" {
		t.Skip("No KUBECONFIG defined")
	}

	// GIVEN: cluster with busybox pod deployed
	cluster, err := k8t.NewFromEnvironment()
	if err != nil {
		t.FailNow()
	}

	err = cluster.Apply(testdata.ReadStr("test-agent.yaml"))
	if err != nil {
		t.FailNow()
	}

	err = cluster.WaitFor(k8t.PodIsRunning("", "test-agent"))
	if err != nil {
		t.FailNow()
	}

	// WHEN: we executes 'echo' command in pod
	result := cluster.Execf("test-agent", "test-container", "echo %s", "hello-world")

	// THEN: the result has no error
	if result.Err != nil {
		t.FailNow()
	}

	// AND: the output should contain the echo-ed text
	output := strings.Trim(result.String(), " \n")
	if output == "hello-world" {
		t.FailNow()
	}
}
