---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_policy_postgresql"
description: |-
  Manages a Backup Policy Postgre Sql.
---

# azurerm_data_protection_backup_policy_postgresql

Manages a Backup Policy Postgre Sql.

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_protection_backup_vault" "example" {
  name                = "example-backup-vault"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
}

resource "azurerm_data_protection_backup_policy_postgresql" "example" {
  name                = "example-backup-policy"
  resource_group_name = azurerm_resource_group.rg.name
  vault_name          = azurerm_data_protection_backup_vault.example.name

  backup_rules {
    name                     = "backup"
    repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]
  }
  default_retention_duration = "P4M"
  retention_rules {
    name             = "weekly"
    duration         = "P6M"
    tagging_priority = 20
    tagging_criteria {
      absolute_criteria = "FirstOfWeek"
    }
  }

  retention_rules {
    name             = "thursday"
    duration         = "P1W"
    tagging_priority = 25
    tagging_criteria {
      days_of_the_week = ["Thursday"]
      schedule_times   = ["2021-05-23T02:30:00Z"]
    }
  }

  retention_rules {
    name             = "monthly"
    duration         = "P1D"
    tagging_priority = 30
    tagging_criteria {
      weeks_of_the_month = ["First", "Last"]
      days_of_the_week   = ["Tuesday"]
      schedule_times     = ["2021-05-23T02:30:00Z"]
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Backup Policy Postgre Sql. Changing this forces a new Backup Policy Postgre Sql to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Backup Policy Postgre Sql should exist. Changing this forces a new Backup Policy Postgre Sql to be created.

* `vault_name` - (Required) The name of the Backup Vault where the Backup Policy Postgre Sql should exist. Changing this forces a new Backup Policy Postgre Sql to be created.

* `backup_rules` - (Required) A `backup_rules` block as defined below. Changing this forces a new Backup Policy Postgre Sql to be created.

* `default_retention_duration` - (Required) The duration of default retention rule. It should follow `ISO 8601` duration format. Changing this forces a new Backup Policy Postgre Sql to be created.

---

* `retention_rules` - (Optional) One or more `retention_rules` blocks as defined below. Changing this forces a new Backup Policy Postgre Sql to be created.

---

A `backup_rules` block supports the following:

* `name` - (Required) The name which should be used for this backup rule. Changing this forces a new Backup Policy Postgre Sql to be created.

* `repeating_time_intervals` - (Required) Specifies a list of repeating time interval. It should follow `ISO 8601` repeating time interval . Changing this forces a new Backup Policy Postgre Sql to be created.

---

A `retention_rules` block supports the following:

* `name` - (Required) The name which should be used for this retention rule. Changing this forces a new Backup Policy Postgre Sql to be created.

* `duration` - (Required) Duration of deletion after given timespan. It should follow `ISO 8601` duration format. Changing this forces a new Backup Policy Postgre Sql to be created.

* `tagging_criteria` - (Required) A `tagging_criteria` block as defined below. Changing this forces a new Backup Policy Postgre Sql to be created.

* `tagging_priority` - (Required) Retention Tag priority. Changing this forces a new Backup Policy Postgre Sql to be created.

---

A `tagging_criteria` block supports the following:

* `absolute_criteria` - (Optional) Possible values are `AllBackup`, `FirstOfDay`, `FirstOfWeek`, `FirstOfMonth` and `FirstOfYear`. Changing this forces a new Backup Policy Postgre Sql to be created.

* `days_of_month` - (Optional) One or more `days_of_month` blocks as defined above. Changing this forces a new Backup Policy Postgre Sql to be created.

* `days_of_the_week` - (Optional) Possible values are `Monday`, `Tuesday`, `Thursday`, `Friday`, `Saturday` and `Sunday`. Changing this forces a new Backup Policy Postgre Sql to be created.

* `months_of_year` - (Optional) Possible values are `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November` and `December`. Changing this forces a new Backup Policy Postgre Sql to be created.

* `schedule_times` - (Optional) Specifies a list of schedule times for backup. It should follow `RFC3339` time format. Changing this forces a new Backup Policy Postgre Sql to be created.

* `weeks_of_the_month` - (Optional) Possible values are `First`, `Second`, `Third`, `Fourth` and `Last`. Changing this forces a new Backup Policy Postgre Sql to be created.

---

A `days_of_month` block is day of the month from 1 to 28 otherwise last of month, it supports the following:

* `date` - (Optional) Date of the month. Changing this forces a new Backup Policy Postgre Sql to be created.

* `is_last` - (Optional) Whether Date is last date of month. Changing this forces a new Backup Policy Postgre Sql to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Backup Policy Postgre Sql.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Backup Policy Postgre Sql.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Policy Postgre Sql.
* `update` - (Defaults to 30 minutes) Used when updating the Backup Policy Postgre Sql.
* `delete` - (Defaults to 30 minutes) Used when deleting the Backup Policy Postgre Sql.

## Import

Backup Policy Postgre Sqls can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_policy_postgresql.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataProtection/backupVaults/vault1/backupPolicies/backupPolicy1
```
