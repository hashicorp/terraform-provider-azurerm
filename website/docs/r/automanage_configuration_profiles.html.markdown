---
subcategory: "Automanage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automanage_configuration_profile"
description: |-
Manages a automanage ConfigurationProfile.
---

# azurerm_automanage_configuration_profile

Manages a automanage ConfigurationProfile.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-automanage"
  location = "West Europe"
}

resource "azurerm_automanage_configuration_profile" "example" {
  name                = "example-configurationprofile"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  configuration_json  = jsonencode({
    "Antimalware/Enable":false,
    "AzureSecurityCenter/Enable":true,
    "Backup/Enable":false,
    "BootDiagnostics/Enable":true,
    "ChangeTrackingAndInventory/Enable":true,
    "GuestConfiguration/Enable":true,
    "LogAnalytics/Enable":true,
    "UpdateManagement/Enable":true,
    "VMInsights/Enable":true
  })
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this automanage ConfigurationProfile. Changing this forces a new automanage ConfigurationProfile to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the automanage ConfigurationProfile should exist. Changing this forces a new automanage ConfigurationProfile to be created.

* `location` - (Required) The Azure Region where the automanage ConfigurationProfile should exist. Changing this forces a new automanage ConfigurationProfile to be created.

* `configuration_json` - (Required) configuration dictionary of the configuration profile. Changing this forces a new automanage ConfigurationProfile to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the automanage ConfigurationProfile.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the automanage ConfigurationProfile.

* `type` - The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts".

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the automanage ConfigurationProfile.
* `read` - (Defaults to 5 minutes) Used when retrieving the automanage ConfigurationProfile.
* `update` - (Defaults to 30 minutes) Used when updating the automanage ConfigurationProfile.
* `delete` - (Defaults to 30 minutes) Used when deleting the automanage ConfigurationProfile.

## Import

automanage ConfigurationProfiles can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automanage_configuration_profile.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automanage/configurationProfiles/configurationProfile1
```