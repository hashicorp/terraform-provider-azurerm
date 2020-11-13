---
subcategory: "TODO - pick from: Load Balancer|Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_virtual_wan"
description: |-
  Gets information about an existing Virtual Wan.
---

# Data Source: azurerm_virtual_wan

Use this data source to access information about an existing Virtual Wan.

## Example Usage

```hcl
data "azurerm_virtual_wan" "example" {
  name = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_virtual_wan.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Virtual Wan.

* `resource_group_name` - (Required) The name of the Resource Group where the Virtual Wan exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Virtual Wan.

* `address_prefix` - TODO.

* `location` - The Azure Region where the Virtual Wan exists.

* `tags` - A mapping of tags assigned to the Virtual Wan.

* `virtual_wan_id` - The ID of the TODO.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Wan.