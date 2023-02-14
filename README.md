# s3gw-operator

This is a demo of a Kubernetes operator using s3gw to provision buckets.

More information on Kubebuilder framework can be found via the
[Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

This project has been initialized with:

```sh
mkdir s3gw-operator
cd s3gw-operator
kubebuilder init --domain s3gw.io --repo github.com/giubacc/s3gw-operator
kubebuilder create api --group s3 --version v1 --kind Bucket
```

## Description

You can create and delete S3 buckets using the CRD: Bucket

Example

```sh
kubectl apply -f bucket.yaml
kubectl get buckets.s3.s3gw.io bucket-sample
kubectl delete buckets.s3.s3gw.io bucket-sample
```

## Getting Started

You’ll need a Kubernetes cluster to run against.
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

Example

```sh
make docker-build docker-push IMG=ghcr.io/giubacc/s3gw-operator:latest
```

- Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/s3gw-operator:tag
```

Example

```sh
make deploy IMG=ghcr.io/giubacc/s3gw-operator
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
