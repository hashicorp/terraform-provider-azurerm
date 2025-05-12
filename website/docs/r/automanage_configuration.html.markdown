---
subcategory: "Automanage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automanage_configuration"
description: |-
  Manages an Automanage Configuration.
---

# azurerm_automanage_configuration

Manages an Automanage Configuration.

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

  azure_security_baseline {
    assignment_type = "ApplyAndAutoCorrect"
  }

  automation_account_enabled = true

  backup {
    policy_name                        = "acctest-backup-policy-%d"
    time_zone                          = "UTC"
    instant_rp_retention_range_in_days = 2

    schedule_policy {
      schedule_run_frequency = "Daily"
      schedule_run_days      = ["Monday", "Tuesday"]
      schedule_run_times     = ["12:00"]
      schedule_policy_type   = "SimpleSchedulePolicy"
    }

    retention_policy {
      retention_policy_type = "LongTermRetentionPolicy"

      daily_schedule {
        retention_times = ["12:00"]
        retention_duration {
          count         = 7
          duration_type = "Days"
        }
      }

      weekly_schedule {
        retention_times = ["14:00"]
        retention_duration {
          count         = 4
          duration_type = "Weeks"
        }
      }
    }
  }

  boot_diagnostics_enabled    = true
  defender_for_cloud_enabled  = true
  guest_configuration_enabled = true
  log_analytics_enabled       = true
  status_change_alert_enabled = true

  tags = {
    "env" = "test"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Automanage Configuration. Changing this forces a new Automanage Configuration to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Automanage Configuration should exist. Changing this forces a new Automanage Configuration to be created.

* `location` - (Required) The Azure Region where the Automanage Configuration should exist. Changing this forces a new Automanage Configuration to be created.

* `antimalware` - (Optional) A `antimalware` block as defined below.

* `azure_security_baseline` - (Optional) A `azure_security_baseline` block as defined below.

* `backup` - (Optional) A `backup` block as defined below.

* `automation_account_enabled` - (Optional) Whether the automation account is enabled. Defaults to `false`.

* `boot_diagnostics_enabled` - (Optional) Whether the boot diagnostics are enabled. Defaults to `false`.

* `defender_for_cloud_enabled` - (Optional) Whether the defender for cloud is enabled. Defaults to `false`.

* `guest_configuration_enabled` - (Optional) Whether the guest configuration is enabled. Defaults to `false`.

* `log_analytics_enabled` - (Optional) Whether log analytics are enabled. Defaults to `false`.

* `status_change_alert_enabled` - (Optional) Whether the status change alert is enabled. Defaults to `false`.

-> **Note:** `status_change_alert_enabled` requires that the `AutomanageAlertsEnabled` feature is enabled. To enable this feature for your subscription, use the following command: `az feature register --namespace Microsoft.Automanage --name AutomanageAlertsEnabled`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `antimalware` block supports the following:

* `exclusions` - (Optional) A `exclusions` block as defined below.

* `real_time_protection_enabled` - (Optional) Whether the real time protection is enabled. Defaults to `false`.

* `scheduled_scan_enabled` - (Optional) Whether the scheduled scan is enabled. Defaults to `false`.

* `scheduled_scan_type` - (Optional) The type of the scheduled scan. Possible values are `Quick` and `Full`. Defaults to `Quick`.

* `scheduled_scan_day` - (Optional) The day of the scheduled scan. Possible values are `0` to `8` where `0` is daily, `1` to `7` are the days of the week and `8` is Disabled. Defaults to `8`.

* `scheduled_scan_time_in_minutes` - (Optional) The time of the scheduled scan in minutes. Possible values are `0` to `1439` where `0` is 12:00 AM and `1439` is 11:59 PM. 

---

The `exclusions` block supports the following:

* `extensions` - (Optional) The extensions to exclude from the antimalware scan, separated by `;`. For example `.ext1;.ext2`.

* `paths` - (Optional) The paths to exclude from the antimalware scan, separated by `;`. For example `C:\\Windows\\Temp;D:\\Temp`.

* `processes` - (Optional) The processes to exclude from the antimalware scan, separated by `;`. For example `svchost.exe;notepad.exe`.

---

The `azure_security_baseline` block supports the following:

* `assignment_type` - (Optional) The assignment type of the azure security baseline. Possible values are `ApplyAndAutoCorrect`, `ApplyAndMonitor`, `Audit` and `DeployAndAutoCorrect`. Defaults to `ApplyAndAutoCorrect`.

---

The `backup` block supports the following:

* `policy_name` - (Optional) The name of the backup policy.

* `time_zone` - (Optional) The timezone of the backup policy. Defaults to `UTC`.

* `instant_rp_retention_range_in_days` - (Optional) The retention range in days of the backup policy. Defaults to `5`.

* `schedule_policy` - (Optional) A `schedule_policy` block as defined below.

* `retention_policy` - (Optional) A `retention_policy` block as defined below.

---

The `schedule_policy` block supports the following:

* `schedule_run_frequency` - (Optional) The schedule run frequency of the backup policy. Possible values are `Daily` and `Weekly`. Defaults to `Daily`.

* `schedule_run_times` - (Optional) The schedule run times of the backup policy.

* `schedule_run_days` - (Optional) The schedule run days of the backup policy. Possible values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` and `Saturday`.

* `schedule_policy_type` - (Optional) The schedule policy type of the backup policy. Possible value is `SimpleSchedulePolicy`. Defaults to `SimpleSchedulePolicy`.

---

The `retention_policy` block supports the following:

* `retention_policy_type` - (Optional) The retention policy type of the backup policy. Possible value is `LongTermRetentionPolicy`. Defaults to `LongTermRetentionPolicy`.

* `daily_schedule` - (Optional) A `daily_schedule` block as defined below.

* `weekly_schedule` - (Optional) A `weekly_schedule` block as defined below.

---

The `daily_schedule` block supports the following:

* `retention_times` - (Optional) The retention times of the backup policy.

* `retention_duration` - (Optional) A `retention_duration` block as defined below.

---

The `weekly_schedule` block supports the following:

* `retention_times` - (Optional) The retention times of the backup policy.

* `retention_duration` - (Optional) A `retention_duration` block as defined below.

---

The `retention_duration` block supports the following:

* `count` - (Optional) The count of the retention duration of the backup policy. Valid value inside `daily_schedule` is `7` to `9999` and inside `weekly_schedule` is `1` to `5163`.

* `duration_type` - (Optional) The duration type of the retention duration of the backup policy. Valid value inside `daily_schedule` is `Days` and inside `weekly_schedule` is `Weeks`. Defaults to `Days`.

---
## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Automanage Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automanage Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automanage Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the Automanage Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Automanage Configuration.

## Import

Automanage Configuration can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automanage_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AutoManage/configurationProfiles/configurationProfile1
```
