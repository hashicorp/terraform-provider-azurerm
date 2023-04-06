---
subcategory: "ArcKubernetes"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_arc_kubernetes_cluster"
description: |-
  Manages an Arc Kubernetes Cluster.
---

# azurerm_arc_kubernetes_cluster

Manages an Arc Kubernetes Cluster.

-> **Note:** Installing and configuring the Azure Arc Agent on your Kubernetes Cluster to establish connectivity is outside the scope of this document. For more details refer to [Deploy agents to your cluster](https://learn.microsoft.com/en-us/azure/azure-arc/kubernetes/conceptual-agent-overview#deploy-agents-to-your-cluster) and [Connect an existing Kubernetes Cluster](https://learn.microsoft.com/en-us/azure/azure-arc/kubernetes/quickstart-connect-cluster?tabs=azure-cli#connect-an-existing-kubernetes-cluster). If you encounter issues connecting your Kubernetes Cluster to Azure Arc, we'd recommend opening a ticket with Microsoft Support.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_arc_kubernetes_cluster" "example" {
  name                         = "example-akcc"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = "West Europe"
  agent_public_key_certificate = filebase64("testdata/public.cer")

  identity {
    type = "SystemAssigned"
  }

  tags = {
    ENV = "Test"
  }
}
```

-> **Note:** An extensive example on connecting the `azurerm_arc_kubernetes_cluster` to an external kubernetes cluster can be found in [the `./examples/arckubernetes` directory within the GitHub Repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/arckubernetes).

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Arc Kubernetes Cluster. Changing this forces a new Arc Kubernetes Cluster to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Arc Kubernetes Cluster should exist. Changing this forces a new Arc Kubernetes Cluster to be created.

* `agent_public_key_certificate` - (Required) Specifies the base64-encoded public certificate used by the agent to do the initial handshake to the backend services in Azure. Changing this forces a new Arc Kubernetes Cluster to be created.

* `identity` - (Required) An `identity` block as defined below. Changing this forces a new Arc Kubernetes Cluster to be created.

* `location` - (Required) Specifies the Azure Region where the Arc Kubernetes Cluster should exist. Changing this forces a new Arc Kubernetes Cluster to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Arc Kubernetes Cluster.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity assigned to this Arc Kubernetes Cluster. At this time the only possible value is `SystemAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Arc Kubernetes Cluster.

* `agent_version` - Version of the agent running on the cluster resource.

* `distribution` - The distribution running on this Arc Kubernetes Cluster.

* `identity` - An `identity` block as defined below.

* `infrastructure` - The infrastructure on which the Arc Kubernetes Cluster is running on.

* `kubernetes_version` - The Kubernetes version of the cluster resource.

* `offering` - The cluster offering.

* `total_core_count` - Number of CPU cores present in the cluster resource.

* `total_node_count` - Number of nodes present in the cluster resource.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Arc Kubernetes Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Arc Kubernetes Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Arc Kubernetes Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Arc Kubernetes Cluster.

## Import

Arc Kubernetes Cluster can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_arc_kubernetes_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Kubernetes/connectedClusters/cluster1
```
