#

## Checklist for your GKE

- [X] Have a high level architecture of your applications and map them into a Kubernetes cluster, think about how they're going to deploy and run. 

- [X] Public vs priavte GKE Cluster
    - Control plane (API) endpoint accessbility
    - Public vs. private IPs for nodes
    - Node service accounts
    - Accessibility of node metadata


- [X] Planning networking
    - Plan Pod density per node
    - Planning IP ranges for: nodes, pods, services, control plane, load balancer (maybe), etc.
    - Routes-based or VPC-native
    - Single VPC, shared VPC, VPC Peering, or Open-Hybrid
    - NATs


- [X] Cluster security
    - Restrict access to the control plane
    - Limit access with IAM service accounts
    - Periodically rotate credentials
    - Secret encryption
    - Isolation based on communication needs for multi-tenancy
    - Workload identity
    - RBAC


- [X] Capacity and scalability
    - Capacity of nodes based on your requirement, such as CPU, GPU or TPU, etc.
    - Node Auto Provison (NAP)
    - Node pool autoscaling
    - Horizontal Pod autoscaling (HPA)
    - Vertical Pod autoscaling (VPA)
    - Multidimensional Pod autoscaling
    - Customized cluster autoscaler with community edition


- [X] Storage
    - Zonal or regional persistent storage
    - Shared storage
    - Object storage
    - Require curtain IOPS of storage 


- [X] Operational capability
    - Enroll to a right channel: Rapid, Regular or Stable
    - GKE or GKE Autopilot
    - Zonal or regional cluster
    - Maintainace window
    - Metrics, logging and tracing
    - Backup for GKE
    - Runtime security monitoring - Security Posture
    - CI/CD

- [X] Running services on GKE
    - Using Service Mesh
    - Multi-cluster
    - Service discovery across multiple clusters
    - KubeDNS or Cloud DNS
    - Configuration
    - Ingress or Gateway API

## References
- [Smooth sailing with Kubernetes](https://cloud.google.com/kubernetes-engine/kubernetes-comic)
- [Preparing a Google Kubernetes Engine environment for production ](https://cloud.google.com/architecture/prep-kubernetes-engine-for-prod)