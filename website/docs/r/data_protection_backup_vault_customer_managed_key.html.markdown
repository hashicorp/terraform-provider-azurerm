---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_vault_customer_managed_key"
description: |-
  Manages a Backup Vault Customer Managed Key.
---

# azurerm_data_protection_backup_vault_customer_managed_key

Manages a Backup Vault Customer Managed Key.

!> **Note:** It is not possible to remove the Customer Managed Key from the Backup Vault once it's been added. To remove the Customer Managed Key, the parent Data Protection Backup Vault must be deleted and recreated.

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


data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "example" {
  name                        = "example-key-vault"
  location                    = azurerm_resource_group.example.location
  resource_group_name         = azurerm_resource_group.example.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = true

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Decrypt",
      "Encrypt",
      "Delete",
      "Get",
      "List",
      "Purge",
      "UnwrapKey",
      "WrapKey",
      "Verify",
      "GetRotationPolicy"
    ]
    secret_permissions = [
      "Set",
    ]
  }

  access_policy {
    tenant_id = azurerm_data_protection_backup_vault.example.identity[0].tenant_id
    object_id = azurerm_data_protection_backup_vault.example.identity[0].principal_id

    key_permissions = [
      "Create",
      "Decrypt",
      "Encrypt",
      "Delete",
      "Get",
      "List",
      "Purge",
      "UnwrapKey",
      "WrapKey",
      "Verify",
      "GetRotationPolicy"
    ]
    secret_permissions = [
      "Set",
    ]
  }
}

resource "azurerm_key_vault_key" "example" {
  name         = "example-key"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_data_protection_backup_vault_customer_managed_key" "example" {
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.example.id
  key_vault_key_id                = azurerm_key_vault_key.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `data_protection_backup_vault_id` - (Required) The ID of the Backup Vault. Changing this forces a new resource to be created.

* `key_vault_key_id` - (Required) The ID of the Key Vault Key which should be used to Encrypt the data in this Backup Vault.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Backup Vault.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Backup Vault Customer Managed Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Vault Customer Managed Key.
* `update` - (Defaults to 30 minutes) Used when updating the Backup Vault Customer Managed Key.
* `delete` - (Defaults to 5 minutes) Used when deleting the Backup Vault Customer Managed Key.

## Import

Backup Vault Customer Managed Keys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_vault_customer_managed_key.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataProtection/backupVaults/vault1
```
