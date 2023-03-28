# Gateway API Hands-on Labs

## Overview
The sample labs in this section demonstrate how to use the Gateway API to configure a Gateway and expose a service to external traffic. The labs are based on the [Gateway API](https://gateway-api.sigs.k8s.io/) and [GKE Gateway](https://cloud.google.com/kubernetes-engine/docs/how-to/deploying-gateways) features.

## Contents
- Lab1. [External Gateway using Certificate Manager with wildcard domain]()
- Lab2. [Configuring a static IP for a Gateway]()
- Lab3. [Multi-cluster gateway]()
- Lab4. [Capacity-based load balancing]()


## Prerequisites
### 1. Enable APIs, such as GKE, GKE Hub, and GKE Multi-cluster Ingress (MCI), etc. in your projects

```shell
# Environment variables
export PROJECT_ID=play-api-service
export PROJECT_NUMBER=374299782509

# Enable APIs
gcloud services enable \
    container.googleapis.com \
    gkehub.googleapis.com \
    multiclusterservicediscovery.googleapis.com \
    multiclusteringress.googleapis.com \
    trafficdirector.googleapis.com \
    --project=${PROJECT_ID}
```

### 2. Provision a GKE cluster with Gateway API, or enable existed GKE clutser

```shell
# Create a GKE cluster with Gateway API
gcloud container clusters create pt-cluster-4 \
    --gateway-api=standard \
    --release-channel=Rapid \
    --region=asia-southeast1


# Enable Gateway API for an existing GKE cluster
gcloud container clusters update pt-cluster-4 \
    --gateway-api=standard \
    --region=asia-southeast1

# Verify the cluster
kubectl get gatewayclass
```



## Lab 1. External Gateway using Certificate Manager with wildcard domain

### 1.1 Create a certificate map

```shell
# Enable Certificate Manager API
gcloud services enable certificatemanager.googleapis.com

# Create a DNS Authorization to validate your certificate. Example here is my personal domain cc4i.xyz and reigstered in GoDaddy, you can use your own domain.
gcloud certificate-manager dns-authorizations create dns-auth-cc4i-xyz \
  --domain="cc4i.xyz"

# !!!Returning DNS resource record needs to be added as CNAME into your DNS configuration !!!

# Create a certificate
gcloud beta certificate-manager certificates create store-cc4i-xyz-cert \
    --domains="cc4i.xyz,*.cc4i.xyz" \
    --dns-authorizations=dns-auth-cc4i-xyz

# Create a certificate map
gcloud beta certificate-manager maps create store-cc4i-xyz-map

# Create a certificate map entry for wildcard domain
gcloud beta certificate-manager maps entries create store-cc4i-xyz-map-entry1 \
    --map=store-cc4i-xyz-map \
    --hostname="*.cc4i.xyz" \
    --certificates=store-cc4i-xyz-cert

# Create a certificate map entry for root domain
gcloud beta certificate-manager maps entries create store-cc4i-xyz-map-entry2 \
    --map=store-cc4i-xyz-map \
    --hostname="cc4i.xyz" \
    --certificates=store-cc4i-xyz-cert
```

>References: 
>- [DNS Authorization](https://cloud.google.com/certificate-manager/docs/deploy-google-managed-dns-auth#create_a_google-managed_certificate_referencing_the_dns_authorization)
>- [Secure a gateway](https://cloud.google.com/kubernetes-engine/docs/how-to/secure-gateway)


### 1.2 Create a Gateway

```shell
# Create a Gateway
kubectl apply -f single-https/gateway.yaml

# Verify the Gateway
kubectl describe gateways.gateway.networking.k8s.io external-http
```

### 1.3 Deployment demo application

```shell
# Deploy the demo application
kubectl apply -f single-https/store.yaml
kubectl get pods
kubectl get service

# Check Gateway IP
kubectl get gateways.gateway.networking.k8s.io external-http -o=jsonpath="{.status.addresses[0].value}"

# Apply the http route
kubectl apply -f single-https/store-route-external.yaml
```

> References:
>- [Demo applications](https://cloud.google.com/kubernetes-engine/docs/how-to/deploying-gateways#deploy_the_demo_applications_2)
>- [HTTPRoute](https://gateway-api.sigs.k8s.io/api-types/httproute/)


### 1.4 Use shared Gateways

```shell
# Used as a shared Gateway
kubectl apply -f single-https/site.yaml
# Apply the http route
kubectl apply -f single-https/site-route-external.yaml
```

# Lab 2. Configuring a static IP for a Gateway

## 2.1 Create a static IP address

```shell
# https://cloud.google.com/kubernetes-engine/docs/how-to/deploying-gateways#gateway_ip_addressing


gcloud compute addresses create test-public-ip-cc4i-xyz \
    --global \
    --project=${PROJECT_ID}
```
### 2.2 Create a Gateway with static IP

```shell
# Apply the Gateway with static IP
kubectl apply -f single-https/named-ip-gateway.yaml 
```

## Lab 3. Multi-cluster gateway


## 3.1 Craete multiple GKE clusters in different regions

```shell
# Create multiple GKE clusters in different regions
gcloud container clusters create gke-west-1 \
    --gateway-api=standard \
    --zone=us-west1-a \
    --workload-pool=${PROJECT_ID}.svc.id.goog \
    --cluster-version=1.25.6-gke.1000 \
    --project=${PROJECT_ID}

gcloud container clusters create gke-west-2 \
    --gateway-api=standard \
    --zone=us-west1-a \
    --workload-pool=${PROJECT_ID}.svc.id.goog \
    --cluster-version=1.25.6-gke.1000 \
    --project=${PROJECT_ID}

gcloud container clusters create gke-east-1 \
    --gateway-api=standard \
    --zone=us-east1-b \
    --workload-pool=${PROJECT_ID}.svc.id.goog \
    --cluster-version=1.25.6-gke.1000 \
    --project=${PROJECT_ID}

# Rename the context
kubectl config rename-context gke_${PROJECT_ID}_us-west1-a_gke-west-1 gke-west-1
kubectl config rename-context gke_${PROJECT_ID}_us-west1-a_gke-west-2 gke-west-2
kubectl config rename-context gke_${PROJECT_ID}_us-east1-b_gke-east-1 gke-east-1
```

## 3.2 Register to the fleet

```shell
# Register to the fleet
gcloud container fleet memberships register gke-west-1 \
     --gke-cluster us-west1-a/gke-west-1 \
     --enable-workload-identity \
     --project=${PROJECT_ID}

gcloud container fleet memberships register gke-west-2 \
     --gke-cluster us-west1-a/gke-west-2 \
     --enable-workload-identity \
     --project=${PROJECT_ID}

gcloud container fleet memberships register gke-east-1 \
     --gke-cluster us-east1-b/gke-east-1 \
     --enable-workload-identity \
     --project=${PROJECT_ID}

# Try following command to add permission to service account if you got error ralated to permission!!!
gcloud projects add-iam-policy-binding ${PROJECT_ID} --member \
    serviceAccount:service-${PROJECT_NUMBER}@gcp-sa-gkehub.iam.gserviceaccount.com \
    --role "roles/gkehub.connect"

# Validate the membership
gcloud container fleet memberships list --project=${PROJECT_ID}

```


## 3.3 Enable multi-cluster services

```shell
# Enable multi-cluster services API
gcloud container fleet multi-cluster-services enable \
    --project ${PROJECT_ID}

# Grant IAM permission required for multi-cluster services
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
     --member "serviceAccount:${PROJECT_ID}.svc.id.goog[gke-mcs/gke-mcs-importer]" \
     --role "roles/compute.networkViewer" \
     --project=${PROJECT_ID}

# Choose member cluster, which would host the resource for multi-cluster Gateway
gcloud container fleet ingress enable \
    --config-membership=gke-west-1 \
    --project=${PROJECT_ID}

# Validate the ingress
gcloud container fleet ingress describe --project=${PROJECT_ID}

# Grant IAM permission required for multi-cluster services
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member "serviceAccount:service-${PROJECT_NUMBER}@gcp-sa-multiclusteringress.iam.gserviceaccount.com" \
    --role "roles/container.admin" \
    --project=${PROJECT_ID}

# Validate the GatewayClass
kubectl get gatewayclasses --context=gke-west-1
```

## 3.4 Deploy demo application into all clusters 

```shell
# Deploy demo application into all clusters
kubectl apply --context gke-west-1 -f multi/multi-store.yaml
kubectl apply --context gke-west-2 -f multi/multi-store.yaml
kubectl apply --context gke-east-1 -f multi/multi-store.yaml

kubectl apply -f store-west-service.yaml --context gke-west-1
kubectl apply -f store-east-service.yaml --context gke-east-1 --namespace store

# Validate the deployment
kubectl get serviceexports --context gke-east-1 -n store
kubectl get serviceimports --context gke-east-1 -n store
```

## 3.5 Deploy external http Gateway

```shell
# Deploy external http Gateway
kubectl apply -f multi/external-http-gateway.yaml --context gke-west-1 --namespace store

# Validate the Gateway
kubectl describe gateways.gateway.networking.k8s.io external-http --context gke-west-1 --namespace store
```

## 3.6 Validate deployment
```shell
# Get the IP of the Gateway
kubectl get gateway -n store
kubectl get gateways.gateway.networking.k8s.io external-http -o=jsonpath="{.status.addresses[0].value}" --context gke-west-1 --namespace store

# Take a while to sync all routes into GLB

curl -H "host: store.example.com" http://34.98.80.105
curl -H "host: store.example.com" http://34.98.80.105/west
curl -H "host: store.example.com" http://34.98.80.105/east
```

> Reference:
> - [Enable Multi-cluster Gateway](https://cloud.google.com/kubernetes-engine/docs/how-to/enabling-multi-cluster-gateways)
> - [Fleet](https://cloud.google.com/kubernetes-engine/docs/fleets-overview#introducing_fleets)


## Lab 4. Capacity-based load balancing

## 4.1 Prepare clusters if not done in previous steps
```shell
kubectl get gatewayclasses --context=gke-west-1
```

## 4.2 Deploy the demo application

```shell
# Deploy the application into all clusters
kubectl apply -f capacity-based-glb/store-traffic-deploy.yaml --context gke-west-1
kubectl apply -f capacity-based-glb/store-traffic-deploy.yaml  --context gke-east-1
kubectl apply -f capacity-based-glb/store-service.yaml --context gke-west-1
kubectl apply -f capacity-based-glb store-service.yaml --context gke-east-1
```

## 4.3 Deploy external http Gateway

```shell
# Deploy external http Gateway
kubectl apply -f store-route.yaml --context gke-west-1
# Validate the Gateway
kubectl describe gateways.gateway.networking.k8s.io store -n traffic-test --context gke-west-1
```

## 4.4 Testing the Gateway

```shell
# Get the IP of the Gateway
curl http://34.160.103.232

# Take a while to sync all routes into GLB
kubectl get gateways.gateway.networking.k8s.io store -n traffic-test --context=gke-west-1 -o=jsonpath="{.metadata.annotations.networking\.gke\.io/url-maps}"


# Add a load generator to the cluster
kubectl run --context=gke-west-1 -i --tty --rm loadgen  \
    --image=cyrilbkr/httperf  \
    --restart=Never  \
    -- /bin/sh -c 'httperf  \
    --server=34.160.103.232  \
    --hog --uri="/zone" --port 80  --wsess=100000,1,1 --rate 10'


# Fetch metrics from Monitoring & observe the traffic distribution
fetch https_lb_rule
| metric 'loadbalancing.googleapis.com/https/backend_request_count'
| filter (resource.url_map_name =='gkemcg1-traffic-test-store-armvfyupay1t')
| align rate(1m)
| every 1m
| group_by [resource.backend_scope],
    [value_backend_request_count_aggregate:
       aggregate(value.backend_request_count)]


# Add more traffic to cluster & observe the traffic distribution
kubectl run --context=gke-west-1 -i --tty --rm loadgen  \
    --image=cyrilbkr/httperf  \
    --restart=Never  \
    -- /bin/sh -c 'httperf  \
    --server=34.160.103.232  \
    --hog --uri="/zone" --port 80  --wsess=100000,1,1 --rate 30'
```

> Reference:
> - [Capacity-based load balancing](https://cloud.google.com/kubernetes-engine/docs/how-to/deploying-multi-cluster-gateways#capacity-load-balancing)
