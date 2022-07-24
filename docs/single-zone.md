## Place Pods into nodes in single zone with high availability

### Instruction

```sh


# Clone repo
git clone https://github.com/cc4i/multi-k8s.git
cd multi-k8s
cd hack && ./gke-affinity.sh

# Apply manifests
cd ..
kubectl apply -f manifests/ns.yaml
kubectl apply -f manifests/service-account.yaml
kubectl apply -f deploy-affinity
kubectl apply -f service

```