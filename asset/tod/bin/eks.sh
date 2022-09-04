#!/bin/sh 

export EKS_CLUSTER=eks-sculpture
export REGION=ap-southeast-1

# 1. Create an EKS cluster
eksctl create cluster --name ${EKS_CLUSTER} --region ${REGION}
node_group=`eksctl get nodegroups --cluster ${EKS_CLUSTER} -o json | jq -r ".[].Name"`
eksctl scale nodegroup --name ${node_group} --cluster ${EKS_CLUSTER} -N 6 -M 8

# Create an EKS cluster 
aws eks describe-cluster ${EKS_CLUSTER}

# Get OIDC endpoint
oidc_issuer=`aws eks describe-cluster --name ${EKS_CLUSTER} \
 --region ${REGION} \
 --query "cluster.identity.oidc.issuer" \
 --output text`



# 2. Atached the EKS cluster to GCP Fleet
gcloud container fleet memberships register ${EKS_CLUSTER} \
  --context=chuancc-code@${EKS_CLUSTER}.${REGION}.eksctl.io \
  --kubeconfig=~/.kube/config \
  --enable-workload-identity \
  --public-issuer-url=${oidc_issuer}


# 3. In order to login into EKS cluster in GCP Console, we need to create a bear token for that
# References::
#   - https://cloud.google.com/anthos/identity/setup/bearer-token-auth

ksa_name=gcp-sa
kubectl create serviceaccount ${ksa_name}
kubectl create clusterrolebinding gcp-sa-view-role-binding \
   --clusterrole view --serviceaccount default:${ksa_name}
kubectl create clusterrolebinding gcp-sa-read-role-binding \
   --clusterrole cloud-console-reader --serviceaccount default:${ksa_name}
kubectl create clusterrolebinding gcp-sa-admin-role-binding \
   --clusterrole cluster-admin --serviceaccount default:gcp-sa


secret_name=gcp-sa-token

kubectl apply -f - << __EOF__
apiVersion: v1
kind: Secret
metadata:
  name: "${secret_name}"
  annotations:
    kubernetes.io/service-account.name: "${ksa_name}"
type: kubernetes.io/service-account-token
__EOF__

until [[ $(kubectl get -o=jsonpath="{.data.token}" "secret/${secret_name}") ]]; do
  echo "waiting for token..." >&2;
  sleep 1;
done

bear_token=`kubectl get secret ${secret_name} -o jsonpath='{$.data.token}' | base64 --decode`

echo "The following token is to login into the EKS cluster ::"
echo ${bear_token}
