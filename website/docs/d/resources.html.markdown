---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resources"
description: |-
  Gets information about an existing Resources.
---

# Data Source: azurerm_resources

Use this data source to access information about existing resources.

## Example Usage

```hcl
# Get Resources from a Resource Group
data "azurerm_resources" "example" {
  resource_group_name = "example-resources"
}

# Get Resources with specific Tags
data "azurerm_resources" "example" {
  resource_group_name = "example-resources"

  required_tags = {
    environment = "production"
    role        = "webserver"
  }
}

# Get resources by type, create spoke vNet peerings
data "azurerm_resources" "spokes" {
  type = "Microsoft.Network/virtualNetworks"

  required_tags = {
    environment = "production"
    role        = "spokeNetwork"
  }
}

resource "azurerm_virtual_network_peering" "spoke_peers" {
  count = length(data.azurerm_resources.spokes.resources)

  name                      = "hub2${data.azurerm_resources.spokes.resources[count.index].name}"
  resource_group_name       = azurerm_resource_group.hub.name
  virtual_network_name      = azurerm_virtual_network.hub.name
  remote_virtual_network_id = data.azurerm_resources.spokes.resources[count.index].id
}
```

## Argument Reference

~> **Note:** At least one of `name`, `resource_group_name` or `type` must be specified.

* `name` - (Optional) The name of the Resource.

* `resource_group_name` - (Optional) The name of the Resource group where the Resources are located.

* `type` - (Optional) The Resource Type of the Resources you want to list (e.g. `Microsoft.Network/virtualNetworks`). A resource type's name follows the format: `{resource-provider}/{resource-type}`. The resource type for a key vault is `Microsoft.KeyVault/vaults`. A full list of available Resource Providers can be found [here](https://docs.microsoft.com/azure/azure-resource-manager/azure-services-resource-providers). A full list of Resources Types can be found [here](https://learn.microsoft.com/en-us/azure/templates/#find-resources).

* `required_tags` - (Optional) A mapping of tags which the resource has to have in order to be included in the result.

## Attributes Reference

* `resources` - One or more `resource` blocks as defined below.

---

The `resource` block exports the following:

* `name` - The name of this Resource.

* `id` - The ID of this Resource.

* `resource_group_name` - The name of the Resource Group in which this Resource exists.

* `type` - The type of this Resource. (e.g. `Microsoft.Network/virtualNetworks`).

* `location` - The Azure Region in which this Resource exists.

* `tags` - A map of tags assigned to this Resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Resources.
