#

## Description

If you use StorageClass to provision regional PV in GKE, your Pod can mount PV back even it's reboot in node at different zona. However when you use zonal persistent disk to provision your PV (by default), and then you need to manually tranfer your storage across zone. Here's guide to help you transfer your zonal PV in GKE. 

## Guide


### 1. Take a snapshot for your orginal volumn.
```bash

```

### 2. Create a new disk using the snapthot.
```bash
```

### 3. Create a new PV and PVC from the disk.
```bash
```

### 4. Launch a example deployment using PVC created above.
```bash
```

## Clean up

## References

- [Create clones of persistent volumes](https://cloud.google.com/kubernetes-engine/docs/how-to/persistent-volumes/volume-cloning)
- [Create a snapshot of a zonal persistent disk](https://cloud.google.com/compute/docs/disks/create-snapshots#create_zonal_snapshot)
- [Persistent volumes and dynamic provisioning](https://cloud.google.com/kubernetes-engine/docs/concepts/persistent-volumes)