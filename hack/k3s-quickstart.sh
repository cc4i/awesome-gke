#!/bin/bash

export SERVICE_ACCOUNT_ID=sa-k3s-nodes
export PROJECT_ID=play-with-anthos-340801
export PROJECT_NUMBER=`gcloud projects list --filter PROJECT_ID=${PROJECT_ID} --format "value(PROJECT_NUMBER)"`
export INSTANCE_TYPE=e2-medium
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
gcloud compute instance-templates describe k3s-server-mig-template
if [ $? != 0 ]
then 
    gcloud compute instance-templates create k3s-server-mig-template \
        --project=${PROJECT_ID} --machine-type=${INSTANCE_TYPE} \
        --network-interface=network=${NETWROK},network-tier=PREMIUM,address="" \
        --maintenance-policy=MIGRATE \
        --provisioning-model=STANDARD \
        --service-account=${SERVICE_ACCOUNT_ID}@${PROJECT_ID}.iam.gserviceaccount.com \
        --scopes=https://www.googleapis.com/auth/cloud-platform \
        --create-disk=auto-delete=yes,boot=yes,device-name=k3s-instance-mig-template,image=projects/debian-cloud/global/images/debian-11-bullseye-v20220719,mode=rw,size=50,type=pd-balanced \
        --no-shielded-secure-boot \
        --shielded-vtpm --shielded-integrity-monitoring --reservation-affinity=any \
        --metadata=enable-oslogin=true
fi


# 2. Create MIG for server group
gcloud compute instance-groups describe k3s-server-instance-group --zone ${ZONE}
if [ $? != 0 ]
then
    gcloud compute instance-groups managed create k3s-server-instance-group \
        --zone ${ZONE} \
        --template k3s-server-mig-template \
        --size 1
    # Waiting for the managed instance group is ready
    while [ $? != 0 ] 
    do
        gcloud compute instance-groups managed list-instances k3s-server-instance-group --zone ${ZONE} --format="json"
    done
fi

k3s_server=`gcloud compute instance-groups managed list-instances k3s-server-instance-group --zone ${ZONE} --format="json" |jq -r ".[].instance"`
echo "${k3s_server} is ready to provision K3s."

# 3. Install K3s server 
gcloud compute ssh ${k3s_server} -- "sudo curl -sfL https://get.k3s.io | sh -s - server --disable servicelb --disable-cloud-controller --https-listen-port 443"
k3s_server_ip=`gcloud compute instances describe ${k3s_server} --format="json" |jq -r ".networkInterfaces[].networkIP"`
k3s_server_external_ip=`gcloud compute instances describe k3s-server-1 --format="json" --zone asia-southeast1-b |jq -r ".networkInterfaces[].accessConfigs[].natIP"`
k3s_server_token=`gcloud compute ssh k3s-server-1 --zone asia-southeast1-b -- "sudo cat /var/lib/rancher/k3s/server/node-token"`
# Retrieve k3s.yaml for kubectl
gcloud compute ssh ${k3s_server} --zone ${ZONE} -- "sudo cat /etc/rancher/k3s/k3s.yaml">k3s.yaml

# Inject cloud.config file for CCM
cloud_config=$(cat << EOF
[global]\nnode-tags = k3s-cluster-node\nmultizone = true\n
EOF
)
gcloud compute ssh ${k3s_server} --zone ${ZONE} -- "sudo echo ${cloud_config}>./cloud.config;sudo mkdir /etc/kubernetes;sudo cp ./cloud.config /etc/kubernetes/cloud.config"


# 4. Create agent instances template 
startup_script=$(cat << EOF
#! /bin/bash
curl -sfL https://get.k3s.io | K3S_URL=https://${k3s_server_ip}:443 K3S_TOKEN=${k3s_server_token} sh -
EOF
)

gcloud compute instance-templates describe k3s-agent-mig-template
if [ $? != 0 ]
then 
    gcloud compute instance-templates create k3s-agent-mig-template \
        --project=${PROJECT_ID} --machine-type=${INSTANCE_TYPE} \
        --network-interface=network=${NETWROK},network-tier=PREMIUM,address="" \
        --maintenance-policy=MIGRATE \
        --provisioning-model=STANDARD \
        --service-account=${SERVICE_ACCOUNT_ID}@${PROJECT_ID}.iam.gserviceaccount.com \
        --scopes=https://www.googleapis.com/auth/cloud-platform \
        --create-disk=auto-delete=yes,boot=yes,device-name=k3s-instance-mig-template,image=projects/debian-cloud/global/images/debian-11-bullseye-v20220719,mode=rw,size=50,type=pd-balanced \
        --no-shielded-secure-boot \
        --shielded-vtpm --shielded-integrity-monitoring --reservation-affinity=any \
        --metadata=startup-script=${startup_script},enable-oslogin=true
fi

# 5. Cerate agent group & register through startup-script
gcloud compute instance-groups describe k3s-agent-instance-group --zone ${ZONE}
if [ $? != 0 ]
then
    gcloud compute instance-groups managed create k3s-agent-instance-group \
        --zone ${ZONE} \
        --template k3s-agent-mig-template \
        --size 2
    # Waiting for the managed instance group is ready
    while [ $? != 0 ] 
    do
        gcloud compute instance-groups managed list-instances k3s-agent-instance-group --zone ${ZONE} --format="json"
    done
fi


# 6. Taint server node
sed -e 's/127.0.0.1/${k3s_server_external_ip}/g' k3s.yaml > k3s-r.yaml
kubectl --kubeconfig=k3s-r.yaml  get nodes

kubectl --kubeconfig=k3s-r.yaml taint nodes ${k3s_server} node-role.kubernetes.io/control-plane:NoSchedule
kubectl --kubeconfig=k3s-r.yaml get pods -n kube-system


# 7. Deploy CCM for GCE into K3s cluster
# kubectl --kubeconfig=k3s-r.yaml apply -f ../manifests/ccm-k3s/extension-apiserver-authentication.yaml
kubectl --kubeconfig=k3s-r.yaml apply -f ../manifests/ccm-k3s/role.yaml
kubectl --kubeconfig=k3s-r.yaml apply -f ../manifests/ccm-k3s/sa.yaml
kubectl --kubeconfig=k3s-r.yaml apply -f ../manifests/ccm-k3s/rb.yaml
kubectl --kubeconfig=k3s-r.yaml apply -f ../manifests/ccm-k3s/gce.yaml
kubectl --kubeconfig=k3s-r.yaml get pods -n kube-system