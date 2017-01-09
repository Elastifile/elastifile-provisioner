# Kubernetes dynamic volume provisioner for Elastifile

## Configuration

Setting a Kubernetes secret with your password:

    echo -n changeme > password.txt
    kubectl create secret generic elastifile-rest --from-file=password.txt
