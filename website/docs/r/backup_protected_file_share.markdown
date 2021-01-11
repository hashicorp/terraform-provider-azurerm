---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_backup_protected_file_share"
description: |-
  Manages an Azure Backup Protected File Share.
---

# azurerm_backup_protected_file_share

Manages an Azure Backup Protected File Share to enable backups for file shares within an Azure Storage Account

-> **NOTE:** Azure Backup for Azure File Shares is currently in public preview. During the preview, the service is subject to additional limitations and unsupported backup scenarios. [Read More](https://docs.microsoft.com/en-us/azure/backup/backup-azure-files#limitations-for-azure-file-share-backup-during-preview)

-> **NOTE** Azure Backup for Azure File Shares does not support Soft Delete at this time. Deleting this resource will also delete all associated backup data. Please exercise caution. Consider using [`prevent_destroy`](https://www.terraform.io/docs/configuration/resources.html#prevent_destroy) to guard against accidental deletion.

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "tfex-recovery_vault"
  location = "West US"
}

resource "azurerm_recovery_services_vault" "vault" {
  name                = "tfex-recovery-vault"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  sku                 = "Standard"
}

resource "azurerm_storage_account" "sa" {
  name                     = "examplesa"
  location                 = azurerm_resource_group.rg.location
  resource_group_name      = azurerm_resource_group.rg.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "example" {
  name                 = "example-share"
  storage_account_name = azurerm_storage_account.sa.name
}

resource "azurerm_backup_container_storage_account" "protection-container" {
  resource_group_name = azurerm_resource_group.rg.name
  recovery_vault_name = azurerm_recovery_services_vault.vault.name
  storage_account_id  = azurerm_storage_account.sa.id
}

resource "azurerm_backup_policy_file_share" "example" {
  name                = "tfex-recovery-vault-policy"
  resource_group_name = azurerm_resource_group.rg.name
  recovery_vault_name = azurerm_recovery_services_vault.vault.name

  backup {
    frequency = "Daily"
    time      = "23:00"
  }
  
  retention_daily {
    count = 10
  }
}

resource "azurerm_backup_protected_file_share" "share1" {
  resource_group_name       = azurerm_resource_group.rg.name
  recovery_vault_name       = azurerm_recovery_services_vault.vault.name
  source_storage_account_id = azurerm_backup_container_storage_account.protection-container.storage_account_id
  source_file_share_name    = azurerm_storage_share.example.name
  backup_policy_id          = azurerm_backup_policy_file_share.example.id
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group in which to create the Azure Backup Protected File Share. Changing this forces a new resource to be created.

* `recovery_vault_name` - (Required) Specifies the name of the Recovery Services Vault to use. Changing this forces a new resource to be created.

* `source_storage_account_id` - (Required) Specifies the ID of the storage account of the file share to backup. Changing this forces a new resource to be created.

-> **NOTE** The storage account must already be registered with the recovery vault in order to backup shares within the account. You can use the `azurerm_backup_container_storage_account` resource or the [Register-AzRecoveryServicesBackupContainer PowerShell cmdlet](https://docs.microsoft.com/en-us/powershell/module/az.recoveryservices/register-azrecoveryservicesbackupcontainer?view=azps-3.2.0) to register a storage account with a vault.

* `source_file_share_name` - (Required) Specifies the name of the file share to backup. Changing this forces a new resource to be created.

* `backup_policy_id` - (Required) Specifies the ID of the backup policy to use. The policy must be an Azure File Share backup policy. Other types are not supported.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Backup File Share.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 80 minutes) Used when creating the Backup File Share.
* `update` - (Defaults to 80 minutes) Used when updating the Backup File Share.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup File Share.
* `delete` - (Defaults to 80 minutes) Used when deleting the Backup File Share.

## Import

Azure Backup Protected File Shares can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_backup_protected_file_share.item1 "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/example-recovery-vault/backupFabrics/Azure/protectionContainers/StorageContainer;storage;group2;example-storage-account/protectedItems/AzureFileShare;example-share"
```

Note the ID requires quoting as there are semicolons
