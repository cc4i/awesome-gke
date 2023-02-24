# 

## Decription

Somehow you want to delete a specific node from the node pool in GKE cluster, but you can't do that in the console or through 'gcloud container'. Normally you probably cordon and drain the node, and then delete the instance as a normal VM, howerver you can do that through a single command from managed instance group. The node pool in GKE is related to a specific managed instance group, that why we can leverage command like 'gcloud compute instance-groups managed delete-instances'.

## Guide

```sh

# Provison a simple zonal cluster with 3 nodes
gcloude containers create <Cluster Name> \
    --zone <Zone>

# List MIG & nodes
gcloud container node-pools list --cluster <Cluster Name> \
    --zone <Zone>

# Delete a specific node through 'gcloud compute instance-groups managed delete-instances'
gcloud compute instance-groups managed delete-instances <Managed Instance Group> \
    --instances=<Name of Instance1, Name of Instance2, ...> \
    --zone <Zone>

```


## Clean up

```sh

# Delete the GKE cluster
gcloud container clusters delete <Name of Cluster> \
    --zone <Zone>
    --async

```

## References
- https://pminkov.github.io/blog/removing-a-node-from-a-kubernetes-cluster-on-gke-google-container-engine.html