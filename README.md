<img src="./docs/images/skipper.png" alt="multi-k8s" width="200"/>

[GKE](#gke) | [K3s](#k3s) | Hybrid Cloud | Config Management | Security | Operation | Development
<br/>

# Kubernetes Cluster on Google Cloud

Create a place to help people better leveraging GKE & Anthos products in actions, includes various demo and tutorials for different scenarios. 


## GKE
### 1. Infrastructure
#### 1.1. [Place Pods into nodes in single zone with high availability](./docs/single-zone.md)
> The cluster has multiple node pools for cross different zones, one zone for primary and one for standby. Using Affinity/Anti-affinity to place Pods into nodes in primary zone and shift to standby zone when there's zonal failure.

#### 1.2 [Run boortrap scripts when launching nodes in GKE](./docs/startup-script.md)
> To run bootstrap scripts for your nodes in GKE such as initialize something, add iptable entry, etc., you can run a quick DeamonSet to achieve that.

#### 1.3 Using Cloud DNS instead of Kube DNS
> Use much more reliable and robust option Cloud DNS (100% SLO) ot replace Kube DNS. Refernce from [blog](https://medium.com/google-cloud/dns-on-gke-everything-you-need-to-know-b961303f9153). 

#### 1.4 Validating GKE clusters against configuration best practices
> [GKE Policy Automation](https://github.com/google/gke-policy-automation) from Google, contains the tool and the policy library for validating GKE clusters against configuration best practices.
<br>

### 2. Autoscaling 
#### 2.1 Scale the cluster with customized Autoscaler 
<br>

### 3. Observability
#### 3.1. Managed Promestheus with Dataproc on GKE
<br>

### 4. Networking

#### 4.1 Repalce Ingress by Gateway API

#### 4.2 [Anthos Service Mesh (ASM) + Locality Setting](./docs/asm-locality.md)

> Visualize pods across the nodes in different zones to supports highly available and scalable, as well as leverage weight distrubution of Istio to reduce inter-zone traffic & cost.


#### 4.3 [GLB + Anthos Service Mesh (ASM) + Locality Setting](./docs/glb-locality.md)

> Using GLB and ASM to implement multi-cluster traffic managment across different Cloud Providers. In each individual k8s cluster we leverage weight distrubution of Istio to reduce inter-zone traffic & cost.

#### 4.4 Multi-cluster mesh outside Google Cloud
> https://cloud.google.com/service-mesh/docs/unified-install/off-gcp-multi-cluster-setup

### 5.Storage

#### 5.1 Using Google Cloud Storage

> https://github.com/GoogleCloudPlatform/gcsfuse

#### 5.2 Using Google Filestore

> https://github.com/kubernetes-sigs/gcp-filestore-csi-driver

#### 5.3 Using Google Persisten Disk 

> https://github.com/kubernetes-sigs/gcp-compute-persistent-disk-csi-driver

## K3s

### 1. Infrastructure
#### 1.1 Provision K3s on GCE with CCM

> K3s is a lightweight Kubernetes distributuion and packaged as a signle binary, is good for IoT, Edge compute, DevOps, customized k8s compute, etc. Just jump into [here](./blueprint/k3s/play-k3s.md) to know about better leveraging K3s on GCP.