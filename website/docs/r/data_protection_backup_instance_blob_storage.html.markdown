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
  name                                   = "example-backup-policy"
  vault_id                               = azurerm_data_protection_backup_vault.example.id
  operational_default_retention_duration = "P30D"
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

* `auto_protection_enabled` - (Optional) Whether all present and future containers of the source Storage Account should be backed up. Defaults to `false`.

!> **Note:** Enabling `auto_protection_enabled` is irreversible. Once enabled it cannot be disabled again, and `storage_account_container_names` can only be specified again by recreating the Backup Instance, which deletes its vaulted recovery points.

~> **Note:** Azure only protects new containers until the number of protected containers reaches `1000`. If the Storage Account contains more than `1000` containers, `excluded_container_name_prefixes` must be used to reduce the number of protected containers to `1000` or fewer.

* `excluded_container_name_prefixes` - (Optional) A list of container name prefixes to exclude from auto protection. Containers whose names start with one of these prefixes are not backed up. Can only be specified when `auto_protection_enabled` is set to `true`.

* `storage_account_container_names` - (Optional) The list of the container names of the source Storage Account.

-> **Note:** For a vaulted backup policy or an operational and vaulted hybrid backup policy, either `storage_account_container_names` or `auto_protection_enabled` should be specified. Removing `storage_account_container_names` forces a new resource to be created, unless `auto_protection_enabled` is set to `true` at the same time.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Backup Instance Blob Storage.

* `protection_state` - The protection state of the Backup Instance Blob Storage. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Backup Instance Blob Storage.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Instance Blob Storage.
* `update` - (Defaults to 30 minutes) Used when updating the Backup Instance Blob Storage.
* `delete` - (Defaults to 30 minutes) Used when deleting the Backup Instance Blob Storage.

## Import

Backup Instance Blob Storages can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_instance_blob_storage.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataProtection/backupVaults/vault1/backupInstances/backupInstance1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DataProtection` - 2026-03-01
