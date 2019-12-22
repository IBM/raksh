# Overview
Raksh (`protect`) project is created with the aim to secure Kubernetes deployed
workload along with its specification by leveraging hardware assisted security.
This development was primarily to leverage Protected Execution Facility (PEF)
capability (similar to AMD SEV) provided by Power (Power9) processors.  Simply
put, PEF provides the ability to secure data-in-use by protecting access to
specific memory regions.

PEF introduces a concept of secure VM (SVM). Anything running inside the VM is
protected. More details on PEF is available here Ref:
https://www.youtube.com/watch?v=pKh_mPPo9X4

Our goal was to leverage the PEF  functionality with Kubernetes to provide a
more secure option for our clients to deploy their containerized application.
Since the protection and isolation is provided by the virtualization layer (KVM
with support for SVM), the natural choice was to leverage Kata containers as the
basis. There are already examples of Kata integration with different
virtualization technologies for improved security and isolation (firecracker
etc)

We also had a need for strong coupling between the container image and the VM
image and go a step further to protect the application spec as well.

We leverage Kata containers. However we use a modified Kata agent with the following
functionalities

1. Support for Decrypting the spec inside the VM
2. Creating the containers based on the decrypted spec


For more details refer to the following links

- Kubecon Europe 2019 [video](https://www.youtube.com/watch?v=pGMdSFlD0_E)
- Kubecon Europe 2019 [ppt](https://static.sched.com/hosted_files/kccnceu19/5c/KubeCon-Europe-2019-protected-memory.pdf)

## Team
[Harshal Patil](https://github.com/harche)
[Manjunath Kumatagi](https://github.com/mkumatag)
[Nitesh Konkar](https://github.com/nitkon)
[Pradipta Banerjee](https://github.com/bpradipt)
[Suhail Anjum](https://github.com/suhailgray)

# Prerequisites

1. Fedora/RHEL/CentOS or Ubuntu host with KVM support
2. Golang 1.12+
3. Docker or Podman for building the containers
4. Kubernetes with Kata Containers version 1.9 as runtime. Ensure Runtimeclass name is set as "kata-containers"
5. CRIO or containerd

# How to build

```shell
# Clone the repository
$ mkdir -p $GOPATH/src/github.com/ibm
$ cd $GOPATH/src/github.com/ibm
$ git clone https://github.com/ibm/raksh.git
$ git checkout -b 1.9.1-raksh origin/1.9.1-raksh

# For building binaries
$ cd raksh
$ make build-binary

# Install rakshctl binary
$ install -D build/_output/bin/rakshctl /usr/bin/

# Build and push docker images to external registry:
$ REGISTRY=docker.io && docker login $REGISTRY
$ REGISTRY=docker.io ORG=projectraksh make build-image
$ REGISTRY=docker.io ORG=projectraksh make push-image
$ REGISTRY=docker.io ORG=projectraksh IMAGE=sc-scratch  make push-manifest
$ VERSION=latest REGISTRY=docker.io ORG=projectraksh IMAGE=sc-scratch make push-manifest
$ REGISTRY=docker.io ORG=projectraksh IMAGE=securecontainer-operator  make push-manifest
$ VERSION=latest REGISTRY=docker.io ORG=projectraksh IMAGE=securecontainer-operator make push-manifest
```

# Building the Kata (Raksh) agent and initrd

Follow the steps mentioned [here](https://github.com/ibm/raksh-agent/README.md) to build the agent and initrd

# Quick Start
## Deploy the securecontainer-operator

Create image registry secret for your setup

```shell
$ kubectl create secret docker-registry regcred --docker-server=<image-registry> --docker-username=<user-name> --docker-password=<password> --docker-email=<email>

```

```shell
# Setup Service Account
$ kubectl create -f deploy/service_account.yaml

# Setup RBAC
$ kubectl create -f deploy/role.yaml
$ kubectl create -f deploy/role_binding.yaml

# Setup the CRD
$ kubectl create -f deploy/crds/securecontainers.k8s.io_securecontainerimageconfigs_crd.yaml
$ kubectl create -f deploy/crds/securecontainers.k8s.io_securecontainerimages_crd.yaml
$ kubectl create -f deploy/crds/securecontainers.k8s.io_securecontainers_crd.yaml

# Deploy the securecontainer-operator operator:
$ kubectl create -f deploy/operator.yaml
```

> Note: Set the image-registry and change the image name for the securecontainer-operator in deploy/operator.yaml to override

## Delete the securecontainer-operator

```shell
$ kubectl delete -f deploy/operator.yaml

$ kubectl delete -f deploy/crds/securecontainers.k8s.io_securecontainerimageconfigs_crd.yaml
$ kubectl delete -f deploy/crds/securecontainers.k8s.io_securecontainerimages_crd.yaml
$ kubectl delete -f deploy/crds/securecontainers.k8s.io_securecontainers_crd.yaml

$ kubectl delete -f deploy/role.yaml
$ kubectl delete -f deploy/role_binding.yaml

$ kubectl delete -f deploy/service_account.yaml
```

## How to build the SecureContainerImage
```
$ rakshctl image create --image <SecureContainerImage_RESOURCE_NAME> --initrd <PATH_TO_KATA-INITRD_IMAGE> --vmlinux <PATH_TO_KATA-KERNEL> --symmKeyFile <PATH_TO_SYMM_KEY_FILE>  --filename <PATH_TO_DEPLOYMENT_FILE> --scratch <IMAGE_REGISTRY/ORG/SCRATCH_IMAGE_NAME> --push <IMAGE_REGISTRY/ORG/IMAGE_NAME>

$ rakshctl image create --image nginx-securecontainerimage --initrd /usr/share/kata-containers/kata-containers-initrd.img --vmlinux /usr/share/kata-containers/vmlinuz.container --symmKeyFile /root/key_file --filename /securecontainers/sample/nginx.yaml --scratch docker.io/projectraksh/sc-scratch:latest --push docker.io/projectraksh/nginx-securecontainerimage:latest 
```

> **Note**:
> For Intel, use compressed kernel image- vmlinuz.  For Power, use uncompressed kernel image- vmlinux.  
> To generate symmKey use: `openssl rand -rand /dev/urandom 32 > symmKey`

## Create SecureContainerImageConfig custom resource

```shell
$ kubectl create -f imageconfig.yaml
```

Example:
```yaml
$ cat imageconfig.yaml
apiVersion: securecontainers.k8s.io/v1alpha1
kind: SecureContainerImageConfig
metadata:
  name: nginx-securecontainerimageconfig
spec:
  imageDir: /var/lib/kubelet/secure-images
  runtimeClassName: kata-containers
```

## Create SecureContainerImage custom resource:

```shell
$ kubectl create -f securecontainerimage.yaml
```

Example:
```yaml
$ cat securecontainerimage.yaml
apiVersion: securecontainers.k8s.io/v1alpha1
kind: SecureContainerImage
metadata:
  name: nginx-securecontainerimage
spec:
  vmImage: projectraksh/nginx-securecontainerimage:latest
  imagePullSecrets:
    - name: regcred
  SecureContainerImageConfigRef:
    name: nginx-securecontainerimageconfig
```

> Note: Replace the vmImage with proper image build in step [How to build the SecureContainerImage](#How-to-build-the-SecureContainerImage)

## Convert the user workload to secure workload and deploy

If you wish to use vault server for key management, please refer to the doc on [Vault Integration with raksh](docs/vault.md).

Example workload:
```yaml
$ cat nginx.yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    app: nginx
  namespace: default
  name: nginx
spec:
  containers:
    - image: nginx:latest
      imagePullPolicy: IfNotPresent
      name: nginx
      ports:
      - containerPort: 80
        protocol: TCP
```

The above [command](#How-to-build-the-SecureContainerImage) will generate a secure workload in nginx-sc.yaml file

```yaml
$ cat nginx-sc.yaml
---
apiVersion: v1
data:
  nginx: ZZ9P65ZlLoK/f9WZDFLq4j8F4piH9yAANmHg+CwC/5oFatk35E77p+DXYY9DX1HU3OqyZ++7+UHV/7XoWuUfgO0p0eVRT8nkF2VZRRSwWIgeBH7RdIYluPqt0TAQt4AFdOc3E1bFyKheWtT+l/JBJLmjSFTSKaZMaV+hb9Ev3WYVd7VDLoI5fhh9v2LE8bH1GfPFRKo=
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: secure-configmap-nginx
  namespace: default
---
apiVersion: securecontainers.k8s.io/v1alpha1
kind: SecureContainer
metadata:
  creationTimestamp: null
  name: secure-nginx
object:
  apiVersion: v1
  kind: Pod
  metadata:
    creationTimestamp: null
    labels:
      app: nginx
    name: nginx
    namespace: default
  spec:
    containers:
    - image: projectraksh/sc-scratch:latest
      imagePullPolicy: IfNotPresent
      name: nginx
      ports:
      - containerPort: 80
        protocol: TCP
      resources: {}
      volumeMounts:
      - mountPath: /etc/raksh
        name: secure-volume-nginx
        readOnly: true
    imagePullSecrets:
    - name: regcred
    volumes:
    - configMap:
        items:
        - key: nginx
          path: raksh.properties
        name: secure-configmap-nginx
      name: secure-volume-nginx
  status: {}
spec:
  SecureContainerImageRef:
    name: nginx-securecontainerimage
status: {}
```

Deploy the secure workload:

```shell
$ kubectl create -f nginx-sc.yaml
```
