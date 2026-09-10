---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_policy_data_lake_storage"
description: |-
  Manages a Backup Policy to Azure Data Lake Storage.
---

# azurerm_data_protection_backup_policy_data_lake_storage

Manages a Backup Policy to Azure Data Lake Storage.

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

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_data_protection_backup_policy_data_lake_storage" "example" {
  name                            = "example-backup-policy"
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.example.id
  backup_schedule                 = ["R/2021-05-23T02:30:00+00:00/P1W"]
  time_zone                       = "India Standard Time"

  default_retention_duration = "P4M"

  retention_rule {
    name              = "weekly"
    duration          = "P6M"
    absolute_criteria = "FirstOfWeek"
  }

  retention_rule {
    name                   = "thursday"
    duration               = "P1W"
    days_of_week           = ["Thursday"]
    scheduled_backup_times = ["2021-05-23T02:30:00Z"]
  }

  retention_rule {
    name                   = "monthly"
    duration               = "P1D"
    weeks_of_month         = ["First", "Last"]
    days_of_week           = ["Tuesday"]
    scheduled_backup_times = ["2021-05-23T02:30:00Z"]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Backup Policy for the Azure Backup Policy Data Lake Storage. Changing this forces a new resource to be created.

* `data_protection_backup_vault_id` - (Required) The ID of the Backup Vault where the Azure Backup Policy Data Lake Storage should exist. Changing this forces a new resource to be created.

* `backup_schedule` - (Required) Specifies a list of repeating time interval, also known as the backup schedule. It supports daily & weekly backup. It should follow [`ISO 8601` recurring time interval format](https://en.wikipedia.org/wiki/ISO_8601#Recurring_intervals), for example: `R/2021-05-23T02:30:00+00:00/P1W`. Changing this forces a new resource to be created.

* `default_retention_duration` - (Required) The retention duration up to which the backups are to be retained in the data stores. It should follow `ISO 8601` duration format. Changing this forces a new resource to be created.

* `retention_rule` - (Optional) One or more `retention_rule` blocks as defined below. The priority of each rule is determined by its order in the list, where the first rule has the highest priority. Changing this forces a new resource to be created.

* `time_zone` - (Optional) Specifies the Time Zone which should be used by the backup schedule. Changing this forces a new resource to be created. Possible values are `Afghanistan Standard Time`,`Alaskan Standard Time`,`Aleutian Standard Time`,`Altai Standard Time`,`Arab Standard Time`,`Arabian Standard Time`,`Arabic Standard Time`,`Argentina Standard Time`,`Astrakhan Standard Time`,`Atlantic Standard Time`,`AUS Central Standard Time`,`Aus Central W. Standard Time`,`AUS Eastern Standard Time`,`Azerbaijan Standard Time`,`Azores Standard Time`,`Bahia Standard Time`,`Bangladesh Standard Time`,`Belarus Standard Time`,`Bougainville Standard Time`,`Canada Central Standard Time`,`Cape Verde Standard Time`,`Caucasus Standard Time`,`Cen. Australia Standard Time`,`Central America Standard Time`,`Central Asia Standard Time`,`Central Brazilian Standard Time`,`Central Europe Standard Time`,`Central European Standard Time`,`Central Pacific Standard Time`,`Central Standard Time`,`Central Standard Time (Mexico)`,`Chatham Islands Standard Time`,`China Standard Time`,`Coordinated Universal Time`,`Cuba Standard Time`,`Dateline Standard Time`,`E. Africa Standard Time`,`E. Australia Standard Time`,`E. Europe Standard Time`,`E. South America Standard Time`,`Easter Island Standard Time`,`Eastern Standard Time`,`Eastern Standard Time (Mexico)`,`Egypt Standard Time`,`Ekaterinburg Standard Time`,`Fiji Standard Time`,`FLE Standard Time`,`Georgian Standard Time`,`GMT Standard Time`,`Greenland Standard Time`,`Greenwich Standard Time`,`GTB Standard Time`,`Haiti Standard Time`,`Hawaiian Standard Time`,`India Standard Time`,`Iran Standard Time`,`Israel Standard Time`,`Jordan Standard Time`,`Kaliningrad Standard Time`,`Kamchatka Standard Time`,`Korea Standard Time`,`Libya Standard Time`,`Line Islands Standard Time`,`Lord Howe Standard Time`,`Magadan Standard Time`,`Magallanes Standard Time`,`Marquesas Standard Time`,`Mauritius Standard Time`,`Mid-Atlantic Standard Time`,`Middle East Standard Time`,`Montevideo Standard Time`,`Morocco Standard Time`,`Mountain Standard Time`,`Mountain Standard Time (Mexico)`,`Myanmar Standard Time`,`N. Central Asia Standard Time`,`Namibia Standard Time`,`Nepal Standard Time`,`New Zealand Standard Time`,`Newfoundland Standard Time`,`Norfolk Standard Time`,`North Asia East Standard Time`,`North Asia Standard Time`,`North Korea Standard Time`,`Omsk Standard Time`,`Pacific SA Standard Time`,`Pacific Standard Time`,`Pacific Standard Time (Mexico)`,`Pakistan Standard Time`,`Paraguay Standard Time`,`Qyzylorda Standard Time`,`Romance Standard Time`,`Russia Time Zone 10`,`Russia Time Zone 11`,`Russia Time Zone 3`,`Russian Standard Time`,`SA Eastern Standard Time`,`SA Pacific Standard Time`,`SA Western Standard Time`,`Saint Pierre Standard Time`,`Sakhalin Standard Time`,`Samoa Standard Time`,`Sao Tome Standard Time`,`Saratov Standard Time`,`SE Asia Standard Time`,`Singapore Standard Time`,`South Africa Standard Time`,`South Sudan Standard Time`,`Sri Lanka Standard Time`,`Sudan Standard Time`,`Syria Standard Time`,`Taipei Standard Time`,`Tasmania Standard Time`,`Tocantins Standard Time`,`Tokyo Standard Time`,`Tomsk Standard Time`,`Tonga Standard Time`,`Transbaikal Standard Time`,`Turkey Standard Time`,`Turks And Caicos Standard Time`,`Ulaanbaatar Standard Time`,`US Eastern Standard Time`,`US Mountain Standard Time`,`UTC`,`UTC-02`,`UTC-08`,`UTC-09`,`UTC-11`,`UTC+12`,`UTC+13`,`Venezuela Standard Time`,`Vladivostok Standard Time`,`Volgograd Standard Time`,`W. Australia Standard Time`,`W. Central Africa Standard Time`,`W. Europe Standard Time`,`W. Mongolia Standard Time`,`West Asia Standard Time`,`West Bank Standard Time`,`West Pacific Standard Time`,`Yakutsk Standard Time` and `Yukon Standard Time`.

---

A `retention_rule` block supports the following:

* `name` - (Required) Specifies the name of the retention rule. Changing this forces a new resource to be created.

* `duration` - (Required) The retention duration up to which the backups are to be retained in the data stores. It should follow `ISO 8601` duration format. Changing this forces a new resource to be created.

* `absolute_criteria` - (Optional) Specifies the absolute criteria for the retention rule. Possible values include `AllBackup`, `FirstOfDay`, `FirstOfWeek`, `FirstOfMonth`, and `FirstOfYear`. These values mean the first successful backup of the day/week/month/year. Changing this forces a new resource to be created.

* `days_of_week` - (Optional) Specifies a list of days of the week on which the retention rule applies. Possible values include `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, and `Sunday`. Changing this forces a new resource to be created.

* `weeks_of_month` - (Optional) Specifies a list of weeks of the month on which the retention rule applies. Possible values include `First`, `Second`, `Third`, `Fourth`, and `Last`. Changing this forces a new resource to be created.

* `months_of_year` - (Optional) Specifies a list of months of the year on which the retention rule applies. Possible values include `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, and `December`. Changing this forces a new resource to be created.

* `scheduled_backup_times` - (Optional) Specifies a list of backup times for backup in the `RFC3339` format. Changing this forces a new resource to be created.

~> **Note:** At least one of `absolute_criteria` or `days_of_week` must be specified. `weeks_of_month` and `months_of_year` are optional and can be supplied together. Multiple intervals may be set using multiple `retention_rule` blocks.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Backup Policy Data Lake Storage.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Backup Policy Data Lake Storage.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Backup Policy Data Lake Storage.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Backup Policy Data Lake Storage.

## Import

Azure Backup Policy Data Lake Storages can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_policy_data_lake_storage.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataProtection/backupVaults/vault1/backupPolicies/backupPolicy1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DataProtection` - 2025-07-01
