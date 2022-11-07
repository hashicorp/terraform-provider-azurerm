---
subcategory: "Automanage"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_automanage_configuration_profile_assignment"
description: |-
  Gets information about an existing automanage ConfigurationProfileAssignment.
---

# Data Source: azurerm_automanage_configuration_profile_assignment

Use this data source to access information about an existing automanage ConfigurationProfileAssignment.

## Example Usage

```hcl
data "azurerm_automanage_configuration_profile_assignment" "example" {
  name                = "example-configurationprofileassignment"
  resource_group_name = "existing"
  vm_name             = "existing"
}

output "id" {
  value = data.azurerm_automanage_configuration_profile_assignment.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this automanage ConfigurationProfileAssignment.

* `resource_group_name` - (Required) The name of the Resource Group where the automanage ConfigurationProfileAssignment exists.

* `vm_name` - (Required) The name of the virtual machine.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the automanage ConfigurationProfileAssignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the automanage ConfigurationProfileAssignment.
