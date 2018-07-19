
PROJECT=blah
for i in {0..99}; do ID=$i PROJECT=${PROJECT} envsubst < templates/test-scenario.yaml.template | kubectl create -f -; done
