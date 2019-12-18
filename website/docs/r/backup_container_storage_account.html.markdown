---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_backup_container_storage_account"
sidebar_current: "docs-azurerm-backup-container-storage-account"
description: |-
    Manages a storage account container in an Azure Recovery Vault
---

# azurerm_backup_container_storage_account

Manages registration of a storage account with Azure Backup. Storage accounts must be registered with an Azure Recovery Vault in order to backup file shares within the storage account. Registering a storage account with a vault creates what is known as a protection container within Azure Recovery Services. Once the container is created, Azure file shares within the storage account can be backed up using the `azurerm_backup_protected_file_share` resource.

-> **NOTE:** Azure Backup for Azure File Shares is currently in public preview. During the preview, the service is subject to additional limitations and unsupported backup scenarios. [Read More](https://docs.microsoft.com/en-us/azure/backup/backup-azure-files#limitations-for-azure-file-share-backup-during-preview)

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "tfex-network-mapping-primary"
  location = "West US"
}

resource "azurerm_recovery_services_vault" "vault" {
  name                = "example-recovery-vault"
  location            = "${azurerm_resource_group.rg.location}"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  sku                 = "Standard"
}

resource "azurerm_storage_account" "sa" {
  name                     = "examplesa"
  location                 = "${azurerm_resource_group.rg.location}"
  resource_group_name      = "${azurerm_resource_group.rg.name}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_backup_container_storage_account" "container" {
  resource_group_name = "${azurerm_resource_group.rg.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.vault.name}"
  storage_account_id  = "${azurerm_storage_account.sa.id}"
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) Name of the resource group where the vault is located.

* `recovery_vault_name` - (Required) The name of the vault where the storage account will be registered.

* `storage_account_id` - (Required) Azure Resource ID of the storage account to be registered

-> **NOTE** Azure Backup places a Resource Lock on the storage account that will cause deletion to fail until the account is unregistered from Azure Backup

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

Azure Backup Storage Account Containers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_backup_container_storage_account.mycontainer "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/backupFabrics/Azure/protectionContainers/StorageContainer;storage;storage-rg-name;storage-account"
```

Note the ID requires quoting as there are semicolons
