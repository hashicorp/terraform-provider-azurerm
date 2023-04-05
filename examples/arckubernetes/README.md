## Example: Azure Arc Kubernetes

This example provisions the following Resources:

## Creates

1. A [Resource Group](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/resource_group)
2. A [Virtual Network](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/virtual_network)
3. A [Subnet](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/subnet)
4. A [Linux Virtual Machine](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/linux_virtual_machine)
5. An [Arc Kubernetes Cluster](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/arc_kubernetes_cluster)
6. A [Kind Cluster](https://kind.sigs.k8s.io/) in [Linux Virtual Machine](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/linux_virtual_machine) by [remote-exec Provisioner](https://developer.hashicorp.com/terraform/language/resources/provisioners/remote-exec)
7. [Azure Arc agents](https://learn.microsoft.com/en-us/azure/azure-arc/kubernetes/conceptual-agent-overview) in [Kind Cluster](https://kind.sigs.k8s.io/) by [remote-exec Provisioner](https://developer.hashicorp.com/terraform/language/resources/provisioners/remote-exec)

~> **NOTE:** To connect an existing Kubernetes cluster to Azure Arc, the following conditions must be met:

* An [Arc Kubernetes Cluster](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/arc_kubernetes_cluster) must be created in Azure
* [Azure Arc agents](https://learn.microsoft.com/en-us/azure/azure-arc/kubernetes/conceptual-agent-overview) must be installed in the Kubernetes cluster which is connected to Azure

## Usage

- Provide values to all variables
- Create with `terraform apply`
- Destroy all with `terraform destroy`
