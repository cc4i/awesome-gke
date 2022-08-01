## Run boortrap scripts when launching nodes in GKE

### Instruction

```sh


# Clone repo
git clone https://github.com/cc4i/multi-k8s.git
cd bootstrap/startup-script

# build image & push into registry
docker build . -t <image url>:<tag>
docker push <image url>:<tag>

# Modify script as per request betweend "# BO:" and "# EO:"

# Deploy to cluster
kustomize build . | kubectl apply -f -

```

### Notes
- The pod container can share the host process ID namespace.
- THe container is required root privilege.