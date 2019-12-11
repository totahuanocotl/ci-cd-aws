#!/bin/bash -eu
DIR=$(dirname ${BASH_SOURCE[0]})

{% set i = inventory.parameters %}
kubeconfig=${KUBECONFIG:-/tmp/kubeconfig}
echo "Running for target cluster {{ i.target_name }}"
kubectl --kubeconfig ${kubeconfig} -n {{ i.hello_world.namespace }} apply -f ../manifests/pre-deploy
kubectl --kubeconfig ${kubeconfig} -n {{ i.hello_world.namespace }} apply -f ../manifests
