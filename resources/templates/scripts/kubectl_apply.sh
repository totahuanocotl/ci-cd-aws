#!/bin/bash -eu

{% set i = inventory.parameters %}
DIR=$(dirname ${BASH_SOURCE[0]})

echo "Running for target cluster {{ i.target_name }}"
echo ${DIR}
kubectl --kubeconfig ./kubeconfig --context {{ i.cluster.context }} -n {{ i.hello_world.namespace }} apply -f ../manifests/pre-deploy
kubectl --kubeconfig ./kubeconfig --context {{ i.cluster.context }} -n {{ i.hello_world.namespace }} apply -f ../manifests
