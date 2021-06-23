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

  my_resource_name = data.azurerm_resource_id.example["instanceName"]
  # set to "MyResource"
}
```

## Argument Reference

* `resource_id` - The Azure resource id to parse.

## Attributes Reference

* `subscription_id` - The parsed Azure subscription.

* `resource_group_name` - The parsed Azure resource group name.

* `resource_type` - The type of the primary resource. (e.g. `Microsoft.Network/virtualNetworks`).

* `secondary_resource_type` - The type of the child resource.

* `parts` - A map of any additional key-value pairs in the path, this includes the resource name, accessed using an index of the key name.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 seconds) Used when parsing a resource id.
