# K3s on GCP

## Provision K3s on GCE with Cloud Controller Manager (CCM)

> Key points:
> - Grant server & agent nodes with permission in order to provision network resources.
> - Allow port:60000 in the firewall rules.
> - Following guidance to build CCM and push into your own image registry.
> - Configure RBAC and allow CCM working properly, see [reference](../../manifests/ccm-k3s/).

> Steps
> - Create service account with proper permissions.
> - Create instance template for K3s server.
> - Create managed instance group with server template.
> - Install K3s server side.
> - Create agent instances template for K3s agent.
> - Cerate managed instance group with agent template.
> - Taint server node.
> - Deploy CCM for GCE into K3s cluster.

<br>

## All-in-one quick-start example
The k3s cluster uses managed instance group for both server and agent nodes for high available purpose, but no auto scaling capability, no external database.

```sh

git clone https://github.com/cc4i/multi-k8s.git
cd multi-k8s/hack && ./k3s-quickstart.sh

```

## High-Availability and auto-scalability K3s Cluster
The K3s cluster leverages managed instance group for both server and agent nodes, integrate with internal load balance with server nodes to provide fixed registeration address. Configure auto scaling policy with CPU utilization by default, using Cloud SQL as external database.


```sh
```


## Clean up

```sh

cd multi-k8s/hack && ./k3s-quickstart-cleanup.sh

```