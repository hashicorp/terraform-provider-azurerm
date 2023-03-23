---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry_agent_pool"
description: |-
  Manages an Azure Container Registry Agent Pool.
---

# azurerm_container_registry_agent_pool

Manages an Azure Container Registry Agent Pool.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europ"
}

resource "azurerm_container_registry" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Premium"
}

resource "azurerm_container_registry_agent_pool" "example" {
  name                    = "example"
  resource_group_name     = azurerm_resource_group.example.name
  location                = azurerm_resource_group.example.location
  container_registry_name = azurerm_container_registry.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `container_registry_name` - (Required) Name of Azure Container Registry to create an Agent Pool for. Changing this forces a new Azure Container Registry Agent Pool to be created.

* `location` - (Required) The Azure Region where the Azure Container Registry Agent Pool should exist. Changing this forces a new Azure Container Registry Agent Pool to be created.

* `name` - (Required) The name which should be used for this Azure Container Registry Agent Pool. Changing this forces a new Azure Container Registry Agent Pool to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Container Registry Agent Pool should exist. Changing this forces a new Azure Container Registry Agent Pool to be created.

---

* `instance_count` - (Optional) VMSS instance count. Defaults to `1`.

* `tier` - (Optional) Sets the VM your agent pool will run on. Valid values are: `S1` (2 vCPUs, 3 GiB RAM), `S2` (4 vCPUs, 8 GiB RAM), `S3` (8 vCPUs, 16 GiB RAM) or `I6` (64 vCPUs, 216 GiB RAM, Isolated). Defaults to `S1`. Changing this forces a new Azure Container Registry Agent Pool to be created.

* `virtual_network_subnet_id` - (Optional) The ID of the Virtual Network Subnet Resource where the agent machines will be running. Changing this forces a new Azure Container Registry Agent Pool to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Container Registry Agent Pool.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Container Registry Agent Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Container Registry Agent Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Container Registry Agent Pool.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Container Registry Agent Pool.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Container Registry Agent Pool.

## Import

Azure Container Registry Agent Pool can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_registry_agent_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ContainerRegistry/registries/registry1/agentPools/agentpool1
```
