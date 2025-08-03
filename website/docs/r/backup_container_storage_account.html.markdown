---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_backup_container_storage_account"
description: |-
    Manages a storage account container in an Azure Recovery Vault
---

# azurerm_backup_container_storage_account

Manages registration of a storage account with Azure Backup. Storage accounts must be registered with an Azure Recovery Vault in order to backup file shares within the storage account. Registering a storage account with a vault creates what is known as a protection container within Azure Recovery Services. Once the container is created, Azure file shares within the storage account can be backed up using the `azurerm_backup_protected_file_share` resource.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tfex-network-mapping-primary"
  location = "West Europe"
}

resource "azurerm_recovery_services_vault" "vault" {
  name                = "example-recovery-vault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "azurerm_storage_account" "sa" {
  name                     = "examplesa"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_backup_container_storage_account" "container" {
  resource_group_name = azurerm_resource_group.example.name
  recovery_vault_name = azurerm_recovery_services_vault.vault.name
  storage_account_id  = azurerm_storage_account.sa.id
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) Name of the resource group where the vault is located. Changing this forces a new resource to be created.

* `recovery_vault_name` - (Required) The name of the vault where the storage account will be registered. Changing this forces a new resource to be created.

* `storage_account_id` - (Required) The ID of the Storage Account to be registered Changing this forces a new resource to be created.

-> **Note:** Azure Backup places a Resource Lock on the storage account that will cause deletion to fail until the account is unregistered from Azure Backup

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the Backup Storage Account Container.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Backup Storage Account Container.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Storage Account Container.
* `delete` - (Defaults to 30 minutes) Used when deleting the Backup Storage Account Container.

## Import

Backup Storage Account Containers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_backup_container_storage_account.mycontainer "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/backupFabrics/Azure/protectionContainers/StorageContainer;storage;storage-rg-name;storage-account"
```

Note the ID requires quoting as there are semicolons
