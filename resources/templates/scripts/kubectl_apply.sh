#!/bin/bash -eu

{% set i = inventory.parameters %}
DIR=$(dirname ${BASH_SOURCE[0]})

echo "Running for target cluster {{ i.target_name }}"
echo ${DIR}
kubectl --context {{ i.cluster.context }} -n {{ i.hello_world.namespace }} apply -f ../manifests/pre-deploy
kubectl --context {{ i.cluster.context }} -n {{ i.hello_world.namespace }} apply -f ../manifests
