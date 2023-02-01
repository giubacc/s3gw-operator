# s3gw-operator

This is a demo of a Kubernetes operator for s3gw to provision buckets.

More information on Kubebuilder framework can be found via the
[Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## Description

You can create S3 buckets using the CRD: Bucket

## Getting Started

You’ll need a Kubernetes cluster to run against.
You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster
for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current
context in your kubeconfig file
(i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster

- Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

- Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/s3gw-operator:tag
```

- Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/s3gw-operator:tag
```

### Uninstall CRDs

To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller

UnDeploy the controller from the cluster:

```sh
make undeploy
```

### How it works

This project aims to follow the Kubernetes
[Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses
[Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources
until the desired state is reached on the cluster.

### Test It Out

- Install the CRDs into the cluster:

```sh
make install
```

- Run your controller (this will run in the foreground, so switch
to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions

If you are editing the API definitions, generate the manifests
such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets
