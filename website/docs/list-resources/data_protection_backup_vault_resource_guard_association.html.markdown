---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_vault_resource_guard_association"
description: |-
  Lists Backup Vault Resource Guard Association resources.
---

# List resource: azurerm_data_protection_backup_vault_resource_guard_association

Lists Backup Vault Resource Guard Association resources.

## Example Usage

```hcl
list "azurerm_data_protection_backup_vault_resource_guard_association" "example" {
  provider = azurerm
  config {
    data_protection_backup_vault_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DataProtection/backupVaults/backupVault1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `data_protection_backup_vault_id` - (Required) The ID of the Data Protection Backup Vault to query.