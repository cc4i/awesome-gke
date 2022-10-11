
#!/bin/bash


# Enviroment variables
export PROJECT_ID=play-with-anthos-340801
export REGION=asia-southeast1
export CLUSTER=auto-k8s
export PROJECT_NUMBER=`gcloud projects list --filter PROJECT_ID=${PROJECT_ID} --format "value(PROJECT_NUMBER)"`

# Provision GKE Autopilot
gcloud container clusters create-auto ${auto-k8s} \
    --release-channel=rapid \
    --scopes=https://www.googleapis.com/auth/cloud-platform \
    --region=${REGION}


# Install managed Anthos Service Mesh with fleet API 
#   - oud.google.com/service-mesh/docs/managed/automatic-management-with-fleet
#   - Waiting for a while after you just creating your Autopilot (!!!avoid the faiure casued by not ready cluster!!!)
gcloud container fleet memberships register ${CLUSTER} \
  --gke-uri=https://container.googleapis.com/v1/projects/${PROJECT_ID}/locations/${REGION}/clusters/${CLUSTER} \
  --enable-workload-identity \
  --project ${PROJECT_ID}

gcloud container fleet memberships list --project ${PROJECT_ID}

# !!! Only for cluster's project differs from your fleet host project!!!
# gcloud projects add-iam-policy-binding "${PROJECT_ID}"  \
#   --member "serviceAccount:service-FLEET_PROJECT_NUMBER@gcp-sa-servicemesh.iam.gserviceaccount.com" \
#   --role roles/anthosservicemesh.serviceAgent

gcloud container clusters update  --project ${PROJECT_ID} ${CLUSTER} \
  --region ${REGION} --update-labels mesh_id=proj-${PROJECT_NUMBER}

gcloud container fleet mesh update \
    --management automatic \
    --memberships ${CLUSTER} \
    --project ${PROJECT_ID}

# Enable managed data plane to fully manage upgrades of the proxies 
#   - The managed data plane is enabled for the automatically provisioned managed control plane revision.
#   - https://cloud.google.com/service-mesh/docs/managed/configure-managed-anthos-service-mesh#managed-data-plane

# Deploy demo for validation
#   - kubectl label namespace book istio-injection=enabled istio.io/rev- --overwrite
#
kubectl create ns book
kubectl label namespace book istio-injection=enabled istio.io/rev- --overwrite
curl -LO https://storage.googleapis.com/gke-release/asm/istio-1.14.4-asm.2-osx.tar.gz
tar xzvf istio-1.14.4-asm.2-osx.tar.gz
cd istio-1.14.4-asm.2
kubectl apply -f samples/bookinfo/platform/kube/bookinfo.yaml -n book
# setup istio-gateway manually (https://cloud.google.com/service-mesh/docs/gateways#deploy_gateways)
kubectl apply -f samples/bookinfo/networking/bookinfo-gateway.yaml -n book


# Using Spot instances for deployment
#   - Spot Pods during preemption is 25 seconds
#   - https://cloud.google.com/kubernetes-engine/docs/how-to/autopilot-spot-pods
