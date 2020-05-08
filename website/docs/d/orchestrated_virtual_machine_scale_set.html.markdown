---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_orchestrated_virtual_machine_scale_set"
description: |-
  Gets information about an existing Orchestrated Virtual Machine Scale Set.
---

# Data Source: azurerm_orchestrated_virtual_machine_scale_set

Use this data source to access information about an existing Orchestrated Virtual Machine Scale Set.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_orchestrated_virtual_machine_scale_set" "example" {
  name                = "example-VMSS"
  resource_group_name = "example-resources"
}

output "id" {
  value = data.azurerm_orchestrated_virtual_machine_scale_set.example.id
}
```

## Argument Reference

* `name` - (Required) The name which should be used for this Orchestrated Virtual Machine Scale Set.

* `resource_group_name` - (Required) The name of the resource group where this Orchestrated Virtual Machine Scale Set exists.

## Attributes Reference

* `id` - The ID of this Orchestrated Virtual Machine Scale Set.

* `location` - The location where this Orchestrated Virtual Machine Scale Set exists.

* `tags` - A mapping of tags assigned to this Orchestrated Virtual Machine Scale Set.

* `unique_id` - The Unique ID of this Orchestrated Virtual Machine Scale Set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Orchestrated Virtual Machine Scale Set.
