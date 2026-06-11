---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_data_protection_backup_policy_blob_storage"
description: |-
  Gets information about an existing Backup Policy (Blob Storage).
---

# Data Source: azurerm_data_protection_backup_policy_blob_storage

Use this data source to access information about an existing Backup Policy (Blob Storage).

## Example Usage

```hcl
data "azurerm_data_protection_backup_policy_blob_storage" "example" {
  name     = "existing-backup-policy"
  vault_id = data.azurerm_data_protection_backup_vault.example.id
}

output "id" {
  value = data.azurerm_data_protection_backup_policy_blob_storage.example.id
}
```

## Arguments Reference

* `name` - (Required) Specifies the name of the Backup Policy.

* `vault_id` - (Required) Specifies the ID of the Backup Vault.

## Attributes Reference

* `id` - The ID of the Backup Policy.

* `backup_repeating_time_intervals` - The backup repeating time intervals.

* `operational_default_retention_duration` - The duration of operational default retention rule.

* `retention_rule` - A `retention_rule` block as defined below.

* `time_zone` - The time zone used by the backup schedule.

* `vault_default_retention_duration` - The duration of vault default retention rule.

---

A `retention_rule` block exports the following:

* `criteria` - A `criteria` block as defined below.

* `life_cycle` - A `life_cycle` block as defined below.

* `name` - The name of the retention rule.

* `priority` - The priority of the rule.

---

A `criteria` block exports the following:

* `absolute_criteria` - The absolute criteria.

* `days_of_month` - The days of the month.

* `days_of_week` - The days of the week.

* `months_of_year` - The months of the year.

* `scheduled_backup_times` - The scheduled backup times.

* `weeks_of_month` - The weeks of the month.

---

A `life_cycle` block exports the following:

* `data_store_type` - The type of data store.

* `duration` - The retention duration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Policy.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.DataProtection` - 2025-07-01
