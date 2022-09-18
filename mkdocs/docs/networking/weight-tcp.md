# 

## Description

Control TCP L4 communication is very important in Kubernetes world, such as migrating TCP traffic gradually from an older version of a microservice to a new one. In GKE, you can accomplish that by leveraging Anthos Service Mesh and shfiting a percentage of TCP traffic from one destination to another. 

## Guide

```sh
# Clone repo
git clone https://github.com/cc4i/multi-k8s.git
cd multi-k8s/asset/tod/bin && ./gke.sh


# Following pinting to install ASM for GKE, https://cloud.google.com/service-mesh/docs/unified-install/install-anthos-service-mesh#install_mesh_ca


#  Build application 
skaffold build 
tag=`skaffold build --dry-run --output='{{json .}}' --quiet |jq '.builds[].tag' -r`


# Deploy demo applications
cd multi-k8s/asset/tod/manifests/examples/tcp
kustomize build . | kubectl apply -f -

# Modify ../../weight-tcp/serving.yaml as by comments

# Validate through 'telnet' 
telnet <ingress-gateway> 8008
# Modify weight and observe the status of current connection & new connection.

```


## Operate with External TCP Proxy Load Balancer

We can use [External TCP Proxy Load Balancer](https://cloud.google.com/load-balancing/docs/forwarding-rule-concepts#tcp_proxy) with NEG to achieve that, but it's much tedious tasks.

Create NEG(Network Endpoint Group) with annotations, like following code snippet:

```yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    cloud.google.com/neg: '{"exposed_ports": {"8008":{"name": "svc-1-v2-tcp-neg"}}}'

```
And then create backend service with those NEGs, follow by forwarding rule with TCP protocol. Whenever we need to deploy new version, we can deployment new version with different labels, create new NEGs. Manipulate those NEGs, remove older one and add new one. We can leverage Connection Draining Timeout (maxmium 3600 seconds) to wait flight connection to complete and do not accept new requests during a connection drain.

>There's potential connection risk to disconnect with old service during update configuration of backend. Change was slow compare to weight-based technique.


## References

- TCP Traffic Shifting by Istio
