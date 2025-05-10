---
subcategory: "Maintenance"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_maintenance_configuration"
description: |-
  Manages a Maintenance Configuration.
---

# azurerm_maintenance_configuration

Manages a maintenance configuration.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_maintenance_configuration" "example" {
  name                = "example-mc"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  scope               = "SQLDB"

  tags = {
    Env = "prod"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Maintenance Configuration. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Maintenance Configuration should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specified the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `scope` - (Required) The scope of the Maintenance Configuration. Possible values are `Extension`, `Host`, `InGuestPatch`, `OSImage`, `SQLDB` or `SQLManagedInstance`.

* `visibility` - (Optional) The visibility of the Maintenance Configuration. The only allowable value is `Custom`. Defaults to `Custom`.

* `window` - (Optional) A `window` block as defined below.

* `install_patches` - (Optional) An `install_patches` block as defined below.

-> **Note:** `install_patches` must be specified when `scope` is `InGuestPatch`.

* `in_guest_user_patch_mode` - (Optional) The in guest user patch mode. Possible values are `Platform` or `User`. Must be specified when `scope` is `InGuestPatch`.

* `properties` - (Optional) A mapping of properties to assign to the resource.

* `tags` - (Optional) A mapping of tags to assign to the resource. The key could not contain upper case letter.

---

A `window` block supports:

* `start_date_time` - (Required) Effective start date of the maintenance window in YYYY-MM-DD hh:mm format.

* `expiration_date_time` - (Optional) Effective expiration date of the maintenance window in YYYY-MM-DD hh:mm format.

* `duration` - (Optional) The duration of the maintenance window in HH:mm format.

* `time_zone` - (Required) The time zone for the maintenance window. A list of timezones can be obtained by executing [System.TimeZoneInfo]::GetSystemTimeZones() in PowerShell.

* `recur_every` - (Optional) The rate at which a maintenance window is expected to recur. The rate can be expressed as daily, weekly, or monthly schedules.

---

A `install_patches` block supports:

* `linux` - (Optional) A `linux` block as defined above. This property only applies when `scope` is set to `InGuestPatch`

* `windows` - (Optional) A `windows` block as defined above. This property only applies when `scope` is set to `InGuestPatch`

* `reboot` - (Optional) Possible reboot preference as defined by the user based on which it would be decided to reboot the machine or not after the patch operation is completed. Possible values are `Always`, `IfRequired` and `Never`. This property only applies when `scope` is set to `InGuestPatch`.

---

A `linux` block supports:

* `classifications_to_include` - (Optional) List of Classification category of patches to be patched. Possible values are `Critical`, `Security` and `Other`.

* `package_names_mask_to_exclude` - (Optional) List of package names to be excluded from patching.

* `package_names_mask_to_include` - (Optional) List of package names to be included for patching.

---

A `windows` block supports:

* `classifications_to_include` - (Optional) List of Classification category of patches to be patched. Possible values are `Critical`, `Security`, `UpdateRollup`, `FeaturePack`, `ServicePack`, `Definition`, `Tools` and `Updates`.

* `kb_numbers_to_exclude` - (Optional) List of KB numbers to be excluded from patching.

* `kb_numbers_to_include` - (Optional) List of KB numbers to be included for patching.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Maintenance Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Maintenance Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Maintenance Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the Maintenance Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Maintenance Configuration.

## Import

Maintenance Configuration can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_maintenance_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Maintenance/maintenanceConfigurations/example-mc
```
