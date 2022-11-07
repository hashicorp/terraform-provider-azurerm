---
subcategory: "Automanage"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_automanage_configuration_profile"
description: |-
  Gets information about an existing automanage ConfigurationProfile.
---

# Data Source: azurerm_automanage_configuration_profile

Use this data source to access information about an existing automanage ConfigurationProfile.

## Example Usage

```hcl
data "azurerm_automanage_configuration_profile" "example" {
  name                = "example-configurationprofile"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_automanage_configuration_profile.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this automanage ConfigurationProfile.

* `resource_group_name` - (Required) The name of the Resource Group where the automanage ConfigurationProfile exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the automanage ConfigurationProfile.

* `location` - The Azure Region where the automanage ConfigurationProfile exists.

* `tags` - A mapping of tags assigned to the automanage ConfigurationProfile.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the automanage ConfigurationProfile.
