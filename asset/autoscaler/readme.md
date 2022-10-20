

helm install autoscaler autoscaler/cluster-autoscaler \
    --set "autoscalingGroupsnamePrefix[0].name=default-,autoscalingGroupsnamePrefix[0].maxSize=10,autoscalingGroupsnamePrefix[0].minSize=1" \
    --set autoDiscovery.clusterName=test-autoscaler-1 \
    --set cloudProvider=gce

gcloud container clusters update custom-autoscaler-1 --zone asia-east2-a --autoscaling-profile=profile-unspecified

gcloud container clusters get-credentials test-autoscaler-3 --zone us-central1-c



gcloud container clusters create "test-autoscaler-2" --scopes "https://www.googleapis.com/auth/cloud-platform"



helm install autoscaler2 autoscaler/cluster-autoscaler \
    --set "autoscalingGroupsnamePrefix[0].name=default-,autoscalingGroupsnamePrefix[0].maxSize=10,autoscalingGroupsnamePrefix[0].minSize=1" \
    --set autoDiscovery.clusterName=test-autoscaler-2 \
    --set cloudProvider=gce


helm install autoscaler2 autoscaler/cluster-autoscaler \
    --set "autoscalingGroupsnamePrefix[0].name=default-,autoscalingGroupsnamePrefix[0].maxSize=10,autoscalingGroupsnamePrefix[0].minSize=1" \
    --set autoDiscovery.clusterName=test-autoscaler-3 \
    --set cloudProvider=gce



helm install --dry-run autoscaler2 autoscaler/cluster-autoscaler \
    --set "autoscalingGroupsnamePrefix[0].name=default-,autoscalingGroupsnamePrefix[0].maxSize=10,autoscalingGroupsnamePrefix[0].minSize=1" \
    --set "autoscalingGroupsnamePrefix[1].name=default2-,autoscalingGroupsnamePrefix[1].maxSize=10,autoscalingGroupsnamePrefix[1].minSize=1" \
    --set autoDiscovery.clusterName=test-autoscaler-1 \
    --set cloudProvider=gce