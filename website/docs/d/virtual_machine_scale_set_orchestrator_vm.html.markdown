---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_scale_set_orchestrator_vm"
description: |-
  Gets information about an existing Virtual Machine Scale Set Orchestrator VM.
---

# Data Source: azurerm_virtual_machine_scale_set_orchestrator_vm

Use this data source to access information about an existing Virtual Machine Scale Set Orchestrator VM.

## Example Usage

```hcl
provider "azurerm" {
features {}
}

data "azurerm_virtual_machine_scale_set_orchestrator_vm" "example" {
name = "example-VMSS"
resource_group_name = "example-resources"
}

output "id" {
  value = data.azurerm_virtual_machine_scale_set_orchestrator_vm.example.id
}
```

## Argument Reference

* `name` - (Required) The name which should be used for this Virtual Machine Scale Set Orchestrator VM.

* `resource_group_name` - (Required) The name of the resource group where this Virtual Machine Scale Set Orchestrator VM exists.

## Attributes Reference

* `id` - The ID of this Virtual Machine Scale Set Orchestrator VM.

* `location` - The location where this Virtual Machine Scale Set Orchestrator VM exists.

* `tags` - A mapping of tags assigned to this Virtual Machine Scale Set Orchestrator VM.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Machine Scale Set Orchestrator VM.

