---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_policy_blob_storage"
description: |-
  Manages a Backup Policy Blob Storage.
---

# azurerm_data_protection_backup_policy_blob_storage

Manages a Backup Policy Blob Storage.

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_protection_backup_vault" "example" {
  name                = "example-backup-vault"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
}

resource "azurerm_data_protection_backup_policy_blob_storage" "example" {
  name               = "example-backup-policy"
  vault_id           = azurerm_data_protection_backup_vault.example.id
  retention_duration = "P30D"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Backup Policy Blob Storage. Changing this forces a new Backup Policy Blob Storage to be created.

* `vault_id` - (Required) The ID of the Backup Vault within which the Backup Policy Blob Storage should exist. Changing this forces a new Backup Policy Blob Storage to be created.
  
* `retention_duration` - (Required) Duration of deletion after given timespan. It should follow `ISO 8601` duration format. Changing this forces a new Backup Policy Blob Storage to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Backup Policy Blob Storage.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Backup Policy Blob Storage.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Policy Blob Storage.
* `update` - (Defaults to 30 minutes) Used when updating the Backup Policy Blob Storage.
* `delete` - (Defaults to 30 minutes) Used when deleting the Backup Policy Blob Storage.

## Import

Backup Policy Blob Storages can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_policy_blob_storage.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataProtection/backupVaults/vault1/backupPolicies/backupPolicy1
```
