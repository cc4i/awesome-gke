## Anthos Service Mesh (ASM) + Locality Setting 

> Deploy pods across the nodes in different zones to supports highly available and scalable, as well as leverage weight distrubution of Istio to reduce inter-zone traffic & cost.

### Instruction

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