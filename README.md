# K8T - your kit for e2e Kubernetes testing

The goal of this project is provide hi-level easy-to-use functions you can 
use for e2e testing of Kubernetes.

The vision of this kit is:

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

Doing E2E testing for Kubernetes isn't easy. It require lot of code. Yes, 
there is e2s-kubernetes framework. My main problem with this framework it
want to be test framework. 

I don't want to use yet-another test framework. TBH I'm very happy with BDD 
style of E2E testing. I'm happy what Ginkgo/Gomega or GoDog are doing. 

What I wish to have is just set of super-simple hi-level functions I can use 
in my fawourite BDD-like framework, or classy unit-tests, or just build
own CLI that will do testing (e.g. like `cilium connectivity test`).

## Getting started

The K8T is simple Go module, thus could be integrated directly in your 
tests. Simply pull the K8T Go module:

```
go get sn3d.com/sn3d/k8t
```

### Apply and delete resource from manifest

The module relly on cluster instance which is offering you basic functions like
`Apply()`, `Get()`, `List()` or `Delete()`. The cluster instance also offers you
`WaitFor()` function that accept various checking functions.

Following example apply the manifest `manifests/busybox.yaml`, check if busybox 
pod is running and at the end, delete the pod:

```
import (
   "github.com/sn3d/k8t"
   _ "embed"
)


//go:embed manifests/busybox.yaml
var busyboxYAML string

func main() {
   
   // get the instance for tested cluster (from KUBECONFIG)
   cluster,_ := k8t.NewFromEnvironment()

   // apply manifest
   cluster.Apply(busyboxYAML)

   // wait until pod is running
   err := cluster.WaitFor(k8t.PodIsRunning("","busybox-pod"))
   if err != nil {
      panic("Pod is not running")
   }

   // delete the applied pod
   cluster.Delete(busyboxYAML)
}
```
