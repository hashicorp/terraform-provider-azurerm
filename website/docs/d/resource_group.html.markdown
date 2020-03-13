---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_resource_group"
description: |-
  Gets information about an existing Resource Group.
---

# Data Source: azurerm_resource_group

Use this data source to access information about an existing Resource Group.

## Example Usage

```hcl
data "azurerm_resource_group" "example" {
  name = "existing"
}

output "id" {
  value = data.azurerm_resource_group.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The Name of this Resource Group.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Resource Group.

* `location` - The Azure Region where the Resource Group exists.

* `tags` - A mapping of tags assigned to the Resource Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Group.
