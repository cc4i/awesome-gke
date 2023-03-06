# 

## Description
The issue with managing the core dumps in each pod, in case the Pod/Container/Node restart… then we will be losing all the data which was generated in core-dump folder, those are important data which cannot be compromised. 
By default the container is generating the core dumps in the path “/var/lib/systemd/coredump”, this path can vary as per the core dump pattern defined in your container. We can create the core dump location specified in the file “/proc/sys/kernel/core_pattern” by doing something like below.


```bash
mkdir -p /tmp/cores
chmod a+rwx /tmp/cores
echo “/tmp/cores/core.%e.%p.%h.%t” > /proc/sys/kernel/core_pattern
```

Here I will be describing how we can manage core dumps for a GKE Cluster and storage core dumps into GCS for further process.

## Guide

```bash

# 1. Create service account to access storage bucket.
# 2. Create a Google Storage bucket.
# 3. Persistent volume mount on core dumps location inside application container.
# 4. Build a core dumps agent to collect dump files.
# 5. Deploy core dump agent as a daemonset.
# 6. Run a sample application to generate core dump. 

```

## References 
- [Core Dump handler](https://github.com/IBM/core-dump-handler)
- [How to generate coredump for containers running with k8s](https://shuanglu1993.medium.com/how-to-generate-coredump-for-containers-running-with-k8s-1a3f4a7e75b2)