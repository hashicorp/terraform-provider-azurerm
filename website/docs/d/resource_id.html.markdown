---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_id"
description: |-
  Parses an Azure resource path into its component parts.
---

# Data Source: azurerm_resource_id

Use this data source to access individual components of an Azure resource id.

## Example Usage

```hcl
# Get Resources from a Resource Group
data "azurerm_resource_id" "example" {
  resource_id = "/subscriptions/c90e9ba4-9a69-49d6-be99-2110471ec1a4/resourceGroups/SomeResourceGroup/providers/Microsoft.ResourceProvider/instanceName/MyResource"
}

locals {
  my_subscription_id = data.azurerm_resource_id.example.subscription_id
  # set to "c90e9ba4-9a69-49d6-be99-2110471ec1a4"

  my_resource_name = data.azurerm_resource_id.example.name
  # set to "MyResource"
}
```

## Argument Reference

* `resource_id` - Resource id to parse.

## Attributes Reference

* `subscription_id` - Resource subscription id.

* `resource_group_name` - Resource group name.

* `provider_namespace` - Resource namespace (e.g. `Microsoft.Network`).

* `resource_type` - Resource type (e.g. `virtualNetworks`).

* `name` - Resource name.

* `full_resource_type` - Full resource type (including parent types if applicable).

* `parent_resources` - A map of parent resource types and names


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 seconds) Used when parsing a resource id.
