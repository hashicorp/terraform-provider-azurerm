---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_instance_cosmosdb_account"
description: |-
  Lists Data Protection (Backup Vault) Backup Instance resources for Cosmos DB Database Accounts in a Backup Vault.
---

# List resource: azurerm_data_protection_backup_instance_cosmosdb_account

Lists Data Protection (Backup Vault) Backup Instance resources for Cosmos DB Database Accounts in a Backup Vault.

## Example Usage

### List Backup Instances for Cosmos DB Database Accounts in a Backup Vault

```hcl
list "azurerm_data_protection_backup_instance_cosmosdb_account" "example" {
  provider = azurerm
  config {
    data_protection_backup_vault_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DataProtection/backupVaults/backupVault1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `data_protection_backup_vault_id` - (Required) The ID of the Data Protection Backup Vault to query.
