# Overview
The project Consists of 
- Kubernetes Operator and CRDs
- Custom Kata agent
- CLI tool (`rakshctli`)

# Operator
Secure container is implemented using the Operator pattern. Following are the Custom Resource Definitions (CRDs)

## SecureContainerImageConfig 
- SVM (Kata VM) kernel and initrd location
- Kata Runtime Class to use

## SecureContainerImage 
- Container image with SVM (Kata VM) kernel and initrd

The image gets downloaded and extracted to location specified by SecureContainerImageConfig 

## SecureContainer
- Secure Container resource to deploy app container encapsulated in secure VM (SVM)


# Custom Kata Agent 
This handles container image lifecycle - download, extract, decrypt, execute inside the SVM 

# CLI 
`rakshctl` CLI is to convert an existing  Kubernetes application specification  (YAML) to support secure containers


# Workflow

Convert App YAML to secure APP YAML -> Deploy the secure app YAML to Kubernetes Cluster 

The following high level steps explains the conversion of the applciation YAML 

1. Read the container specification from the YAML and create a ConfigMap for the spec
2. Encrypt the ConfigMap using symmetric key and a random unique string (nonce)
3. Store the symmetric key, nonce in Kata initrd and modify it for SVM (lockbox creation, encrypt rootfs etc) . 
    * When using it for development encrypted rootfs is not needed. Additionally the key can be stored and retrieved from vault
4. Create container image consisting of Kata kernel and initrd. This will be the SecureContainerImage resource
5. Output modified YAML. This will use the SecureContainer resource

## Sample - Original YAML
```
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

## Encrypted ConfigMap of the original YAML
```
apiVersion: v1
data:
  nginx: 6qvygg8md7bXfyX3Y9cpZxUp4eZA0kKmWBirrpJv/WEGkrdLYrdtqxdqm4cGLG4++06d2iGTaB+5SDjjDwf05T+9a2iUAdHmRngHcQNAzkKK2RCnR4Zkt0cXDaEP+w5mbugH0xdqGm8SoX4IgvWGi2toq1CUcc8OmgTX42g0NruTZbrNv5NccyS7+kR7Iib6vaMI24E=
kind: ConfigMap
metadata:
   name: secure-configmap-nginx
```

## Sample - Modified (secure app) YAML
```
apiVersion: securecontainers.k8s.io/v1alpha1
kind: SecureContainer
metadata:
  name: secure-nginx
object:
  apiVersion: v1
  kind: Pod
  metadata:
    labels:
      app: nginx
    name: nginx
spec:
    containers:
    - image: sc-scratch:latest
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
    volumes:
    - configMap:
        items:
        - key: nginx
          path: raksh.properties
        name: secure-configmap-nginx
      name: secure-volume-nginx
spec:
  SecureContainerImageRef:
    name: nginx-securecontainerimage
```


