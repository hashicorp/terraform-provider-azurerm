---
subcategory: "Kubernetes"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_connected_cluster"
description: |-
  Manages a kubernetes connected cluster.
---

# azurerm_kubernetes_connected_cluster

Manages a kubernetes connected cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-kubernetes"
  location = "West Europe"
}

resource "azurerm_kubernetes_connected_cluster" "test" {
  name                         = "acctest-k-%d"
  resource_group_name          = azurerm_resource_group.test.name
  agent_public_key_certificate = "xxxxxxxxxxxxxxxxxxx"
  distribution                 = "gke"

  identity = {
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

* `name` - (Required) The name which should be used for this Kubernetes Connected Cluster. Changing this forces a new Kubernetes Connected Cluster to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the kubernetes connected cluster should exist. Changing this forces a new kubernetes connected cluster to be created.

* `agent_public_key_certificate` - (Required) Base64 encoded public certificate used by the agent to do the initial handshake to the backend services in Azure. Changing this forces a new kubernetes connected cluster to be created.

* `identity` - (Required) An `identity` block as defined below.

* `location` - (Required) The Azure Region where the kubernetes connected cluster should exist. Changing this forces a new kubernetes connected cluster to be created.

* `distribution` - (Optional) The Kubernetes distribution running on this connected cluster. Changing this forces a new kubernetes connected cluster to be created.

* `infrastructure` - (Optional) The infrastructure on which the Kubernetes cluster represented by this connected cluster is running on. Changing this forces a new kubernetes connected cluster to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the kubernetes connected cluster.

---

An `identity` block exports the following:

* `type` - (Required) Specifies the type of Managed Service Identity. Possible values are `SystemAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the kubernetes connected cluster.

* `identity` - An `identity` block as defined below.

* `provisioning_state` - The current deployment state of connected clusters.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the kubernetes connected cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the kubernetes connected cluster.
* `update` - (Defaults to 30 minutes) Used when updating the kubernetes connected cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the kubernetes connected cluster.

## Import

kubernetes connected cluster can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_connected_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Kubernetes/connectedClusters/cluster1
```
