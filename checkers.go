package k8t

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// enables the WaitFor() function to receive and execute custom
// checking functions, allowing the caller to define their own
// conditions for waiting on a specific resource or state to be reached.
type Checker func(context.Context, *Cluster) (bool, error)

// Optional options for WaitForWithOpts function
type WaitForOpts struct {

	// here you can pass your context. It it's not set, the default
	// context.Background() will be used
	Context context.Context

	// maximum duration for which the function will wait for the condition
	// to be met. If it's not set, the default 2 minutes will be used
	Timeout time.Duration

	// represents the duration between invocations of the checker function
	// If it's not set, the defauls 2 seconds will be used
	Interval time.Duration
}

// simple, sort version of WaitForOpts with default values. The timeout is
// 2 minutes and interval is set to 2 seconds. Is't equivalent to:
//
//	WaitForOpts(check, WaitForOpts{})
func (c *Cluster) WaitFor(check Checker) error {
	return c.WaitForWithOpts(check, WaitForOpts{})
}

// the function will check if system met given condition.
// The checking is triggered every n seconds defined by interval in opts.
// The invocation of checking function continue until:
//   - function ends with true
//   - reach timeout limit defined in opts
//
// The function returns nil if the given checking function met the condition.
// If we reach the timeout or some error occured during checking, the function
// returns error.
func (c *Cluster) WaitForWithOpts(check Checker, opts WaitForOpts) error {

	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 120 * time.Second
	}

	interval := opts.Interval
	if interval == 0 {
		interval = 2 * time.Second
	}

	err := wait.PollUntilContextTimeout(ctx, interval, timeout, true, func(innerCtx context.Context) (done bool, err error) {
		return check(innerCtx, c)
	})

	return err
}

// check if resource of given kind and name exists in the given
// cluster, in the cluster's test namespace
func ResourceExist(apiVersion, kind, name string) Checker {
	return func(ctx context.Context, c *Cluster) (bool, error) {
		res, err := c.GetWithOpts(apiVersion, kind, name, GetOpts{Context: ctx})
		if err != nil {
			return false, err
		}

		if res == nil {
			return false, nil
		}
		return true, nil
	}
}

// check if resource of given kind and name doesn't exists in the given
// cluster, in the cluster's test namespace
func ResourceNotExist(apiVersion, kind, name string) Checker {
	return func(ctx context.Context, c *Cluster) (bool, error) {
		res, err := c.GetWithOpts(apiVersion, kind, name, GetOpts{Context: ctx})
		if err != nil {
			return true, nil
		}

		if res == nil {
			return true, nil
		}
		return false, nil
	}
}

// check if given pod is running, that means it's in `Running` phase.
// The function expect also namespace. If you provide empty string, then
// cluster's default test namespace will be used.
func PodIsRunning(namespace, name string) Checker {

	return func(ctx context.Context, c *Cluster) (bool, error) {
		podNamespace := namespace
		if podNamespace == "" {
			podNamespace = c.testNamespace
		}

		pod, err := c.k8sClient.CoreV1().Pods(podNamespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if pod.Status.Phase == corev1.PodRunning {
			return true, nil
		}

		return false, nil
	}
}
