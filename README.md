# K8s Clusters on Google Cloud

Create a place to help people better leveraging GKE & Anthos products in actions, includes various demo and tutorials for different scenarios. 

## 1. Infrastructure
### 1.1. [Place Pods into nodes in single zone with high availability](./docs/single-zone.md)
> The cluster has multiple node pools for cross different zones, one zone for primary and one for standby. Using Affinity/Anti-affinity to place Pods into nodes in primary zone and shift to standby zone when there's zonal failure.
<br>

## 2. Autoscaling 
<br>

## 3. Observability
### 3.1. Managed Promestheus with Dataproc on GKE
<br>

## 4. Networking

### 4.1 Repalce Ingress by Gateway API

### 4.2 [Anthos Service Mesh (ASM) + Locality Setting](./docs/asm-locality.md)

> Deploy pods across the nodes in different zones to supports highly available and scalable, as well as leverage weight distrubution of Istio to reduce inter-zone traffic & cost.


### 4.3 [GLB + Anthos Service Mesh (ASM) + Locality Setting](./docs/glb-locality.md)

> Using GLB and ASM to implement multi-cluster traffic managment across different Cloud Providers. In each individual k8s cluster we leverage weight distrubution of Istio to reduce inter-zone traffic & cost.

### 4.4 Multi-cluster mesh outside Google Cloud
> https://cloud.google.com/service-mesh/docs/unified-install/off-gcp-multi-cluster-setup