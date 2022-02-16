---
subcategory: "HybridKubernetes"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hybrid_kubernetes_connected_cluster"
description: |-
  Manages a HybridKubernetes connected cluster.
---

# azurerm_hybrid_kubernetes_connected_cluster

Manages a HybridKubernetes connected cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_hybrid_kubernetes_connected_cluster" "test" {
  name                         = "example-connected-cluster"
  resource_group_name          = azurerm_resource_group.test.name
  agent_public_key_certificate = "xxxxxxxxxxxxxxxxxxx"
  distribution                 = "gke"

  identity {
    type = "SystemAssigned"
  }
  infrastructure = "gcp"
  location       = "West Europe"

  tags = {
    ENV = "Test"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this HybridKubernetes Connected Cluster. Changing this forces a new HybridKubernetes Connected Cluster to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the HybridKubernetes connected cluster should exist. Changing this forces a new HybridKubernetes connected cluster to be created.

* `agent_public_key_certificate` - (Required) Base64 encoded public certificate used by the agent to do the initial handshake to the backend services in Azure. Changing this forces a new HybridKubernetes connected cluster to be created.

* `identity` - (Required) An `identity` block as defined below.

* `location` - (Required) The Azure Region where the HybridKubernetes connected cluster should exist. Changing this forces a new HybridKubernetes connected cluster to be created.

* `distribution` - (Optional) The HybridKubernetes distribution running on this connected cluster. Changing this forces a new HybridKubernetes connected cluster to be created.

* `infrastructure` - (Optional) The infrastructure on which the HybridKubernetes cluster represented by this connected cluster is running on. Changing this forces a new HybridKubernetes connected cluster to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the HybridKubernetes connected cluster.

---

An `identity` block exports the following:

* `type` - (Required) Specifies the type of Managed Service Identity. Possible values are `SystemAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the HybridKubernetes connected cluster.

* `identity` - An `identity` block as defined below.

* `provisioning_state` - The current deployment state of connected clusters.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the HybridKubernetes connected cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the HybridKubernetes connected cluster.
* `update` - (Defaults to 30 minutes) Used when updating the HybridKubernetes connected cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the HybridKubernetes connected cluster.

## Import

HybridKubernetes connected cluster can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_hybrid_kubernetes_connected_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Kubernetes/connectedClusters/cluster1
```
