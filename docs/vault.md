# Vault Integration with raksh

This document describes how the Vault server can be used to manage and utilize the required keys used by raksh

## Installation

Please follow [this](https://www.vaultproject.io/docs/platform/k8s/run.html) documentation to learn how to run `Vault` on Kubernetes.

Once you have deployed vault in your Kubernetes cluster by following the above documentation you'll need to get the vault service details to be used later on.


```sh
$ kubectl get svc
NAME                 TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)                   AGE
vault                 ClusterIP    10.96.38.205    <none>        8200/TCP,8201/TCP         14h
```
From the above,  the vault service is available at the following location `http://10.96.38.205:8200`

If vault is running on the same cluster, you can exec to the vault container and run the following command to create a vault token.

```
$ kubectl exec -it vault-0 sh
$ vault token create
s.YNAqtnnRiEuKy1y6uEBSavvG
```

## Required Key(s)

To conceal the workload details, an encrypted configMap is attached to the pod yaml which contains the details of the actual workload. This configMap is encrypted with AES 256 symmetric key.

The key needs to be uploaded the `Vault` server by the user before it can be used by `raksh`. If you already have the required symmetric key in your `Vault` server feel free to skip to section [Kubernetes Secret](##Kubernetes-Secret)

### Generating and Uploading the Key to the Vault

One of the ways to generate this symmetric key is,

```sh
$ dd bs=1 if=/dev/random of=/tmp/symm_key count=32
```
or
```sh
$ openssl rand -rand /dev/urandom 32 > /tmp/symm_key
```

This key, located at `/tmp/symm_key`, needs to be uploaded to the vault server.

Let's begin by generating a `base64` string of this key.

```sh
$ base64 /tmp/symm_key
tfi/lC30L0JsgT0RVFYi+p9+EWsge9IVQqqb+euYhW4=
```

Now let's upload the `base64` encoded key to the vault server.

```sh
$ vault kv put secret/hello symm_key=tfi/lC30L0JsgT0RVFYi+p9+EWsge9IVQqqb+euYhW4=
```
This writes the `symm_key` to the path `secret/hello`. Learn more about creating vault secrets [here](https://www.vaultproject.io/intro/getting-started/first-secret.html).

## Kubernetes Secret

raksh agent needs to fetch the symmetric key to decrypt the configMap from the Vault server. We do not want to add vault access credentials directly to the workload's deployment yaml. Doing so will increase the risk of unauthorized access to sensitive vault credentials which can result in unauthorized access to the symmetric key stored in the vault server.

We will create a Kubernetes Secret to store the Vault realted information. All key values in a Kubernetes Secret need to be base64 encoded.

1. Vault IP address and port - This is required to let the kata agent know where we are hosting our Vault server. If Vault is running on the Kubernetes server and exposed as a service available at http://10.96.38.205:8200 as mentioned before, To generate `base64` encoded string,

```sh
$  echo -n "http://10.96.38.205:8200" | base64 -w 0 ; echo ""
aHR0cDovLzEwLjk2LjM4LjIwNTo4MjAw
```
2. Vault token - A token is used to authenticate requests coming from vault clients. Just like above, let's convert the token into it's `base64` encoded string.
```sh
$ echo -n "s.YNAqtnnRiEuKy1y6uEBSavvG" | base64 -w 0 ; echo ""
cy5ZTkFxdG5uUmlFdUt5MXk2dUVCU2F2dkc=
```
3. Secret Name - The path used by the vault server to fetch the secret. If you had pushed the secret to vault server using path `secret/hello` then your secret name would be `secret/data/hello`. Please refer Vault [API documentation](https://www.vaultproject.io/api/secret/kv/kv-v2.html#read-secret-version) for more details.
```sh
$ echo -n "secret/data/hello" | base64 -w 0 ; echo ""
c2VjcmV0L2RhdGEvaGVsbG8=
```
4. Key Name - This is the name of the key stored at the vault path specified above.
```sh
$ echo -n "symm_key" | base64 -w 0 ; echo ""
c3ltbV9rZXk=
```

Now we will take all these vaules and embed them into a Kubernetes Secret,

```yaml
apiVersion: v1
kind: Secret
metadata:
 name: mysecret
type: Opaque
data:
 vaultAdd: aHR0cDovLzEwLjk2LjM4LjIwNTo4MjAw
 vaultToken: cy5ZTkFxdG5uUmlFdUt5MXk2dUVCU2F2dkc=
 secretName: c2VjcmV0L2RhdGEvaGVsbG8=
 keyName: c3ltbV9rZXk=
```
Save this file, say `vault.yaml`, to create a secret using `kubectl create -f vault.yaml` command.

## Modifying the Deployment or Pod yaml

Once we have the kubernetes secret in place, we are ready to modify our deployment or pod yaml so that our runtime agent can read the vault access credentials securely.

```sh
rakshctl image create --image nginx-securecontainerimage --initrd /usr/share/kata-containers/kata-containers-initrd.img --vmlinux /usr/share/kata-containers/vmlinux.container --symmKeyFile /root/key_file --filename /securecontainers/sample/nginx.yaml --scratch <image-registry>/sc-image:latest --push --vaultSecret mysecret <image-registry>/nginx-securecontainerimage:latest
```
Where,

* _nginx.yaml_ is the workload pod yaml
* _&lt;image-registry&gt;/nginx-securecontainerimage:latest_ is the name of the secure image that holds modified kata initrd
* _mysecret_ is the name of the kubernetes secret that hold vault access credentials

For more info on the arguments please refer to the main [README](../README.md)
