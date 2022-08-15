#!/bin/bash

export SERVICE_ACCOUNT_ID=sa-k3s-nodes
export PROJECT_ID=play-with-anthos-340801
export PROJECT_NUMBER=`gcloud projects list --filter PROJECT_ID=${PROJECT_ID} --format "value(PROJECT_NUMBER)"`
export INSTANCE_TYPE=e2-medium
export NETWROK=default
export REGION=asia-southeast1
export ZONE=asia-southeast1-b	

# Delete Service Account
gcloud iam service-accounts delete ${SERVICE_ACCOUNT_ID}@${PROJECT_ID}.iam.gserviceaccount.com --quiet


# Delete k3s server groups
gcloud compute instance-templates delete k3s-server-instance-group --quiet
gcloud compute instance-groups managed delete k3s-server-instance-group --zone ${ZONE} --quiet


# Delete k3s agent groups 
gcloud compute instance-templates delete k3s-agent-mig-template --quiet
gcloud compute instance-groups managed delete k3s-agent-instance-group --zone ${ZONE} --quiet

# Delete kubeconfig for k3s
rm -f ./k3s-r.yaml
rm -f ./k3s.yaml