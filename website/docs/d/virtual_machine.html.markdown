---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine"
sidebar_current: "docs-azurerm-datasource-virtual-machine"
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
  value = "${data.azurerm_virtual_machine.example.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Virtual Machine.
* `resource_group_name` - (Required) Specifies the name of the resource group the Virtual Machine is located in.

## Attributes Reference

* `id` - The ID of the Virtual Machine.
