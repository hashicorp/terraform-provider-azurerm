---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_instance_cosmosdb_database_account"
description: |-
  Manages a Data Protection Backup Instance for a Cosmos DB Database Account.
---

# azurerm_data_protection_backup_instance_cosmosdb_database_account

Manages a Data Protection Backup Instance for a Cosmos DB Database Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_cosmosdb_account" "example" {
  name                = "example-cosmosdb-account"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level = "Session"
  }

  backup {
    type = "Continuous"
    tier = "Continuous7Days"
  }

  geo_location {
    location          = azurerm_resource_group.example.location
    failover_priority = 0
  }
}

resource "azurerm_data_protection_backup_vault" "example" {
  name                = "example-data-protection-backup-vault"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_data_protection_backup_policy_cosmosdb_database_account" "example" {
  name                            = "example-data-protection-backup-policy-cosmosdb-database-account"
  vault_id                        = azurerm_data_protection_backup_vault.example.id
  backup_repeating_time_intervals = ["R/2026-02-08T10:00:00+00:00/P1W"]
  time_zone                       = "UTC"

  default_retention_rule {
    life_cycle {
      duration        = "P10Y"
      data_store_type = "VaultStore"
    }
  }
}

resource "azurerm_role_assignment" "reader" {
  scope                = azurerm_resource_group.example.id
  role_definition_name = "Reader"
  principal_id         = azurerm_data_protection_backup_vault.example.identity[0].principal_id
}

resource "azurerm_role_assignment" "cosmos_operator" {
  scope                = azurerm_cosmosdb_account.example.id
  role_definition_name = "Cosmos DB Operator"
  principal_id         = azurerm_data_protection_backup_vault.example.identity[0].principal_id
}

resource "azurerm_data_protection_backup_instance_cosmosdb_database_account" "example" {
  name                = "example-data-protection-backup-instance-cosmosdb-database-account"
  location            = azurerm_resource_group.example.location
  vault_id            = azurerm_data_protection_backup_vault.example.id
  backup_policy_id    = azurerm_data_protection_backup_policy_cosmosdb_database_account.example.id
  cosmosdb_account_id = azurerm_cosmosdb_account.example.id

  depends_on = [
    azurerm_role_assignment.reader,
    azurerm_role_assignment.cosmos_operator,
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Data Protection Backup Instance for the Cosmos DB Database Account. Changing this forces a new resource to be created.

* `vault_id` - (Required) The ID of the Data Protection Backup Vault where the Backup Instance should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region of the Cosmos DB Database Account. Changing this forces a new resource to be created.

* `backup_policy_id` - (Required) The ID of the Data Protection Backup Policy for Cosmos DB Database Accounts.

* `cosmosdb_account_id` - (Required) The ID of the Cosmos DB Database Account to protect. Changing this forces a new resource to be created.

~> **Note:** The Cosmos DB Database Account must use Continuous backup mode.

-> **Note:** The Backup Vault's system-assigned identity requires the `Reader` role on the source Resource Group and the `Cosmos DB Operator` role on the Cosmos DB Database Account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Protection Backup Instance for the Cosmos DB Database Account.

* `protection_state` - The protection state of the Cosmos DB Database Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Data Protection Backup Instance for the Cosmos DB Database Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Protection Backup Instance for the Cosmos DB Database Account.
* `update` - (Defaults to 1 hour) Used when updating the Data Protection Backup Instance for the Cosmos DB Database Account.
* `delete` - (Defaults to 1 hour) Used when deleting the Data Protection Backup Instance for the Cosmos DB Database Account.

## Import

A Data Protection Backup Instance for a Cosmos DB Database Account can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_instance_cosmosdb_database_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DataProtection/backupVaults/backupVault1/backupInstances/backupInstance1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DataProtection` - 2026-06-01
