---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_backup_policy_file_share"
description: |-
  Gets information about an existing existing File Share Backup Policy.
---

# Data Source: azurerm_backup_policy_file_share

Use this data source to access information about an existing File Share Backup Policy.

## Example Usage

```hcl
data "azurerm_backup_policy_file_share" "policy" {
  name                = "policy"
  recovery_vault_name = "recovery_vault"
  resource_group_name = "resource_group"
}
```

## Arguments Reference

The following arguments are supported:

- `name` - Specifies the name of the File Share Backup Policy.

- `recovery_vault_name` - Specifies the name of the Recovery Services Vault.

- `resource_group_name` - The name of the resource group in which the File Share Backup Policy resides.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the File Share Backup Policy.

- `backup` - A `backup` block as defined below.

- `backup_tier` - The backup tier to use.

- `retention_daily` - A `retention_daily` block as defined below.

- `retention_monthly` - A `retention_monthly` block as defined below.

- `retention_weekly` - A `retention_weekly` block as defined below.

- `retention_yearly` - A `retention_yearly` block as defined below.

- `snapshot_retention_in_days` - The number of days to retain the snapshots.

- `timezone` - The timezone.

---

A `backup` block exports the following:

- `frequency` - The backup frequency.

- `hourly` - An `hourly` block as defined below.

- `time` - The time of day to perform the backup.

---

An `hourly` block exports the following:

- `interval` - The interval at which backup is triggered.

- `start_time` - The start time of the hourly backup.

- `window_duration` - The duration of the backup window in hours.

---

A `retention_daily` block exports the following:

- `count` - The number of daily backups to keep.

---

A `retention_weekly` block exports the following:

- `count` - The number of weekly backups to keep.

- `weekdays` - The weekday backups retained.

---

A `retention_monthly` block exports the following:

- `count` - The number of monthly backups to keep.

- `days` - The days of the month backups are retained on.

- `include_last_days` - Whether the last day of the month is included.

- `weekdays` - The weekday backups retained.

- `weeks` - The weeks of the month backups are retained on.

---

A `retention_yearly` block exports the following:

- `count` - The number of yearly backups to keep.

- `days` - The days of the month backups are retained on.

- `include_last_days` - Whether the last day of the month is included.

- `months` - The months of the year backups are retained on.

- `weekdays` - The weekday backups retained.

- `weeks` - The weeks of the month backups are retained on.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Recovery Services File Share Protection Policy.
