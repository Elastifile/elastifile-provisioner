FROM gcr.io/google-appengine/debian9 AS build

RUN apt-get update && apt-get install -y --no-install-recommends gettext

# Handle custom arguments
ADD schema.yaml /tmp/schema.yaml

# Provide registry prefix and tag for default values for images.
ARG REGISTRY
ARG TAG
RUN cat /tmp/schema.yaml | env -i "REGISTRY=$REGISTRY" "TAG=$TAG" envsubst > /tmp/schema.yaml.new && mv /tmp/schema.yaml.new /tmp/schema.yaml

#FROM gcr.io/google-marketplace-tools/k8s/deployer_envsubst
FROM gcr.io/cloud-marketplace-tools/k8s/deployer_envsubst

# Install tools for debugging - can be removed later
RUN apt-get update && apt-get install -y less vim screen

COPY config/storageclass.yaml /tmp/storageclass.yaml
COPY config/serviceaccount.yaml /tmp/serviceaccount.yaml
COPY config/clusterrole.yaml /tmp/clusterrole.yaml
COPY config/clusterrolebinding.yaml /tmp/clusterrolebinding.yaml
COPY config/deployment.yaml /tmp/deployment.yaml
#COPY config/pvc.yaml /tmp/pvc.yaml
COPY --from=build /tmp/schema.yaml /data/
COPY deploy.sh /tmp/deploy.sh

#COPY manifest /data/manifest
#COPY apptest/deployer /data-test/

#CMD tail -f /dev/null
#ENTRYPOINT ["tail", "-f", "/dev/null"]
#ENTRYPOINT ["/bin/bash", "/tmp/deploy.sh", "-k"]
ENTRYPOINT ["/bin/bash", "/tmp/deploy.sh"]

