#!/bin/bash


# Environments
export PROJECT_ID=play-with-anthos-340801
export PROJECT_NUMBER=828493099439
export GKE_CLUSTER_AFFINITY=gke-affinity
export GKE_CLUSTER_ISTIO=gke-istio
export REGION=asia-southeast1
export ZONE_A=asia-southeast1-a
export ZONE_B=asia-southeast1-b


# 1.Provision GKE
gcloud container --project ${PROJECT_ID} clusters create ${GKE_CLUSTER_ISTIO} \
    --region ${REGION} \
    --no-enable-basic-auth \
    --machine-type "n2d-standard-2" \
    --image-type "COS_CONTAINERD" \
    --disk-type "pd-standard" \
    --disk-size "100" \
    --metadata disable-legacy-endpoints=true \
    --scopes "https://www.googleapis.com/auth/cloud-platform" \
    --logging=SYSTEM,WORKLOAD --monitoring=SYSTEM \
    --enable-ip-alias \
    --network "projects/${PROJECT_ID}/global/networks/default" \
    --subnetwork "projects/${PROJECT_ID}/regions/${REGION}/subnetworks/default" \
    --enable-intra-node-visibility --default-max-pods-per-node "30" \
    --enable-autoscaling --min-nodes "2" --max-nodes "8" \
    --enable-dataplane-v2 \
    --no-enable-master-authorized-networks \
    --addons HorizontalPodAutoscaling,HttpLoadBalancing,GcePersistentDiskCsiDriver \
    --max-surge-upgrade 1 --max-unavailable-upgrade 0 \
    --labels mesh_id=proj-${PROJECT_NUMBER} \
    --workload-pool "${PROJECT_ID}.svc.id.goog" \
    --enable-shielded-nodes \
    --node-locations ${ZONE_A},${ZONE_B}

kubectl config rename-context gke_${PROJECT_ID}_${REGION}_${GKE_CLUSTER_ISTIO} ${GKE_CLUSTER_ISTIO}

# 2.Install ASM/Istio from !!!CloudShell/Linux!!!
#
echo """
mkdir ~/bin && cd ~/bin
curl https://storage.googleapis.com/csm-artifacts/asm/asmcli_1.14 > asmcli
chmod +x asmcli
./asmcli install \
  --project_id ${PROJECT_ID} \
  --cluster_name ${GKE_CLUSTER_ISTIO} \
  --cluster_location ${REGION} \
  --fleet_id ${PROJECT_ID} \
  --output_dir ~/bin \
  --enable_all \
  --ca mesh_ca
rev=`kubectl get deploy -n istio-system -l app=istiod -o \
  "jsonpath={.items[*].metadata.labels['istio\.io/rev']}{'\n'}"`
sed -e s/REVISION/${rev}/g manifests/base/ns.yaml > manifests/base/ns-rev.yaml
"""
#


