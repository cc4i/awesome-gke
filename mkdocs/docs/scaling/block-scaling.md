#

## Description

Sometimes you may face the issue like "Pod is blocking scale down because it has local storage" in GKE, or doesn't have enough enough Pod Disruption Budget, those issues could be solved in following ways:

- Adding the parameter "cluster-autoscaler.kubernetes.io/safe-to-evict": "true" to your Pods

- Specifying a [Pod Disruption Budget](https://kubernetes.io/docs/tasks/run-application/configure-pdb/) for your Pods

- Check out Node-level reason messages for noScaleDown events appear in the noDecisionStatus.noScaleDown.nodes[].reason field and take action accordingly.

## Refernces

- [NoScaleDown node-level reasons](https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-autoscaler-visibility#noscaledown-node-level-reasons)

- [How does scale-down work?](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#how-does-scale-down-work)

- [Does Cluster autoscaler work with PodDisruptionBudget in scale-down?](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#does-ca-work-with-poddisruptionbudget-in-scale-down)

- [What types of Pods can prevent Cluster autoscaler from removing a node?](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#what-types-of-pods-can-prevent-ca-from-removing-a-node)