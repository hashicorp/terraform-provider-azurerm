---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_policy_cosmosdb_account"
description: |-
  Manages a Data Protection Backup Policy for Cosmos DB Database Accounts.
---

# azurerm_data_protection_backup_policy_cosmosdb_account

Manages a Data Protection Backup Policy for Cosmos DB Database Accounts.

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
  backup_schedule                 = ["R/2026-02-08T10:00:00+00:00/P1W"]
  default_retention_duration      = "P10Y"
  time_zone                       = "UTC"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Data Protection Backup Policy for Cosmos DB Database Accounts. Changing this forces a new resource to be created.

-> **Note:** The `name` must be between 3 and 150 characters, contain only letters, numbers, and hyphens, and start with a letter.

* `data_protection_backup_vault_id` - (Required) The ID of the Data Protection Backup Vault where the policy should exist. Changing this forces a new resource to be created.

* `backup_schedule` - (Required) A list of repeating time intervals that define the backup schedule. Changing this forces a new resource to be created.

-> **Note:** Each interval must use the ISO 8601 repeating time interval format.

* `default_retention_duration` - (Required) The duration for which backups are retained by the default retention rule. Changing this forces a new resource to be created.

-> **Note:** The duration must use the ISO 8601 duration format.

* `retention_rule` - (Optional) One or more `retention_rule` blocks as defined below. Changing this forces a new resource to be created.

~> **Note:** Each `retention_rule` block requires at least one of `absolute_criteria` or `days_of_week` to be specified.

* `time_zone` - (Optional) The Windows time zone identifier used by the backup schedule. Changing this forces a new resource to be created. Possible values are `Afghanistan Standard Time`, `Alaskan Standard Time`, `Aleutian Standard Time`, `Altai Standard Time`, `Arab Standard Time`, `Arabian Standard Time`, `Arabic Standard Time`, `Argentina Standard Time`, `Astrakhan Standard Time`, `Atlantic Standard Time`, `AUS Central Standard Time`, `Aus Central W. Standard Time`, `AUS Eastern Standard Time`, `Azerbaijan Standard Time`, `Azores Standard Time`, `Bahia Standard Time`, `Bangladesh Standard Time`, `Belarus Standard Time`, `Bougainville Standard Time`, `Canada Central Standard Time`, `Cape Verde Standard Time`, `Caucasus Standard Time`, `Cen. Australia Standard Time`, `Central America Standard Time`, `Central Asia Standard Time`, `Central Brazilian Standard Time`, `Central Europe Standard Time`, `Central European Standard Time`, `Central Pacific Standard Time`, `Central Standard Time`, `Central Standard Time (Mexico)`, `Chatham Islands Standard Time`, `China Standard Time`, `Cuba Standard Time`, `Dateline Standard Time`, `E. Africa Standard Time`, `E. Australia Standard Time`, `E. Europe Standard Time`, `E. South America Standard Time`, `Easter Island Standard Time`, `Eastern Standard Time`, `Eastern Standard Time (Mexico)`, `Egypt Standard Time`, `Ekaterinburg Standard Time`, `Fiji Standard Time`, `FLE Standard Time`, `Georgian Standard Time`, `GMT Standard Time`, `Greenland Standard Time`, `Greenwich Standard Time`, `GTB Standard Time`, `Haiti Standard Time`, `Hawaiian Standard Time`, `India Standard Time`, `Iran Standard Time`, `Israel Standard Time`, `Jordan Standard Time`, `Kaliningrad Standard Time`, `Kamchatka Standard Time`, `Korea Standard Time`, `Libya Standard Time`, `Line Islands Standard Time`, `Lord Howe Standard Time`, `Magadan Standard Time`, `Magallanes Standard Time`, `Marquesas Standard Time`, `Mauritius Standard Time`, `Mid-Atlantic Standard Time`, `Middle East Standard Time`, `Montevideo Standard Time`, `Morocco Standard Time`, `Mountain Standard Time`, `Mountain Standard Time (Mexico)`, `Myanmar Standard Time`, `N. Central Asia Standard Time`, `Namibia Standard Time`, `Nepal Standard Time`, `New Zealand Standard Time`, `Newfoundland Standard Time`, `Norfolk Standard Time`, `North Asia East Standard Time`, `North Asia Standard Time`, `North Korea Standard Time`, `Omsk Standard Time`, `Pacific SA Standard Time`, `Pacific Standard Time`, `Pacific Standard Time (Mexico)`, `Pakistan Standard Time`, `Paraguay Standard Time`, `Qyzylorda Standard Time`, `Romance Standard Time`, `Russia Time Zone 10`, `Russia Time Zone 11`, `Russia Time Zone 3`, `Russian Standard Time`, `SA Eastern Standard Time`, `SA Pacific Standard Time`, `SA Western Standard Time`, `Saint Pierre Standard Time`, `Sakhalin Standard Time`, `Samoa Standard Time`, `Sao Tome Standard Time`, `Saratov Standard Time`, `SE Asia Standard Time`, `Singapore Standard Time`, `South Africa Standard Time`, `South Sudan Standard Time`, `Sri Lanka Standard Time`, `Sudan Standard Time`, `Syria Standard Time`, `Taipei Standard Time`, `Tasmania Standard Time`, `Tocantins Standard Time`, `Tokyo Standard Time`, `Tomsk Standard Time`, `Tonga Standard Time`, `Transbaikal Standard Time`, `Turkey Standard Time`, `Turks And Caicos Standard Time`, `Ulaanbaatar Standard Time`, `US Eastern Standard Time`, `US Mountain Standard Time`, `UTC`, `UTC-02`, `UTC-08`, `UTC-09`, `UTC-11`, `UTC+12`, `UTC+13`, `Venezuela Standard Time`, `Vladivostok Standard Time`, `Volgograd Standard Time`, `W. Australia Standard Time`, `W. Central Africa Standard Time`, `W. Europe Standard Time`, `W. Mongolia Standard Time`, `West Asia Standard Time`, `West Bank Standard Time`, `West Pacific Standard Time`, `Yakutsk Standard Time`, and `Yukon Standard Time`.

---

A `retention_rule` block supports the following:

* `duration` - (Required) The duration for which backups matching this retention rule are retained. Changing this forces a new resource to be created.

-> **Note:** The duration must use the ISO 8601 duration format.

* `name` - (Required) The name of the retention rule. Changing this forces a new resource to be created.

* `absolute_criteria` - (Optional) The absolute criterion used to identify retained backups. Changing this forces a new resource to be created. Possible values are `AllBackup`, `FirstOfDay`, `FirstOfMonth`, `FirstOfWeek`, and `FirstOfYear`.

* `days_of_week` - (Optional) A set of days of the week on which backups are retained. Changing this forces a new resource to be created. Possible values are `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, and `Sunday`.

* `months_of_year` - (Optional) A set of months in which backups are retained. Changing this forces a new resource to be created. Possible values are `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, and `December`.

* `scheduled_backup_times` - (Optional) A set of scheduled backup times. Changing this forces a new resource to be created.

-> **Note:** Each value must use RFC3339 format.

* `weeks_of_month` - (Optional) A set of weeks of the month in which backups are retained. Changing this forces a new resource to be created. Possible values are `First`, `Second`, `Third`, `Fourth`, and `Last`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Protection Backup Policy for Cosmos DB Database Accounts.

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
