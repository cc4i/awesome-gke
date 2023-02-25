#

## Description

You need to manually tranfer your PV across zone if you don't use reginal persistent disk at very begining, here's guide to help you transfer you storage.

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