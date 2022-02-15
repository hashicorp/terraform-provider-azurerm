---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_backup_protected_file_share"
description: |-
  Gets information about an existing Azure Backup Protected File Share.
---

# Data Source: azurerm_backup_protected_file_share

Use this data source to access information about an existing Backup Protected File Share.

## Example Usage

```hcl
data "azurerm_backup_protected_file_share" "share1" {
  source_file_share_name    = "example-share"
  source_storage_account_id = azurerm_backup_container_storage_account.protection-container.storage_account_id
  recovery_vault_name       = "recovery_vault"
  resource_group_name       = "resource_group"
}
```

## Argument Reference

The following arguments are supported:

* `source_file_share_name` - Specifies the name of the file share that is backed up.

* `resource_group_name` - The name of the resource group in which the file share backup resides.

* `recovery_vault_name` - Specifies the name of the Recovery Services Vault.

* `source_storage_account_id` - Specifies the ID of the storage account of the file share to that is backed up.

## Attributes Reference

In addition to the arguments listed above - the following attributes are exported: 

* `id` - The ID of the file share backup.

* `backup_policy_id` - Specifies the ID of the backup policy used.
