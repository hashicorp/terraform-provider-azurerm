---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_policy_postgresql"
description: |-
  Manages a Backup Policy to back up PostgreSQL.
---

# azurerm_data_protection_backup_policy_postgresql

Manages a Backup Policy to back up PostgreSQL. 

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

  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]

  default_retention_duration = "P4M"

  retention_rule {
    name     = "weekly"
    duration = "P6M"
    priority = 20
    criteria {
      absolute_criteria = "FirstOfWeek"
    }
  }

  retention_rule {
    name     = "thursday"
    duration = "P1W"
    priority = 25
    criteria {
      days_of_week           = ["Thursday"]
      scheduled_backup_times = ["2021-05-23T02:30:00Z"]
    }
  }

  retention_rule {
    name     = "monthly"
    duration = "P1D"
    priority = 15
    criteria {
      weeks_of_month         = ["First", "Last"]
      days_of_week           = ["Tuesday"]
      scheduled_backup_times = ["2021-05-23T02:30:00Z"]
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Backup Policy PostgreSQL. Changing this forces a new Backup Policy PostgreSQL to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Backup Policy PostgreSQL should exist. Changing this forces a new Backup Policy PostgreSQL to be created.

* `vault_name` - (Required) The name of the Backup Vault where the Backup Policy PostgreSQL should exist. Changing this forces a new Backup Policy PostgreSQL to be created.

* `backup_repeating_time_intervals` - (Required) Specifies a list of repeating time interval. It supports weekly back. It should follow `ISO 8601` repeating time interval. Changing this forces a new Backup Policy PostgreSQL to be created.
  
* `default_retention_duration` - (Required) The duration of default retention rule. It should follow `ISO 8601` duration format. Changing this forces a new Backup Policy PostgreSQL to be created.

---

* `retention_rule` - (Optional) One or more `retention_rule` blocks as defined below. Changing this forces a new Backup Policy PostgreSQL to be created.

---

A `retention_rule` block supports the following:

* `name` - (Required) The name which should be used for this retention rule. Changing this forces a new Backup Policy PostgreSQL to be created.

* `duration` - (Required) Duration after which the backup is deleted. It should follow `ISO 8601` duration format. Changing this forces a new Backup Policy PostgreSQL to be created.

* `criteria` - (Required) A `criteria` block as defined below. Changing this forces a new Backup Policy PostgreSQL to be created.

* `priority` - (Required) Specifies the priority of the rule. The priority number must be unique for each rule. The lower the priority number, the higher the priority of the rule. Changing this forces a new Backup Policy Postgre Sql to be created.

---

A `criteria` block supports the following:

* `absolute_criteria` - (Optional) Possible values are `AllBackup`, `FirstOfDay`, `FirstOfWeek`, `FirstOfMonth` and `FirstOfYear`. These values mean the first successful backup of the day/week/month/year. Changing this forces a new Backup Policy PostgreSQL to be created.

* `days_of_week` - (Optional) Possible values are `Monday`, `Tuesday`, `Thursday`, `Friday`, `Saturday` and `Sunday`. Changing this forces a new Backup Policy PostgreSQL to be created.

* `months_of_year` - (Optional) Possible values are `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November` and `December`. Changing this forces a new Backup Policy PostgreSQL to be created.

* `scheduled_backup_times` - (Optional) Specifies a list of backup times for backup in the `RFC3339` format. Changing this forces a new Backup Policy Postgre Sql to be created.

* `weeks_of_month` - (Optional) Possible values are `First`, `Second`, `Third`, `Fourth` and `Last`. Changing this forces a new Backup Policy PostgreSQL to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Backup Policy PostgreSQL.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Backup Policy PostgreSQL.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Policy PostgreSQL.
* `update` - (Defaults to 30 minutes) Used when updating the Backup Policy PostgreSQL.
* `delete` - (Defaults to 30 minutes) Used when deleting the Backup Policy PostgreSQL.

## Import

Backup Policy PostgreSQLs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_policy_postgresql.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataProtection/backupVaults/vault1/backupPolicies/backupPolicy1
```
