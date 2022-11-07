---
subcategory: "Automanage"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_automanage_configuration_profiles_version"
description: |-
  Gets information about an existing automanage ConfigurationProfilesVersion.
---

# Data Source: azurerm_automanage_configuration_profiles_version

Use this data source to access information about an existing automanage ConfigurationProfilesVersion.

## Example Usage

```hcl
data "azurerm_automanage_configuration_profiles_version" "example" {
  name                       = "example-configurationprofilesversion"
  resource_group_name        = "existing"
  configuration_profile_name = "existing"
}

output "id" {
  value = data.azurerm_automanage_configuration_profiles_version.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this automanage ConfigurationProfilesVersion.

* `resource_group_name` - (Required) The name of the Resource Group where the automanage ConfigurationProfilesVersion exists.

* `configuration_profile_name` - (Required) The configuration profile name.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the automanage ConfigurationProfilesVersion.

* `location` - The Azure Region where the automanage ConfigurationProfilesVersion exists.

* `tags` - A mapping of tags assigned to the automanage ConfigurationProfilesVersion.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the automanage ConfigurationProfilesVersion.
