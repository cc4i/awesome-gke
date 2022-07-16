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

### 4.2 Anthos Service Mesh (ASM) + Locality Setting 

> Deploy pods across the nodes in different zones to supports highly available and scalable, as well as leverage weight distrubution of Istio to reduce inter-zone traffic & cost.

```sh
# Clone repo
git clone https://github.com/cc4i/multi-k8s.git
cd multi-k8s
cd hack && ./gke.sh

# Provision & deployment
skaffold build 
tag=`skaffold build --dry-run --output='{{json .}}' --quiet |jq '.builds[].tag' -r`
skaffold deploy --images ${tag}

# Apply locality setting
kubectl apply -f manifests/istio -n run-tracker

# Run
endpoint=`kubectl get svc/istio-ingressgateway -n run-tracker -o "jsonpath={.status.loadBalancer.ingress[0].ip}"`
curl -v http://${endpoint}/trip |jq

```

### 4.3  Multi-Cluster Ingress (MCI) + Anthos Service Mesh (ASM) + Locality Setting

> Using MCI and ASM to implement multi-cluster traffic managment across different Cloud Providers. In each individual k8s cluster we leverage weight distrubution of Istio to reduce inter-zone traffic & cost.