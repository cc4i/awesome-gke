#

## Description

In GKE Autopilot [the ephemeral storage limit](https://cloud.google.com/kubernetes-engine/docs/concepts/autopilot-resource-requests#min-max-requests) must be within 10 MiB and 10 GiB for all compute classes, it's a potential issue for those applicaitons need big size storage. 

There are three diffrent way to surpass the limit and you can choose one of those based on your scenario.

## Use memory for emptyDir

## Use Cloud Filestore

## Use Persistent Disk


## Allowed Storage type 

- "configMap"
- "csi"
- "downwardAPI"
- "emptyDir"
- "gcePersistentDisk"
- "hostPath",
- "nfs"
- "persistentVolumeClaim"
- "projected"
- "secret"

## Refernces

- https://cloud.google.com/kubernetes-engine/docs/concepts/autopilot-overview