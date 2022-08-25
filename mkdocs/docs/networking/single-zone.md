# 
## Description
Placing pods into nodes in single zone with high availability, there're two node pools, one is the primary node pool in signle zone, other one is standby node pool with minmum number of instance is zero.

## Deployment

```sh

# Clone repo
git clone https://github.com/cc4i/multi-k8s.git
cd multi-k8s
cd asset/tod/bin && ./gke-affinity.sh

# Apply manifests
cd ..
kubectl apply -f manifests/ns.yaml
kubectl apply -f manifests/service-account.yaml
kubectl apply -f deploy-affinity
kubectl apply -f service

```

## Notes
- Affinity is a feature from upstream Kubernetes. 