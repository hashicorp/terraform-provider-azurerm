---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry_agent_pool"
description: |-
  Manages a Container Registry Agent Pool.
---

# azurerm_container_registry_agent_pool

Manages a Container Registry Agent Pool.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_container_registry" "example" {
  name                = "example-acr"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Premium"
}

resource "azurerm_container_registry_agent_pool" "example" {
  name                  = "example-agent-pool"
  container_registry_id = azurerm_container_registry.example.id
  tier                  = "S1"
}
```

## Arguments Reference

The following arguments are supported:

* `container_registry_id` - (Required) The ID of the Container Registry that this Agent Pool resides in. Changing this forces a new Container Registry Agent Pool to be created.

* `name` - (Required) The name which should be used for this Container Registry Agent Pool. Changing this forces a new Container Registry Agent Pool to be created.

* `tier` - (Required) The tier of the agent machine for this Container Registry Agent Pool. Possible values are `S1`, `S2`, `S3` and `I6`. Changing this forces a new Container Registry Agent Pool to be created.

---

* `agent_count` - (Optional) The count of the agent machine for this Container Registry Agent Pool. Defaults to `1`.

* `os` - (Optional) The OS of the agent machine for this Container Registry Agent Pool. Defaults to `Linux` (currently the only choice). Changing this forces a new Container Registry Agent Pool to be created.

* `subnet_id` - (Optional) The ID of the Subnet of the agent machine for this Container Registry Agent Pool. Changing this forces a new Container Registry Agent Pool to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Container Registry Agent Pool. Changing this forces a new Container Registry Agent Pool to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Container Registry Agent Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Container Registry Agent Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container Registry Agent Pool.
* `update` - (Defaults to 1 hour) Used when updating the Container Registry Agent Pool.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container Registry Agent Pool.

## Import

Container Registry Agent Pools can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_registry_agent_pool.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.ContainerRegistry/registries/registry1/agentPools/pool1
```
