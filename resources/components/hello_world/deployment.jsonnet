local kube = import "lib/kube.libjsonnet";
local kap = import "lib/kapitan.libjsonnet";
local inventory = kap.inventory();


local helloWorldContainer = kube.Container("hello-world") {
    image: inventory.parameters.hello_world.image,
    args_+: {
        debug: inventory.parameters.hello_world.debug,
        port: inventory.parameters.hello_world.port,
    },
    ports_+: {
        http: {containerPort: inventory.parameters.hello_world.port}
    },
};

local deployment = kube.Deployment("hello-world") {
    spec+: {
        replicas: 1,
        template+: {
            spec+: {
                containers_+: {
                    "hello-world": helloWorldContainer
                },
            }
        }
    }
};

local service = kube.Service("hello-world") {
 target_pod: deployment.spec.template
};

{
    "hello-world-deployment": deployment,
    "hello-world-service": service,

}
