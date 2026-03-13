---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_instance_data_lake_storage"
description: |-
  Manages a Backup Instance to back up Azure Data Lake Storage.
---

# azurerm_data_protection_backup_instance_data_lake_storage

Manages a Backup Instance to back up Azure Data Lake Storage.

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
  is_hns_enabled           = true
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

resource "azurerm_data_protection_backup_policy_data_lake_storage" "example" {
  name                            = "example-backup-policy"
  vault_id                        = azurerm_data_protection_backup_vault.example.id
  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]

  default_retention_rule {
    life_cycle {
      duration        = "P4M"
      data_store_type = "VaultStore"
    }
  }
}

resource "azurerm_data_protection_backup_instance_data_lake_storage" "example" {
  name               = "example-backup-instance"
  location           = azurerm_resource_group.example.location
  vault_id           = azurerm_data_protection_backup_vault.example.id
  storage_account_id = azurerm_storage_account.example.id
  backup_policy_id   = azurerm_data_protection_backup_policy_data_lake_storage.example.id

  depends_on = [azurerm_role_assignment.example]
}
```

~> **Note:** The Storage Account used must have hierarchical namespace enabled (`is_hns_enabled = true`) to support Azure Data Lake Storage backup.

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Backup Instance Data Lake Storage. Changing this forces a new resource to be created.

* `location` - (Required) The location of the source Storage Account. Changing this forces a new resource to be created.

* `backup_policy_id` - (Required) The ID of the Backup Policy.

* `storage_account_id` - (Required) The ID of the source Storage Account. Changing this forces a new resource to be created.

* `vault_id` - (Required) The ID of the Backup Vault within which the Backup Instance Data Lake Storage should exist. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Backup Instance Data Lake Storage.

* `protection_state` - The protection state of the Backup Instance Data Lake Storage.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Backup Instance Data Lake Storage.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Instance Data Lake Storage.
* `update` - (Defaults to 1 hour) Used when updating the Backup Instance Data Lake Storage.
* `delete` - (Defaults to 1 hour) Used when deleting the Backup Instance Data Lake Storage.

## Import

Backup Instance Data Lake Storages can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_instance_data_lake_storage.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataProtection/backupVaults/vault1/backupInstances/backupInstance1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DataProtection` - 2025-09-01
