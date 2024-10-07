---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_policy_blob_storage"
description: |-
  Manages a Backup Policy Blob Storage.
---

# azurerm_data_protection_backup_policy_blob_storage

Manages a Backup Policy Blob Storage.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_protection_backup_vault" "example" {
  name                = "example-backup-vault"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
}

resource "azurerm_data_protection_backup_policy_blob_storage" "example" {
  name                                   = "example-backup-policy"
  vault_id                               = azurerm_data_protection_backup_vault.example.id
  operational_default_retention_duration = "P30D"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Backup Policy Blob Storage. Changing this forces a new Backup Policy Blob Storage to be created.

* `vault_id` - (Required) The ID of the Backup Vault within which the Backup Policy Blob Storage should exist. Changing this forces a new Backup Policy Blob Storage to be created.

* `backup_repeating_time_intervals` - (Optional) Specifies a list of repeating time interval. It should follow `ISO 8601` repeating time interval. Changing this forces a new Backup Policy Blob Storage to be created.

* `operational_default_retention_duration` - (Optional) The duration of operational default retention rule. It should follow `ISO 8601` duration format. Changing this forces a new Backup Policy Blob Storage to be created.

* `retention_rule` - (Optional) One or more `retention_rule` blocks as defined below. Changing this forces a new Backup Policy Blob Storage to be created.

-> **Note:** Setting `retention_rule` also requires setting `vault_default_retention_duration`.

* `time_zone` - (Optional) Specifies the Time Zone which should be used by the backup schedule. Changing this forces a new Backup Policy Blob Storage to be created.

* `vault_default_retention_duration` - (Optional) The duration of vault default retention rule. It should follow `ISO 8601` duration format. Changing this forces a new Backup Policy Blob Storage to be created.

-> **Note:** Setting `vault_default_retention_duration` also requires setting `backup_repeating_time_intervals`. At least one of `operational_default_retention_duration` or `vault_default_retention_duration` must be specified.

---

A `retention_rule` block supports the following:

* `name` - (Required) The name which should be used for this retention rule. Changing this forces a new Backup Policy Blob Storage to be created.

* `duration` - (Required) Duration after which the backup is deleted. It should follow `ISO 8601` duration format. Changing this forces a new Backup Policy Blob Storage to be created.

* `criteria` - (Required) A `criteria` block as defined below. Changing this forces a new Backup Policy Blob Storage to be created.

* `life_cycle` - (Required) A `life_cycle` block as defined below. Changing this forces a new Backup Policy Blob Storage to be created.

* `priority` - (Required) Specifies the priority of the rule. The priority number must be unique for each rule. The lower the priority number, the higher the priority of the rule. Changing this forces a new Backup Policy Blob Storage to be created.

---

A `criteria` block supports the following:

* `absolute_criteria` - (Optional) Possible values are `AllBackup`, `FirstOfDay`, `FirstOfWeek`, `FirstOfMonth` and `FirstOfYear`. These values mean the first successful backup of the day/week/month/year. Changing this forces a new Backup Policy Blob Storage to be created.

* `days_of_month` - (Optional) Must be between `0` and `28`. `0` for last day within the month. Changing this forces a new Backup Policy Blob Storage to be created.

* `days_of_week` - (Optional) Possible values are `Monday`, `Tuesday`, `Thursday`, `Friday`, `Saturday` and `Sunday`. Changing this forces a new Backup Policy Blob Storage to be created.

* `months_of_year` - (Optional) Possible values are `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November` and `December`. Changing this forces a new Backup Policy Blob Storage to be created.

* `scheduled_backup_times` - (Optional) Specifies a list of backup times for backup in the `RFC3339` format. Changing this forces a new Backup Policy Blob Storage to be created.

* `weeks_of_month` - (Optional) Possible values are `First`, `Second`, `Third`, `Fourth` and `Last`. Changing this forces a new Backup Policy Blob Storage to be created.

A `life_cycle` block supports the following:

* `data_store_type` - (Required) The type of data store. The only possible value is `VaultStore`. Changing this forces a new Backup Policy Blob Storage to be created.

* `duration` - (Required) The retention duration up to which the backups are to be retained in the data stores. It should follow `ISO 8601` duration format. Changing this forces a new Backup Policy Blob Storage to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Backup Policy Blob Storage.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Backup Policy Blob Storage.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Policy Blob Storage.
* `delete` - (Defaults to 30 minutes) Used when deleting the Backup Policy Blob Storage.

## Import

Backup Policy Blob Storages can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_policy_blob_storage.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataProtection/backupVaults/vault1/backupPolicies/backupPolicy1
```
