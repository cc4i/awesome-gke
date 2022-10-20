#!/bin/sh

# Setup environment variables

export PROJECT_ID=play-api-service
export PROJECT_NUMBER=`gcloud projects list --filter PROJECT_ID=${PROJECT_ID} --format "value(PROJECT_NUMBER)"`
export REGION=asia-southeast1
export ZONE_A=asia-southeast1-a
export ZONE_B=asia-southeast1-b
export CLUSTER=nap-cluster

# Provision cluster with NAP auto-scaling
gcloud beta container --project "play-api-service" clusters create "nap-cluster" \
    --region "asia-southeast1" \
    --no-enable-basic-auth --cluster-version "1.23.8-gke.1900" --release-channel "regular" \
    --machine-type "e2-medium" \
    --image-type "COS_CONTAINERD" \
    --disk-type "pd-standard" --disk-size "100" --metadata disable-legacy-endpoints=true \
    --scopes "https://www.googleapis.com/auth/cloud-platform" \
    --max-pods-per-node "30" --num-nodes "1" \
    --logging=SYSTEM,WORKLOAD --monitoring=SYSTEM \
    --enable-ip-alias \
    --network "projects/play-api-service/global/networks/default" \
    --subnetwork "projects/play-api-service/regions/asia-southeast1/subnetworks/default" \
    --no-enable-intra-node-visibility \
    --default-max-pods-per-node "30" \
    --no-enable-master-authorized-networks \
    --addons HorizontalPodAutoscaling,HttpLoadBalancing,GcePersistentDiskCsiDriver,BackupRestore,GcpFilestoreCsiDriver \
    --enable-autoupgrade \
    --enable-autorepair \
    --max-surge-upgrade 1 --max-unavailable-upgrade 0 \
    --enable-autoprovisioning --min-cpu 1 --max-cpu 120 --min-memory 1 --max-memory 480 \
    --autoprovisioning-scopes=https://www.googleapis.com/auth/cloud-platform \
    --enable-autoprovisioning-autorepair \
    --enable-autoprovisioning-autoupgrade \
    --autoprovisioning-max-surge-upgrade 1 \
    --autoprovisioning-max-unavailable-upgrade 0 \
    --autoscaling-profile optimize-utilization \
    --workload-pool "play-api-service.svc.id.goog" \
    --enable-shielded-nodes \
    --node-locations "asia-southeast1-a","asia-southeast1-b" \
    --logging-variant=MAX_THROUGHPUT

# Install Anthos Service Mesh with fleet API
gcloud container clusters update  \
    --project play-api-service nap-cluster \
    --region asia-southeast1 --update-labels mesh_id=proj-${PROJECT_NUMBER}

gcloud container fleet mesh update \
    --management automatic \
    --memberships nap-cluster \
    --project play-api-service

# Enable auto injection
# kubectl label namespace NAMESPACE istio-injection=enabled istio.io/rev- --overwrite