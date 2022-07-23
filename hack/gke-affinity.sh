#!/bin/bash


# Environments
export PROJECT_ID=play-with-anthos-340801
export PROJECT_NUMBER=828493099439
export GKE_CLUSTER=gke-affinity
export REGION=asia-southeast1
export ZONE_A=asia-southeast1-a
export ZONE_B=asia-southeast1-b
export ZONE_C=asia-southeast1-c
export INSTANCE_TYPE=e2-medium


# 1.Provision GKE
gcloud container --project ${PROJECT_ID} clusters create ${GKE_CLUSTER} \
    --region ${REGION} \
    --no-enable-basic-auth \
    --metadata disable-legacy-endpoints=true \
    --scopes "https://www.googleapis.com/auth/cloud-platform" \
    --logging=SYSTEM,WORKLOAD --monitoring=SYSTEM \
    --enable-ip-alias \
    --network "projects/${PROJECT_ID}/global/networks/default" \
    --subnetwork "projects/${PROJECT_ID}/regions/${REGION}/subnetworks/default" \
    --enable-intra-node-visibility --default-max-pods-per-node "30" \
    --max-pods-per-node "30" --num-nodes "1" \
    --enable-dataplane-v2 \
    --no-enable-master-authorized-networks \
    --addons HorizontalPodAutoscaling,HttpLoadBalancing,GcePersistentDiskCsiDriver \
    --max-surge-upgrade 1 --max-unavailable-upgrade 0 \
    --labels mesh_id=proj-${PROJECT_NUMBER} \
    --workload-pool "${PROJECT_ID}.svc.id.goog" \
    --enable-shielded-nodes \
    --node-locations ${ZONE_A}

# 2.Added node pools each zone
gcloud beta container --project "${PROJECT_ID}" node-pools create "np-${INSTANCE_TYPE}-${ZONE_B}" \
    --cluster "${GKE_CLUSTER}" \
    --region "${REGION}" \
    --machine-type "${INSTANCE_TYPE}" --image-type "COS_CONTAINERD" --disk-type "pd-standard" --disk-size "100" \
    --metadata disable-legacy-endpoints=true --scopes "https://www.googleapis.com/auth/cloud-platform" \
    --num-nodes "2" \
    --enable-autoupgrade --enable-autorepair \
    --max-surge-upgrade 1 --max-unavailable-upgrade 0 \
    --max-pods-per-node "30" \
    --node-locations "${ZONE_B}"


gcloud beta container --project "${PROJECT_ID}" node-pools create "np-${INSTANCE_TYPE}-${ZONE_C}" \
    --cluster "${GKE_CLUSTER}" \
    --region "${REGION}" --machine-type "${INSTANCE_TYPE}" \
    --image-type "COS_CONTAINERD" --disk-type "pd-standard" --disk-size "100" \
    --metadata disable-legacy-endpoints=true --scopes "https://www.googleapis.com/auth/cloud-platform" \
    --enable-autoscaling --min-nodes "0" --num-nodes "0" --max-nodes "4" \
    --enable-autoupgrade --enable-autorepair \
    --max-surge-upgrade 1 --max-unavailable-upgrade 0 \
    --max-pods-per-node "30" \
    --node-locations ${ZONE_C}

kubectl config rename-context gke_${PROJECT_ID}_${REGION}_${GKE_CLUSTER} ${GKE_CLUSTER}



