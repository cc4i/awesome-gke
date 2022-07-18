# K8s Clusters on Google Cloud

Create a place to help people better leveraging GKE & Anthos products in actions, includes various demo and tutorials for different scenarios. 

## 1. Infrastructure
### 1.1. Place Pods into nodes in same zone with high availability
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