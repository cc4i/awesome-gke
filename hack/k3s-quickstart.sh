#!/bin/bash

export SERVICE_ACCOUNT_ID=sa-k3s-nodes
export PROJECT_ID=play-with-anthos-340801
export PROJECT_NUMBER=
export INSTANCE_TYPE=
export NETWROK=default
export REGION=asia-southeast1
export ZONE=asia-southeast1-b	

# 0. Create SA with proper roles
gcloud iam service-accounts describe ${SERVICE_ACCOUNT_ID}@${PROJECT_ID}.iam.gserviceaccount.com
if [ $? != 0 ]
then
    gcloud iam service-accounts create ${SERVICE_ACCOUNT_ID} \
        --description="Service account for K3s nodes" \
        --display-name=${SERVICE_ACCOUNT_ID}
fi


roles=("roles/artifactregistry.reader" \
"roles/cloudprofiler.agent" "roles/cloudtrace.agent" "roles/logging.logWriter" "roles/monitoring.metricWriter" "roles/stackdriver.resourceMetadata.writer" \
"roles/compute.admin" "roles/compute.loadBalancerAdmin" \
"roles/compute.networkAdmin" "roles/compute.orgFirewallPolicyAdmin" "roles/compute.orgSecurityPolicyAdmin" \
"roles/container.admin" \
"roles/storage.admin")

for r in "${roles[@]}"
do
    echo "Binding a role -> ${r} to ${SERVICE_ACCOUNT_ID}"
    gcloud projects add-iam-policy-binding ${PROJECT_ID} \
        --member="serviceAccount:${SERVICE_ACCOUNT_ID}@${PROJECT_ID}.iam.gserviceaccount.com" \
        --role="${r}"
done 

gcloud projects get-iam-policy ${PROJECT_ID} \
    --flatten="bindings[].members" --format='table(bindings.role)' \
    --filter="bindings.members:${SERVICE_ACCOUNT_ID}@${PROJECT_ID}.iam.gserviceaccount.com"



# 1. Create instance template for server
gcloud compute instance-templates create k3s-server-mig-template \
    --project=${PROJECT_ID} --machine-type=${INSTANCE_TYPE} \
    --network-interface=network=${NETWROK},network-tier=PREMIUM \
    --metadata=enable-oslogin=true --maintenance-policy=MIGRATE \
    --provisioning-model=STANDARD \
    --service-account=${SERVICE_ACCOUNT_ID}@${PROJECT_ID}.iam.gserviceaccount.com \
    --scopes=https://www.googleapis.com/auth/cloud-platform \
    --create-disk=auto-delete=yes,boot=yes,device-name=k3s-instance-mig-template,image=projects/debian-cloud/global/images/debian-11-bullseye-v20220719,mode=rw,size=50,type=pd-balanced \
    --no-shielded-secure-boot \
    --shielded-vtpm --shielded-integrity-monitoring --reservation-affinity=any


# 2. Create MIG for server group
gcloud compute instance-groups managed create k3s-server-instance-group \
    --zone ${ZONE} \
    --template k3s-server-mig-template \
    --size 1

k3s_server=`gcloud compute instance-groups managed list-instances a-master-us-central1-a-simple-k8s-local --zone us-central1-a --format="json" |jq -r ".[].instance"`

# 3. Install K3s server 
gcloud compute ssh ${k3s_server} -- "sudo curl -sfL https://get.k3s.io | sh -s - server --disable servicelb --disable-cloud-controller --https-listen-port 443"
k3s_server_ip=`gcloud compute instances describe ${k3s_server} --format="json" |jq -r ".networkInterfaces[].networkIP"`
k3s_server_token=

# 4. Create agent instances template 
startup_script=$(cat << EOF
#! /bin/bash
curl -sfL https://get.k3s.io | K3S_URL=https://${k3s_server_ip}:443 K3S_TOKEN=${k3s_server_token} sh -
EOF
)

gcloud compute instance-templates create k3s-agent-mig-template \
    --project=${PROJECT_ID} --machine-type=${INSTANCE_TYPE} \
    --network-interface=network=${NETWROK},network-tier=PREMIUM \
    --metadata=enable-oslogin=true --maintenance-policy=MIGRATE \
    --provisioning-model=STANDARD \
    --service-account=${SERVICE_ACCOUNT_ID}@${PROJECT_ID}.iam.gserviceaccount.com \
    --scopes=https://www.googleapis.com/auth/cloud-platform \
    --create-disk=auto-delete=yes,boot=yes,device-name=k3s-instance-mig-template,image=projects/debian-cloud/global/images/debian-11-bullseye-v20220719,mode=rw,size=50,type=pd-balanced \
    --no-shielded-secure-boot \
    --shielded-vtpm --shielded-integrity-monitoring --reservation-affinity=any \
    --metadata=startup-script=${startup_script}

# 5. Cerate agent group & register through startup-script
gcloud compute instance-groups managed create k3s-agent-instance-group \
    --zone ${ZONE} \
    --template k3s-agent-mig-template \
    --size 2


# 6. Taint server node
# 7. Deploy CCM for GCE into K3s cluster