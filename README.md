# Multiple K8s Clusters


## ASM + Locality Setting 

```sh
# Clone repo
git clone https://github.com/cc4i/multi-k8s.git
cd multi-k8s
cd hack && ./gke.sh

# Provision & deployment
skaffold build 
tag=`skaffold build --dry-run --output='{{json .}}' --quiet |jq '.builds[].tag' -r`
skaffold deploy --images ${tag}

# Run
kubectl port-forward svc/svc-1 8000:8000 -n run-tracker
curl -v http://localhost:8000/trip |jq

```