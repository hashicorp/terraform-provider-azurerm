---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_policy_cosmosdb_account"
description: |-
  Manages a Data Protection (Backup Vault) Backup Policy for Cosmos DB Database Accounts.
---

# azurerm_data_protection_backup_policy_cosmosdb_account

Manages a Data Protection (Backup Vault) Backup Policy for Cosmos DB Database Accounts.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_data_protection_backup_vault" "example" {
  name                = "example-data-protection-backup-vault"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
}

resource "azurerm_data_protection_backup_policy_cosmosdb_account" "example" {
  name                            = "example-data-protection-backup-policy-cosmosdb-account"
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.example.id
  backup_schedule                 = "R/2026-02-08T10:00:00+00:00/P1W"
  default_retention_duration      = "P10Y"

  weekly_retention_rule {
    duration = "P4W"
  }

  monthly_retention_rule {
    duration       = "P6M"
    days_of_week   = ["Sunday"]
    weeks_of_month = ["First"]
  }

  yearly_retention_rule {
    duration = "P10Y"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Data Protection Backup Policy for Cosmos DB Database Accounts. Changing this forces a new resource to be created.

-> **Note:** The `name` must be between 3 and 150 characters, contain only letters, numbers, and hyphens, and start with a letter.

* `data_protection_backup_vault_id` - (Required) The ID of the Data Protection Backup Vault where the policy should exist. Changing this forces a new resource to be created.

* `backup_schedule` - (Required) The repeating time interval that specifies the weekday and time when the weekly Full backup runs. Changing this forces a new resource to be created.

-> **Note:** The interval must use `R/YYYY-MM-DDThh:mm:ssZ/P1W`, `R/YYYY-MM-DDThh:mm:ss+hh:mm/P1W`, or `R/YYYY-MM-DDThh:mm:ss-hh:mm/P1W`. Seconds and a time zone suffix are required; `Z` represents UTC. For example, `R/2026-02-08T10:00:00Z/P1W` schedules a weekly Full backup at 10:00 UTC on Sunday. Minute-only timestamps (`Thh:mm`), fractional seconds (`Thh:mm:ss.fff`), and other ISO 8601 variations are not accepted. The `YYYY-MM-DD` component is used to determine the weekday and does not specify the date when backups begin.

* `default_retention_duration` - (Required) The duration for which the default retention rule retains backups. Changing this forces a new resource to be created.

-> **Note:** The duration must use the ISO 8601 duration format.

* `backup_frequency` - (Optional) The frequency at which backups are performed. Possible values are `Daily` and `Weekly`. Defaults to `Daily`. Changing this forces a new resource to be created.

~> **Note:** If `backup_frequency` is set to `Daily`, Azure Backup performs a Full backup once every week, while Incremental backups are attempted on other days. If set to `Weekly`, only the weekly Full backup is scheduled. The weekly Full backup runs on the weekday and at the time specified by `backup_schedule`.

* `monthly_retention_rule` - (Optional) A `monthly_retention_rule` block as defined below. Changing this forces a new resource to be created.

* `weekly_retention_rule` - (Optional) A `weekly_retention_rule` block as defined below. Changing this forces a new resource to be created.

* `yearly_retention_rule` - (Optional) A `yearly_retention_rule` block as defined below. Changing this forces a new resource to be created.

-> **Note:** At most one rule of each category can be configured. Rules select Full backups and are evaluated in yearly, monthly, then weekly priority order.

---

A `monthly_retention_rule` block supports the following:

* `duration` - (Required) The duration for which backups matching this retention rule are retained. Changing this forces a new resource to be created.

-> **Note:** The duration must use the ISO 8601 duration format.

* `days_of_week` - (Optional) A set containing exactly one day of the week whose scheduled backup is retained by this rule. Possible values are `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, and `Sunday`. Changing this forces a new resource to be created.

~> **Note:** The value specified for `days_of_week` must match the day of the week in `backup_schedule`.

* `weeks_of_month` - (Optional) A set of ordinal weeks within each month in which backups matching this rule are retained. Possible values are `First`, `Second`, `Third`, `Fourth`, and `Last`. Changing this forces a new resource to be created.

~> **Note:** Specify `days_of_week` and `weeks_of_month` together to select backups by day and week. When both fields are omitted, this rule retains the first successful full backup of the month.

---

A `weekly_retention_rule` block supports the following:

* `duration` - (Required) The duration for which backups matching this retention rule are retained. Changing this forces a new resource to be created.

-> **Note:** The duration must use the ISO 8601 duration format.

* `days_of_week` - (Optional) A set containing exactly one day of the week whose scheduled backup is retained by this rule. Possible values are `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, and `Sunday`. Changing this forces a new resource to be created.

~> **Note:** The value specified for `days_of_week` must match the day of the week in `backup_schedule`.

-> **Note:** When `days_of_week` is omitted, this rule retains the first successful full backup of the week.

---

A `yearly_retention_rule` block supports the following:

* `duration` - (Required) The duration for which backups matching this retention rule are retained. Changing this forces a new resource to be created.

-> **Note:** The duration must use the ISO 8601 duration format.

* `days_of_week` - (Optional) A set containing exactly one day of the week whose scheduled backup is retained by this rule. Possible values are `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, and `Sunday`. Changing this forces a new resource to be created.

~> **Note:** The value specified for `days_of_week` must match the day of the week in `backup_schedule`.

* `months_of_year` - (Optional) A set of calendar months in which backups matching this rule are retained. Possible values are `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, and `December`. Changing this forces a new resource to be created.

* `weeks_of_month` - (Optional) A set of ordinal weeks within each month in which backups matching this rule are retained. Possible values are `First`, `Second`, `Third`, `Fourth`, and `Last`. Changing this forces a new resource to be created.

~> **Note:** Specify `days_of_week`, `weeks_of_month`, and `months_of_year` together to select backups by day, week, and month. When all three fields are omitted, this rule retains first successful full backup of the year.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Protection Backup Policy for Cosmos DB Database Accounts.

* `incremental_backup_schedules` - A list of ISO 8601 repeating time intervals indicating the days of the week when Incremental backups are attempted.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Protection Backup Policy for Cosmos DB Database Accounts.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Protection Backup Policy for Cosmos DB Database Accounts.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Protection Backup Policy for Cosmos DB Database Accounts.

## Import

A Data Protection Backup Policy for Cosmos DB Database Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_policy_cosmosdb_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DataProtection/backupVaults/backupVault1/backupPolicies/backupPolicy1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DataProtection` - 2026-06-01
