## Example: Kubernetes Cluster with private API server

This example creates a Kubernetes Cluster with private API server.

In a private cluster, the control plane or API server has internal IP addresses that are defined in the RFC1918 - Address Allocation for Private Internets document. By using a private cluster, you can ensure that network traffic between your API server and your node pools remains on the private network only.

The control plane or API server is in an Azure Kubernetes Service (AKS)-managed Azure subscription. A customer's cluster or node pool is in the customer's subscription. The server and the cluster or node pool can communicate with each other through the Azure Private Link service in the API server virtual network and a private endpoint that's exposed in the subnet of the customer's AKS cluster.

The Private Link service is supported on Standard Azure Load Balancer only. Basic Azure Load Balancer isn't supported.

For further information, refer to Azure AKS [documentation](https://docs.microsoft.com/en-us/azure/aks/private-clusters).