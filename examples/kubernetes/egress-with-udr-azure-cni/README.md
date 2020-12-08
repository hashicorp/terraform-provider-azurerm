## Example: Kubernetes Cluster with Azure CNI and cluster egress using a User-Defined Route

This example creates a Kubernetes Cluster with Azure CNI and cluster egress using a User-Defined Route.

Setting outboundType requires AKS clusters with a vm-set-type of VirtualMachineScaleSets and load-balancer-sku of Standard.

Setting outboundType to a value of UDR requires a user-defined route with valid outbound connectivity for the cluster.

For further information, refer to Azure AKS [documentation](https://docs.microsoft.com/en-us/azure/aks/egress-outboundtype).
