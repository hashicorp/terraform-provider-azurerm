---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_backup_policy_vm"
description: |-
  Gets information about an existing VM Backup Policy.
---

# Data Source: azurerm_backup_policy_vm

Use this data source to access information about an existing VM Backup Policy.

## Example Usage

```hcl
data "azurerm_backup_policy_vm" "policy" {
  name                = "policy"
  recovery_vault_name = "recovery_vault"
  resource_group_name = "resource_group"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - Specifies the name of the VM Backup Policy.

* `recovery_vault_name` - Specifies the name of the Recovery Services Vault.

* `resource_group_name` - The name of the resource group in which the VM Backup Policy resides.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Backup VM Protection Policy.

* `backup` - A `backup` block as defined below.

* `consistency_type` - The consistency type for the backup policy.

* `instant_restore_resource_group` - An `instant_restore_resource_group` block as defined below.

* `instant_restore_retention_days` - The instant restore retention range in days.

* `policy_type` - The type of the Backup Policy.

* `retention_daily` - A `retention_daily` block as defined below.

* `retention_monthly` - A `retention_monthly` block as defined below.

* `retention_weekly` - A `retention_weekly` block as defined below.

* `retention_yearly` - A `retention_yearly` block as defined below.

* `tiering_policy` - A `tiering_policy` block as defined below.

* `timezone` - The timezone.

---

A `backup` block exports the following:

* `frequency` - The backup frequency.

* `hour_duration` - The duration of the backup window in hours.

* `hour_interval` - The interval in hours at which backup is triggered.

* `time` - The time of day to perform the backup.

* `weekdays` - The days of the week the backup is performed on.

---

An `instant_restore_resource_group` block exports the following:

* `prefix` - The prefix for the `instant_restore_resource_group` name.

* `suffix` - The suffix for the `instant_restore_resource_group` name.

---

A `retention_daily` block exports the following:

* `count` - The number of daily backups to keep.

---

A `retention_weekly` block exports the following:

* `count` - The number of weekly backups to keep.

* `weekdays` - The weekday backups retained.

---

A `retention_monthly` block exports the following:

* `count` - The number of monthly backups to keep.

* `days` - The days of the month backups are retained on.

* `include_last_days` - Whether the last day of the month is included.

* `weekdays` - The weekday backups retained.

* `weeks` - The weeks of the month backups are retained on.

---

A `retention_yearly` block exports the following:

* `count` - The number of yearly backups to keep.

* `days` - The days of the month backups are retained on.

* `include_last_days` - Whether the last day of the month is included.

* `months` - The months of the year backups are retained on.

* `weekdays` - The weekday backups retained.

* `weeks` - The weeks of the month backups are retained on.

---

A `tiering_policy` block exports the following:

* `archived_restore_point` - An `archived_restore_point` block as defined below.

---

An `archived_restore_point` block exports the following:

* `duration` - The number of days/weeks/months/years backups are retained in the current tier before tiering.

* `duration_type` - The retention duration type.

* `mode` - The tiering mode used to control automatic tiering of recovery points.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Recovery Services VM Protection Policy.
