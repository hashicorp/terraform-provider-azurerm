---
subcategory: "Automanage"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_automanage_configuration_profile_hcrpassignment"
description: |-
  Gets information about an existing automanage ConfigurationProfileHCRPAssignment.
---

# Data Source: azurerm_automanage_configuration_profile_hcrpassignment

Use this data source to access information about an existing automanage ConfigurationProfileHCRPAssignment.

## Example Usage

```hcl
data "azurerm_automanage_configuration_profile_hcrpassignment" "example" {
  name = "example-configurationprofilehcrpassignment"
  resource_group_name = "existing"
  machine_name = "existing"
}

output "id" {
  value = data.azurerm_automanage_configuration_profile_hcrpassignment.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this automanage ConfigurationProfileHCRPAssignment.

* `resource_group_name` - (Required) The name of the Resource Group where the automanage ConfigurationProfileHCRPAssignment exists.

* `machine_name` - (Required) The name of the Arc machine.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the automanage ConfigurationProfileHCRPAssignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the automanage ConfigurationProfileHCRPAssignment.