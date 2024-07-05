---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_instance_blob_storage"
description: |-
  Manages a Backup Instance Blob Storage.
---

# azurerm_data_protection_backup_instance_blob_storage

Manages a Backup Instance Blob Storage.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "storageaccountname"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
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

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Storage Account Backup Contributor"
  principal_id         = azurerm_data_protection_backup_vault.example.identity[0].principal_id
}

resource "azurerm_data_protection_backup_policy_blob_storage" "example" {
  name               = "example-backup-policy"
  vault_id           = azurerm_data_protection_backup_vault.example.id
  retention_duration = "P30D"
}

resource "azurerm_data_protection_backup_instance_blob_storage" "example" {
  name               = "example-backup-instance"
  vault_id           = azurerm_data_protection_backup_vault.example.id
  location           = azurerm_resource_group.example.location
  storage_account_id = azurerm_storage_account.example.id
  backup_policy_id   = azurerm_data_protection_backup_policy_blob_storage.example.id

  depends_on = [azurerm_role_assignment.example]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Backup Instance Blob Storage. Changing this forces a new Backup Instance Blob Storage to be created.

* `location` - (Required) The location of the source Storage Account. Changing this forces a new Backup Instance Blob Storage to be created.

* `vault_id` - (Required) The ID of the Backup Vault within which the Backup Instance Blob Storage should exist. Changing this forces a new Backup Instance Blob Storage to be created.

* `storage_account_id` - (Required) The ID of the source Storage Account. Changing this forces a new Backup Instance Blob Storage to be created.

* `backup_policy_id` - (Required) The ID of the Backup Policy.

* `storage_account_container_names` - (Optional) The list of the container names of the source Storage Account.

-> **Note:** The `storage_account_container_names` should be specified in the vaulted backup policy/operational and vaulted hybrid backup policy. Removing the `storage_account_container_names` will force a new resource to be created since it can't be removed once specified.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Backup Instance Blob Storage.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Backup Instance Blob Storage.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Instance Blob Storage.
* `update` - (Defaults to 30 minutes) Used when updating the Backup Instance Blob Storage.
* `delete` - (Defaults to 30 minutes) Used when deleting the Backup Instance Blob Storage.

## Import

Backup Instance Blob Storages can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_instance_blob_storage.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataProtection/backupVaults/vault1/backupInstances/backupInstance1
```
