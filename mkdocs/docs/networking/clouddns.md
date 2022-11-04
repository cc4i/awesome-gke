# 

## Description 

To improve domain resolving performance and experience in GKE, we can integrate Cloud DNS with GKE instead of using KubeDNS.

## Guide

```sh

# 1. Setup environment variables, replace with your own

export PROJECT_ID=<Project ID>
export ZONE=<Zone>
export REGION=<Region>
export CLUSTER=clouddns-cluster

# 2. Provision GKE with cluster scope DNS

gcloud beta container clusters create clouddns-cluster \
    --cluster-dns=clouddns \
    --cluster-dns-scope=cluster \
    --zone=${ZONE}


# 3. Verify Cloud DNS

kubectl run dns-test --image us-docker.pkg.dev/google-samples/containers/gke/hello-app:2.0
kubectl expose pod dns-test --name dns-test-svc --port 8080
kubectl get svc dns-test-svc
kubectl exec -it dns-test -- nslookup dns-test-svc

# 4. Add A record inside of Cloud DNS, such as "server-kafka-0.server-kafka-headless.default.svc : 192.0.2.99", then you can use follwing command to validate 

kubectl exec -it dns-test -- nslookup server-kafka-0.server-kafka-headless.default.svc.cluster.local

```

## References

- https://medium.com/google-cloud/dns-on-gke-everything-you-need-to-know-b961303f9153
- https://cloud.google.com/kubernetes-engine/docs/how-to/cloud-dns