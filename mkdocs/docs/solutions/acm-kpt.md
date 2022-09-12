#

## Description
The tutorial is to compose the [Kubernetes Resource Model (KRM) blueprints](https://cloud.google.com/anthos-config-management/docs/concepts/blueprints#krm-blueprints) with Anthos Config Management and provision a Google Kubernetes Engine (GKE) cluster and the required networking infrastructure such as a Virtual Private Cloud (VPC) and a subnet to host the GKE cluster, and named IP ranges for pods and services. 

## Prerequisite

- gcloud
- [kpt](https://kpt.dev/)
- kubectl

## Steps

```sh

export PROJECT_ID=play-with-anthos-340801
export ACM_CONTROLLER=acm-controller
export REGION=us-central1
export CONFIG_NAMESPACE=config-control
export VPC_NAME=gke-test-vpc
export SUBNET_NAME=${VPC_NAME}-subnetwork-${REGION}
export GKE_CLUSTER=acm-kpt-cluster

# Enable APIs
gcloud services enable serviceusage.googleapis.com \
    krmapihosting.googleapis.com \
    container.googleapis.com \
    cloudresourcemanager.googleapis.com

# Create Anthos Config Controller, GKE cluster could be private with NAT, more detail - https://cloud.google.com/sdk/gcloud/reference/anthos/config/controller/create
gcloud anthos config controller create ${ACM_CONTROLLER} \
    --location=${REGION}

# Give Config Controller permission to manage Google Cloud resources
export SA_EMAIL="$(kubectl get ConfigConnectorContext -n config-control \
    -o jsonpath='{.items[0].spec.googleServiceAccount}' 2> /dev/null)"
gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member "serviceAccount:${SA_EMAIL}" \
    --role "roles/owner" \
    --project "${PROJECT_ID}"


# Verify that Config Connector is configured and healthy in the project namespace
kubectl get ConfigConnectorContext -n ${CONFIG_NAMESPACE} \
    -o "custom-columns=NAMESPACE:.metadata.namespace,NAME:.metadata.name,HEALTHY:.status.healthy"

# Configure the VPC
kpt pkg get \
  https://github.com/cc4i/multi-k8s.git/asset/acm-kpt/vpc@main \
  ${VPC_NAME}
cd ${VPC_NAME}
( echo "cat <<EOF" ; cat setters.yaml ; echo EOF ) | sh > setters-val.yaml
mv setters-val.yaml setters.yaml
# kpt fn eval --image set-namespace:v0.1 -- namespace=config-control


# Configure subnet 
kpt pkg get \
  https://github.com/cc4i/multi-k8s.git/asset/acm-kpt/subnet@main \
  ${SUBNET_NAME}
cd ${SUBNET_NAME}
( echo "cat <<EOF" ; cat setters.yaml ; echo EOF ) | sh > setters-val.yaml
mv setters-val.yaml setters.yaml
# kpt fn eval --image set-namespace:v0.1 -- namespace=config-control


# Initialize the working directory with kpt, which creates a resource to track changes
cd ..
kpt fn render
kpt live init --namespace config-control
kpt live apply --dry-run
kpt live apply
kpt live status --output table --poll-until current

# Configure GKE
cd ..
kpt pkg get \
  https://github.com/cc4i/multi-k8s.git/asset/acm-kpt/gke@main \
  ${GKE_CLUSTER}
cd ${GKE_CLUSTER}
( echo "cat <<EOF" ; cat setters.yaml ; echo EOF ) | sh > setters-val.yaml
mv setters-val.yaml setters.yaml
kpt fn render
kpt live init --namespace config-control
kpt live apply --dry-run
kpt live apply
kpt live status --output table --poll-until current

```

## Clean Up

```sh
cd ${GKE_CLUSTER}
kpt live destroy

cd ${VPC_NAME}
kpt live destroy

```

## References 

- https://cloud.google.com/anthos-config-management/docs/concepts/blueprints#krm-blueprints
- https://cloud.google.com/anthos-config-management/docs/tutorials/gke-cluster-blueprint 
- https://github.com/GoogleCloudPlatform/blueprints/tree/main/catalog