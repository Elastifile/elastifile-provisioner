# Kubernetes Dynamic Volume Provisioner for Elastifile

`elastifile-provisioner` is an out-of-tree dynamic provisioner for Kubernetes, based on the Kubernetes Incubator's [nfs-provisioner](http://github.com/kubernetes-incubator/nfs-provisioner). It can be used to dynamically provision Kubernetes persistent volumes on an Elastifile ECFS system.

## Configuration

```console
git clone https://github.com/Elastifile/elastifile-provisioner.git
$ cd elastifile-provisioner
```

Customize the Storage Class parameters in `deploy/kube-config/storageclass.yaml` to match your Elastifile setup:



```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1beta1
metadata:
  name: elastifile
provisioner: elastifile.com/nfs
parameters:
  nfsServer: "192.168.0.1"         # Elastifile NFS address
  restURL: "http://10.11.198.128"  # REST API URL
  username: "admin"                # REST API username
  secretName: "elastifile-rest"    # Name of Kubernetes secret storing the REST API password
  secretNamespace: "default"       # Kubernetes namespace for the secret
```
    
Create a Kubernetes secret with your Elastifile password:

```console
$ echo -n "changeme" > password.txt
$ kubectl create secret generic elastifile-rest --from-file=password.txt -n default
secret "elastifile-rest" created
```

Create the Storage Class that you configured above:
```console
$ kubectl create -f deploy/kube-config/storageclass.yaml -n default
storageclass "elastifile" created
```

### RBAC 
```console
$ kubectl create -f deploy/kube-config/serviceaccount.yaml -n default
serviceaccount "elastifile-provisioner" created
$ kubectl create -f deploy/kube-config/clusterrole.yaml -n default
clusterrole "elastifile-provisioner-runner" created
$ kubectl create -f deploy/kube-config/clusterrolebinding.yaml -n default
clusterrolebinding "elastifile-provisioner" created
```

Create the deployment for the provisioner:
```console
$ kubectl create -f deploy/kube-config/deployment.yaml
deployment "elastifile-provisioner" created
```
```console
$ kubectl patch deployment elastifile-provisioner -p '{"spec":{"template":{"spec":{"serviceAccount":"elastifile-provisioner"}}}}'
```

### OpenShift
```console
$ oc create -f deploy/kube-config/serviceaccount.yaml -n default
serviceaccount "elastifile-provisioner" created
$ oc create -f deploy/kube-config/openshift-clusterrole.yaml -n default
clusterrole "elastifile-provisioner-runner" created
$ oadm policy add-scc-to-user hostmount-anyuid system:serviceaccount:default:elastifile-provisioner
$ oadm policy add-cluster-role-to-user nfs-client-provisioner-runner system:serviceaccount:default:nfs-client-provisioner
$ oc patch deployment elastifile-provisioner -p '{"spec":{"template":{"spec":{"serviceAccount":"elastifile-provisioner"}}}}'
```




### Elastifile provisioner validation  
To validate the configuration please Create a `PersistentVolumeClaim` with annotation `volume.beta.kubernetes.io/storage-class: "elastifile"`.
You may want to modify the `storage` size and/or the `accessModes`:

### Example :

```console
$ kubectl create -f deploy/kube-config/pvc.yaml
persistentvolumeclaim "elasti1" created
```

A `PersistentVolume` is provisioned for the `PersistentVolumeClaim`. Now the claim can be consumed by some pod(s) and the backing Elastifile storage read from or written to.
```console



$ kubectl get pvc
NAME      STATUS    VOLUME                                     CAPACITY   ACCESSMODES   AGE
elasti1   Bound     pvc-11819eb8-d66d-11e6-a66b-005056912012   3Gi        RWX           7s

$ kubectl get pv
NAME                                       CAPACITY   ACCESSMODES   RECLAIMPOLICY   STATUS    CLAIM             REASON    AGE
pvc-11819eb8-d66d-11e6-a66b-005056912012   3Gi        RWX           Delete          Bound     default/elasti1             7s
```

Deleting the `PersistentVolumeClaim` will cause the provisioner to delete the `PersistentVolume` and its data.

## Contributing/Developer guide:
#### Prerequisites
[get go](https://golang.org/doc/install)

[get docker](https://docs.docker.com/engine/installation/)

#### Set up the project
To start with, we’ll need the source code for the project. It’s important that it lives inside the GOPATH, so the easiest way to grab the code is with go get: 
```console
go get github.com/elastifile/elastifile-provisioner
```
We’ll do the rest of our work from the project directory:
```console
cd $GOPATH/src/github.com/elastifile/elastifile-provisioner
```

#### Build
To produce a provisioner binary:
```console
$ make build
```
To build the provisioner docker image:
```console
$ make container
```
To push the provisioner docker image:
```console
$ make push
```
