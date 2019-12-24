# Overview
A given Kubernetes YAML needs to be first converted so as to enable deployment leveraging secure containers.
This is a multi-step process automated via `rakshctl` CLI and is described here:

## Convert the YAML to secure YAML
Using the `rakshctl` CLI we'll convert the following YAML to enable it to be deployed using secure containers
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

The command used will be:
```
$ rakshctl image create  -i nginx-securecontainerimage \
             --initrd kata-containers-initrd.img  \
             --symmKey  symm_key \
             --vmlinux vmlinuz \
             --scratch sc-scratch:latest  
             --filename nginx.yaml nginx-securecontainerimage
```

### Create ConfigMap of the container spec
The container spec looks like this:
```
spec:
  containers:
    - image: nginx:latest
      imagePullPolicy: IfNotPresent
      name: nginx
      ports:
      - containerPort: 80
        protocol: TCP
```
The ConfigMap of the above spec looks something like this:
```
apiVersion: v1
data:
  nginx: |
    spec:
      containers:
      - image: nginx:latest
        imagePullPolicy: IfNotPresent
        name: nginx
        ports:
        - containerPort: 80
          protocol: TCP
kind: ConfigMap
metadata:
  name: configmap-nginx
```

### Encrypt the ConfigMap with symmetric key 
The encrypted ConfigMap data looks like this:
```
apiVersion: v1
data:
  nginx: 6qvygg8md7bXfyX3Y9cpZxUp4eZA0kKmWBirrpJv/WEGkrdLYrdtqxdqm4cGLG4++06d2iGTaB+5SDjjDwf05T+9a2iUAdHmRngHcQNAzkKK2RCnR4Zkt0cXDaEP+w5mbugH0xdqGm8SoX4IgvWGi2toq1CUcc8OmgTX42g0NruTZbrNv5NccyS7+kR7Iib6vaMI24E=
kind: ConfigMap
metadata:
   name: secure-configmap-nginx
```
### Create the modified YAML
The final YAML to be used to deploy the application using secure containers looks like this:
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


The above YAML can be deployed as any other Kubernetes YAML 

