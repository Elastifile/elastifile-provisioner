set -x
kubectl delete deployment elastifile-provisioner-1-elastifile-provisioner
kubectl delete job.batch/elastifile-provisioner-1-deployer
kubectl delete job elastifile-provisioner-1-deployer-custom-script
kubectl delete serviceaccount elastifile-provisioner-sa
kubectl delete storageclass elastifile
kubectl delete configmap elastifile-provisioner-1-deployer-config
kubectl delete secret elastifile-rest
#kubectl delete application elastifile-provisioner-1 # hangs
set +x
