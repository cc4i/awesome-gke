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

# gcloud beta container clusters update clouddns-cluster \
#     --cluster-dns=clouddns \
#     --cluster-dns-scope=cluster \
#     --zone=asia-southeast1-b
# gcloud beta container clusters upgrade clouddns-cluster \
#     --node-pool=default-pool \
#     --cluster-version=1.22.12-gke.2300 \
#     --zone=asia-southeast1-b

# Verify Cloud DNS
# - https://cloud.google.com/kubernetes-engine/docs/how-to/cloud-dns#verify

kubectl run dns-test --image us-docker.pkg.dev/google-samples/containers/gke/hello-app:2.0
kubectl expose pod dns-test --name dns-test-svc --port 8080
kubectl get svc dns-test-svc
kubectl exec -it dns-test -- nslookup dns-test-svc

kubectl apply -f - << __EOF__
apiVersion: v1
kind: Pod
metadata:
  name: dnsutils
  namespace: default
spec:
  containers:
  - name: dnsutils
    image: registry.k8s.io/e2e-test-images/jessie-dnsutils:1.3
    command:
      - sleep
      - "infinity"
    imagePullPolicy: IfNotPresent
  restartPolicy: Always
__EOF__

kubectl exec -it dnsutils -- bash

# Add A record inside of Cloud DNS, such as "server-kafka-0.server-kafka-headless.default.svc : 192.0.2.99", then you can use follwing command to validate 

kubectl exec -it dns-test -- nslookup server-kafka-0.server-kafka-headless.default.svc.cluster.local