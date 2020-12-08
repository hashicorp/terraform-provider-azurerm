## Example: Kubernetes Cluster with Kubenet networking and cluster egress using a User-Defined Route

This example creates a Kubernetes Cluster with Kubenet networking and cluster egress using a User-Defined Route.

Setting outboundType requires AKS clusters with a vm-set-type of VirtualMachineScaleSets and load-balancer-sku of Standard.

Setting outboundType to a value of UDR requires a user-defined route with valid outbound connectivity for the cluster.

For further information, refer to Azure AKS [documentation](https://docs.microsoft.com/en-us/azure/aks/egress-outboundtype).

Managed identities are not currently supported with custom route tables in kubenet. (See Azure AKS [documentation](https://docs.microsoft.com/en-us/azure/aks/configure-kubenet#bring-your-own-subnet-and-route-table-with-kubenet))

The service principal for your AKS cluster requires Contributor permissions on the virtual network and route tables.
