#!/bin/sh

# Setup environment variables

export PROJECT_ID=
export ZONE=asia-southeast1-b
export REGION=
export CLUSTER=clouddns-cluster

# Provision GKE with cluster scope DNS
#   - https://cloud.google.com/kubernetes-engine/docs/how-to/cloud-dns#cluster_scope_dns
gcloud beta container clusters create clouddns-cluster \
    --cluster-dns=clouddns \
    --cluster-dns-scope=cluster \
    --zone=asia-southeast1-b


# Verify Cloud DNS
# - https://cloud.google.com/kubernetes-engine/docs/how-to/cloud-dns#verify

kubectl run dns-test --image us-docker.pkg.dev/google-samples/containers/gke/hello-app:2.0
kubectl expose pod dns-test --name dns-test-svc --port 8080
kubectl get svc dns-test-svc
kubectl exec -it dns-test -- nslookup dns-test-svc