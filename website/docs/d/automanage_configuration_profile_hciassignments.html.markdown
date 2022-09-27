---
subcategory: "Automanage"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_automanage_configuration_profile_hciassignment"
description: |-
  Gets information about an existing automanage ConfigurationProfileHCIAssignment.
---

# Data Source: azurerm_automanage_configuration_profile_hciassignment

Use this data source to access information about an existing automanage ConfigurationProfileHCIAssignment.

## Example Usage

```hcl
data "azurerm_automanage_configuration_profile_hciassignment" "example" {
  name = "example-configurationprofilehciassignment"
  resource_group_name = "existing"
  cluster_name = "existing"
}

output "id" {
  value = data.azurerm_automanage_configuration_profile_hciassignment.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this automanage ConfigurationProfileHCIAssignment.

* `resource_group_name` - (Required) The name of the Resource Group where the automanage ConfigurationProfileHCIAssignment exists.

* `cluster_name` - (Required) The name of the Arc machine.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the automanage ConfigurationProfileHCIAssignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the automanage ConfigurationProfileHCIAssignment.