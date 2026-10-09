---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_data_protection_backup_policy_disk"
description: |-
  Gets information about an existing Backup Policy (Disk).
---

# Data Source: azurerm_data_protection_backup_policy_disk

Use this data source to access information about an existing Backup Policy (Disk).

## Example Usage

```hcl
data "azurerm_data_protection_backup_policy_disk" "example" {
  name     = "existing-backup-policy"
  vault_id = data.azurerm_data_protection_backup_vault.example.id
}

output "id" {
  value = data.azurerm_data_protection_backup_policy_disk.example.id
}
```

## Arguments Reference

* `name` - (Required) Specifies the name of the Backup Policy.

* `vault_id` - (Required) Specifies the ID of the Backup Vault.

## Attributes Reference

* `id` - The ID of the Backup Policy.

* `backup_repeating_time_intervals` - The backup repeating time intervals.

* `default_retention_duration` - The duration of default retention rule.

* `retention_rule` - A `retention_rule` block as defined below.

* `time_zone` - The time zone used by the backup schedule.

---

A `retention_rule` block exports the following:

* `criteria` - A `criteria` block as defined below.

* `duration` - The retention duration.

* `name` - The name of the retention rule.

* `priority` - The priority of the rule.

---

A `criteria` block exports the following:

* `absolute_criteria` - The absolute criteria.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Policy.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.DataProtection` - 2025-07-01
