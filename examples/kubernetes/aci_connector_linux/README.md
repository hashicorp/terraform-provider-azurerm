## Example: AKS Cluster with Virtual Nodes (ACI)

This example is used to demonstrate how to provision an AKS cluster that uses virtual nodes (Azure Container Instances).

This example provisions:
- A Virtual Network with two subnets:
  - One subnet for the AKS default node pool
  - One subnet for the ACI node
- A Service Principal to be used as the AKS cluster identity
- Assigns a *Network Contributor* role to the Service Principal over the subnet where ACI nodes will be provisioned
- AKS cluster with a user managed service principal and a *aci_connector_linux* *addon_profile*

Note that the ACI subnet must be configured with a service delegations:
```
service_delegation {
    name    = "Microsoft.ContainerInstance/containerGroups"
    actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
}
```

A sample app is provided. To deploy the sample, download the AKS credentials:
```
az aks get-credentials --resource-group myResourceGroup --name myAKSCluster
```
Deploy the sample app:
```
kubectl apply -f virtual-node.yaml
```

Further references:
- [Create and configure an Azure Kubernetes Services (AKS) cluster to use virtual nodes using the Azure CLI](https://docs.microsoft.com/en-us/azure/aks/virtual-nodes-cli)
### Variables

- `prefix` - (Required) A prefix used for all resources in this example
- `password` - (Required) Service Principal password
- `location` - (Required) Azure Region in which all resources in this example should be provisioned