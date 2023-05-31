---
subcategory: "Automanage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automanage_configuration"
description: |-
  Manages a automanage ConfigurationProfile.
---

# azurerm_automanage_configuration

Manages a automanage ConfigurationProfile.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-automanage"
  location = "West Europe"
}

resource "azurerm_automanage_configuration" "example" {
  name                = "example-acmp"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  antimalware {
    exclusions {
      extensions = "exe;dll"
      paths      = "C:\\Windows\\Temp;D:\\Temp"
      processes  = "svchost.exe;notepad.exe"
    }

    real_time_protection_enabled   = true
    scheduled_scan_enabled         = true
    scheduled_scan_type            = "Quick"
    scheduled_scan_day             = 1
    scheduled_scan_time_in_minutes = 1339
  }

  automation_account_enabled  = true
  boot_diagnostics_enabled    = true
  defender_for_cloud_enabled  = true
  guest_configuration_enabled = true
  status_change_alert_enabled = true

  tags = {
    "env" = "test"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this automanage ConfigurationProfile. Changing this forces a new automanage ConfigurationProfile to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the automanage ConfigurationProfile should exist. Changing this forces a new automanage ConfigurationProfile to be created.

* `location` - (Required) The Azure Region where the automanage ConfigurationProfile should exist. Changing this forces a new automanage ConfigurationProfile to be created.

* `antimalware` - (Optional) A `antimalware` block as defined below.

* `automation_account_enabled` - (Optional) Whether the automation account is enabled. Defaults to `false`.

* `boot_diagnostics_enabled` - (Optional) Whether the boot diagnostics is enabled. Defaults to `false`.

* `defender_for_cloud_enabled` - (Optional) Whether the defender for cloud is enabled. Defaults to `false`.

* `guest_configuration_enabled` - (Optional) Whether the guest configuration is enabled. Defaults to `false`.

* `status_change_alert_enabled` - (Optional) Whether the status change alert is enabled. Defaults to `false`.

---

* `antimalware` supports the following:

* `exclusions` - (Optional) A `exclusions` block as defined below.

* `real_time_protection_enabled` - (Optional) Whether the real time protection is enabled. Defaults to `false`.

* `scheduled_scan_enabled` - (Optional) Whether the scheduled scan is enabled. Defaults to `false`.

* `scheduled_scan_type` - (Optional) The type of the scheduled scan. Possible values are `Quick` and `Full`. Defaults to `Quick`.

* `scheduled_scan_day` - (Optional) The day of the scheduled scan. Possible values are `0` to `8` where `0` is daily, `1` to `7` are the days of the week and `8` is Disabled. Defaults to `8`.

* `scheduled_scan_time_in_minutes` - (Optional) The time of the scheduled scan in minutes. Possible values are `0` to `1439` where `0` is 12:00 AM and `1439` is 11:59 PM. 

---

* `exclusions` supports the following:

* `extensions` - (Optional) The extensions to exclude from the antimalware scan, separated by `;`. For example `.ext1;.ext2`.

* `paths` - (Optional) The paths to exclude from the antimalware scan, separated by `;`. For example `C:\\Windows\\Temp;D:\\Temp`.

* `processes` - (Optional) The processes to exclude from the antimalware scan, separated by `;`. For example `svchost.exe;notepad.exe`.

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
terraform import azurerm_automanage_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automanage/configurationProfiles/configurationProfile1
```
