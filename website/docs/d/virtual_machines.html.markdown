---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machines"
sidebar_current: "docs-azurerm-datasource-virtual-machines-x"
description: |-
  Gets information about existing Virtual Machines.
---

# Data Source: azurerm_virtual_machines

Use this data source to access information about existing Virtual Machines.

## Example Usage

```hcl
data "azurerm_virtual_machines" "test" {
  name_prefix         = "prod-"
  resource_group_name = "networking"
}

output "virtual_machine_ids" {
  value = "${join(", ", data.azurerm_virtual_machine.test.ids)}"
}

output "virtual_machine_names" {
  value = "${join(", ", data.azurerm_virtual_machine.test.names)}"
}
```

## Argument Reference

* `name_prefix` - (Optional) Specifies the prefix to match the name of the Virtual Machines.
* `resource_group_name` - (Required) Specifies the name of the resource group the Virtual Machines are located in.

## Attributes Reference

* `ids` - List of IDs of the Virtual Machines.
* `names` - List of names of the Virtual Machines.
