---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine"
description: |-
  Gets information about an existing Virtual Machine.
---

# Data Source: azurerm_virtual_machine

Use this data source to access information about an existing Virtual Machine.

## Example Usage

```hcl
data "azurerm_virtual_machine" "example" {
  name                = "production"
  resource_group_name = "networking"
}

output "virtual_machine_id" {
  value = data.azurerm_virtual_machine.example.id
}
```

## Argument Reference

* `name` - Specifies the name of the Virtual Machine.
* `resource_group_name` - Specifies the name of the resource group the Virtual Machine is located in.

## Attributes Reference

* `id` - The ID of the Virtual Machine.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Machine.
