---
subcategory: "Automanage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automanage_configuration_profiles_version"
description: |-
  Manages a automanage ConfigurationProfilesVersion.
---

# azurerm_automanage_configuration_profiles_version

Manages a automanage ConfigurationProfilesVersion.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-automanage"
  location = "West Europe"
}

resource "azurerm_automanage_configuration_profile" "example" {
  name = "example-configurationprofile"
  resource_group_name = azurerm_resource_group.example.name
  location = azurerm_resource_group.example.location
}

resource "azurerm_automanage_configuration_profiles_version" "example" {
  name = "example-configurationprofilesversion"
  resource_group_name = azurerm_resource_group.example.name
  location = azurerm_resource_group.example.location
  configuration_profile_name = azurerm_automanage_configuration_profile.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this automanage ConfigurationProfilesVersion. Changing this forces a new automanage ConfigurationProfilesVersion to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the automanage ConfigurationProfilesVersion should exist. Changing this forces a new automanage ConfigurationProfilesVersion to be created.

* `location` - (Required) The Azure Region where the automanage ConfigurationProfilesVersion should exist.

* `configuration_profile_name` - (Required) Name of the configuration profile. Changing this forces a new automanage ConfigurationProfilesVersion to be created.

---

* `configuration` - (Optional) configuration dictionary of the configuration profile.

* `tags` - (Optional) A mapping of tags which should be assigned to the automanage ConfigurationProfilesVersion.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the automanage ConfigurationProfilesVersion.

* `type` - The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts".

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the automanage ConfigurationProfilesVersion.
* `read` - (Defaults to 5 minutes) Used when retrieving the automanage ConfigurationProfilesVersion.
* `delete` - (Defaults to 30 minutes) Used when deleting the automanage ConfigurationProfilesVersion.

## Import

automanage ConfigurationProfilesVersions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automanage_configuration_profiles_version.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automanage/configurationProfiles/configurationProfile1/versions/version1
```