# Kubernetes Dynamic Volume Provisioner for Elastifile

`elastifile-provisioner` is an out-of-tree dynamic provisioner for Kubernetes, based on the Kubernetes Incubator's [nfs-provisioner](http://github.com/kubernetes-incubator/nfs-provisioner). It can be used to dynamically provision Kubernetes persistent volumes on an Elastifile ECFS system.

## Configuration

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
$ kubectl create secret generic elastifile-rest --from-file=password.txt
secret "elastifile-rest" created
```

Create the Storage Class that you configured above:
```console
$ kubectl create -f deploy/kube-config/storageclass.yaml
storageclass "elastifile" created
```

Create the deployment for the provisioner:
```console
$ kubectl create -f deploy/kube-config/deployment.yaml
deployment "elastifile-provisioner" created
```

Create a `PersistentVolumeClaim` with annotation `volume.beta.kubernetes.io/storage-class: "elastifile"`.
You may want to modify the `storage` size and/or the `accessModes`:
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
