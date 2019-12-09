local kube = import "lib/kube.libjsonnet";
local kap = import "lib/kapitan.libjsonnet";
local inventory = kap.inventory();

{
    "00_namespace": kube.Namespace(inventory.parameters.hello_world.namespace),
}