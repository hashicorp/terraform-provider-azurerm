---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_policy_elastic_san_volume"
description: |-
  Lists Data Protection (Backup Vault) Backup Policy resources for Elastic SAN Volumes.
---

# List resource: azurerm_data_protection_backup_policy_elastic_san_volume

Lists Data Protection (Backup Vault) Backup Policy resources for Elastic SAN Volumes.

## Example Usage

### List all Data Protection Backup Policies for Elastic SAN Volumes in a Data Protection Backup Vault

```hcl
list "azurerm_data_protection_backup_policy_elastic_san_volume" "example" {
  provider = azurerm

  config {
    data_protection_backup_vault_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DataProtection/backupVaults/backupVault1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `data_protection_backup_vault_id` - (Required) The ID of the Data Protection Backup Vault to query.
