---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_virtual_machine_scale_set"
description: |-
  Gets information about an existing Virtual Machine Scale Set.
---

# Data Source: azurerm_virtual_machine_scale_set

Use this data source to access information about an existing Virtual Machine Scale Set.

## Example Usage

```hcl
data "azurerm_virtual_machine_scale_set" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_virtual_machine_scale_set.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Virtual Machine Scale Set.

* `resource_group_name` - (Required) The name of the Resource Group where the Virtual Machine Scale Set exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Virtual Machine Scale Set.

* `identity` - A `identity` block as defined below.

---

A `identity` block exports the following:

* `identity_ids` -  The list of User Managed Identity ID's which are assigned to the Virtual Machine Scale Set.

* `principal_id` - The ID of the System Managed Service Principal assigned to the Virtual Machine Scale Set.

* `type` - The identity type of the Managed Identity assigned to the Virtual Machine Scale Set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Machine Scale Set.
