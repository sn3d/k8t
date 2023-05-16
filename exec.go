package k8t

import (
	"bytes"
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/remotecommand"
)

type ExecOpts struct {
	// here you can pass your context. It it's not set, the default
	// context.Background() will be used
	Context context.Context

	// namespace where is pod for exec. If it's not set, the cluster's
	// test namespace will be used.
	Namespace string
}

type ExecResult struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
	Err    error
}

func (e ExecResult) String() string {
	if e.Err != nil {
		return fmt.Sprintf("%s \n %s", e.Stdout.String(), e.Err.Error())
	} else {
		return e.Stdout.String()
	}
}

func errorResult(err error) ExecResult {
	return ExecResult{
		Err: err,
	}
}

// Less verbose version of ExecWithOpts, which will format and executes given
// command with args with '/bin/sh' shell in given pod.
func (c *Cluster) Execf(pod, container, command string, args ...any) ExecResult {
	formatted := fmt.Sprintf(command, args...)
	cmd := []string{
		"/bin/sh",
		"-c",
		formatted,
	}

	return c.ExecWithOpts(pod, container, cmd, ExecOpts{})
}

// function executes command as array (without shell) in given pod
// and container. It's very general function.
func (c *Cluster) ExecWithOpts(pod string, container string, command []string, opts ExecOpts) ExecResult {
	var err error

	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}

	namespace := opts.Namespace
	if namespace == "" {
		namespace = c.testNamespace
	}

	scheme := runtime.NewScheme()
	if err = corev1.AddToScheme(scheme); err != nil {
		return errorResult(err)
	}
	parameterCodec := runtime.NewParameterCodec(scheme)

	execOpts := &corev1.PodExecOptions{
		Command:   command,
		Container: container,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}

	req := c.k8sClient.CoreV1().RESTClient().Post().Resource("pods").Name(pod).Namespace(namespace).SubResource("exec")
	req.VersionedParams(execOpts, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(c.restConfig, "POST", req.URL())
	if err != nil {
		return errorResult(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    true,
	})

	return ExecResult{
		Stdout: stdout,
		Stderr: stderr,
		Err:    err,
	}
}
