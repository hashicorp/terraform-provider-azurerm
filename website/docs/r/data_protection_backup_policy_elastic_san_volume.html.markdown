---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_policy_elastic_san_volume"
description: |-
  Manages a Data Protection (Backup Vault) Backup Policy for Elastic SAN Volumes.
---

# azurerm_data_protection_backup_policy_elastic_san_volume

Manages a Data Protection (Backup Vault) Backup Policy for Elastic SAN Volumes.

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

resource "azurerm_data_protection_backup_policy_elastic_san_volume" "example" {
  name                            = "example-data-protection-backup-policy-elastic-san-volume"
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.example.id
  daily_backup_schedule           = "R/2026-02-08T13:00:00+00:00/P1D"
  default_retention_duration      = "P7D"
  daily_retention_duration        = "P14D"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Data Protection Backup Policy for Elastic SAN Volumes. Changing this forces a new resource to be created.

-> **Note:** The `name` must be between 3 and 150 characters, contain only letters, numbers, and hyphens, and start with a letter.

* `data_protection_backup_vault_id` - (Required) The ID of the Data Protection Backup Vault where the policy should exist. Changing this forces a new resource to be created.

* `default_retention_duration` - (Required) The ISO 8601 duration for retaining backups under the default retention rule. Changing this forces a new resource to be created.

* `daily_backup_schedule` - (Optional) The ISO 8601 repeating time interval for daily backups. Changing this forces a new resource to be created.

-> **Note:** The interval must start with `R/`, contain an ISO 8601 timestamp, and end with `/P1D`, for example `R/2026-02-08T13:00:00+00:00/P1D`.

~> **Note:** Exactly one of `daily_backup_schedule` and `weekly_backup_schedule` must be specified.

* `daily_retention_duration` - (Optional) The ISO 8601 duration for retaining the first backup of each day. Changing this forces a new resource to be created.

~> **Note:** The `daily_retention_duration` argument cannot be specified with `weekly_backup_schedule`.

* `monthly_retention_duration` - (Optional) The ISO 8601 duration for retaining the first backup of each month. Changing this forces a new resource to be created.

~> **Note:** The `monthly_retention_duration` argument cannot be specified with `daily_backup_schedule`.

* `weekly_backup_schedule` - (Optional) The ISO 8601 repeating time interval for weekly backups. Changing this forces a new resource to be created.

-> **Note:** The interval must start with `R/`, contain an ISO 8601 timestamp, and end with `/P1W`, for example `R/2026-02-08T13:00:00+00:00/P1W`.

* `weekly_retention_duration` - (Optional) The ISO 8601 duration for retaining the first backup of each week. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Protection Backup Policy for Elastic SAN Volumes.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Protection Backup Policy for Elastic SAN Volumes.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Protection Backup Policy for Elastic SAN Volumes.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Protection Backup Policy for Elastic SAN Volumes.

## Import

A Data Protection Backup Policy for Elastic SAN Volumes can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_policy_elastic_san_volume.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DataProtection/backupVaults/backupVault1/backupPolicies/backupPolicy1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DataProtection` - 2026-06-01
