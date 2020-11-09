---
subcategory: "AzureStackHCI"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hyper_converged_cluster"
description: |-
  Manages a Hyper Converged Cluster.
---

# azurerm_hyper_converged_cluster

Manages a Hyper Converged Cluster.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_hyper_converged_cluster" "example" {
  name                = "example-cluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  client_id           = data.azurerm_client_config.current.client_id
  tenant_id           = data.azurerm_client_config.current.tenant_id

  tags = {
    ENV = "Prod"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Hyper Converged Cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Hyper Converged Cluster should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Hyper Converged Cluster should exist. Changing this forces a new resource to be created.

* `client_id` - (Required) The Client ID of the Azure Active Directory which is used by the Hyper Converged Cluster. Changing this forces a new resource to be created.

* `tenant_id` - (Required) The Tenant ID of the Azure Active Directory which is used by the Hyper Converged Cluster. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Hyper Converged Cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Hyper Converged Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Hyper Converged Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Hyper Converged Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Hyper Converged Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Hyper Converged Cluster.

## Import

Hyper Converged Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_hyper_converged_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AzureStackHCI/clusters/cluster1
```
