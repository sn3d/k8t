# K8T - your kit for e2e Kubernetes testing


The project's goal is to offer high-level, user-friendly functions for 
e2e testing of Kubernetes. 
The K8T's vision is to:

 - provide hi-level easy-to use functions for Kubernetes testing
 - focus on simplicity
 - not covering all use cases, just the basic-once in simple manner
 - not to be test framework itself
 - do not depend on any specific test framework

> **Warning**
>
> Note that K8T in under active development and the APIs may change 
> in a backwards incompatible manner.

## Motivation


E2E testing for Kubernetes can be challenging and requires writing a lot of 
code. Although there is the e2s-kubernetes framework available, my main issue 
with it is that it aims to be a test framework.

Personally, I prefer using BDD-style E2E testing, and I'm satisfied with 
what Ginkgo/Gomega or GoDog offer in that regard.

What I really desire is a simple collection of high-level functions that can 
seamlessly integrate with my preferred BDD-like framework, traditional unit 
tests, or even enable me to build my own testing CLI (similar to cilium connectivity test).


## Getting started

The K8T is a straightforward Go module that can be directly integrated into 
your tests. Just fetch the K8T Go module by executing the following command:

```shell
go get sn3d.com/sn3d/k8t
```

### Apply and delete resource from manifest
The module relies on a cluster instance that provides essential functions 
such as `Apply()`, `Get()`, `List()`, and `Delete()`. Additionally, the 
cluster instance includes a `WaitFor()` function that can be used with 
different checking functions.

Here's an example that applies the `testdata/busybox.yaml` manifest, verifies 
if the busybox pod is running, and finally deletes the pod:

```go
// get the instance for tested cluster (from KUBECONFIG)
cluster,_ := k8t.NewFromEnvironment()

// apply the manifest
cluster.ApplyFile("testdata/busybox.yaml")

// check if pod is running
err := cluster.WaitFor(k8t.PodIsRunning("","busybox-pod"))
if err != nil {
   panic("Pod is not running")
}

cluster.DeleteFile("testdata/busybox.yaml")
```

### Execute command inside cluster

The K8T module provides the `Execf()` function as well as the more detailed 
`ExecWithOpts()` function. These functions can be used when you need to 
execute commands inside a specific container.

Having the ability to execute commands from the test container within the 
tested cluster is useful for E2E testing. It allows us to verify the proper
functioning of DNS or other network components, check storage, and perform
various other tests.

```go
// execute the command and get the result
result := cluster.Execf("busybox", "busybox-container", "nslookup %s", "google.com")

if result.Err != nil {
   panic("cannot execute command")
}

fmt.Printf(result.String())
}
```

The `Execf()` function is a straightforward execution function that accepts 
arguments such as the pod name, container name, and command. The command can 
be formatted using a syntax similar to `Printf()`. It's important to note 
that the pod should be running in the default test namespace, and the 
command is executed with the `/bin/sh` shell.

If you require additional control over the execution, you can use the 
`ExecWithOpts()` function. This function allows you to modify the namespace 
and provides more options for customization.

### Installing Helm charts

K8T also provides support for installing Helm charts.

Helm is extensively used in the Kubernetes ecosystem, and for E2E testing 
scenarios, it's often necessary to install components using Helm. Let's say 
we have a Helm chart located in the testdata/my-helm folder that we want to 
install. We also want to customize certain values, such as 
`deployment.replicaCount`. Here's an example code snippet to accomplish this:

```go
// get the instance for tested cluster (from KUBECONFIG)
cluster,_ := k8t.NewFromEnvironment()

// set values for Helm release
vals := helm.Value{
   "deployment": helm.Value{
      "replicaCount": 3,
   },
}

// install helm chart with values
err := helm.Install(cluster, "testdata/my-helm", vals)
if err != nil {
   panic("chart cannot be installed")
}
```

### Using K8T with Ginkgo/Gomega

One of its notable features is its seamless integration with Ginkgo, a popular 
BDD testing framework. With K8T, you can effortlessly incorporate Kubernetes 
testing capabilities into your Ginkgo test suite.

```go
var _ = Describe("My firt K8T test", func() {
   var cluster *k8t.Cluster

   BeforeEach(func() {
      cluster,_ := k8t.NewFromEnvironment()
   })

   It("should apply manifest", func() {
      err := cluster.ApplyFile("testdata/my-manifest.yaml")
      Expect(err).NotTo(HaveOccurred())
   })
}
```

## Feedback & Bugs

Feedback is more than welcome. Did you found a bug? Is something not behaving 
as expected? Feature or bug, feel free to create [issue](https://github.com/sn3d/kconf/issues).
