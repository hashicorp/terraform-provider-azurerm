---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_vault_resource_guard_association"
description: |-
  Manages a Backup Vault Resource Guard Association.
---

# azurerm_data_protection_backup_vault_resource_guard_association

Manages a Backup Vault Resource Guard Association.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_data_protection_backup_vault" "example" {
  name                = "example-data-protection-backup-vault"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
  soft_delete         = "Off"
}

resource "azurerm_data_protection_resource_guard" "example" {
  name                = "example-data-protection-resource-guard"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_data_protection_backup_vault_resource_guard_association" "example" {
  data_protection_backup_vault_id  = azurerm_data_protection_backup_vault.example.id
  data_protection_resource_guard_id = azurerm_data_protection_resource_guard.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `data_protection_backup_vault_id` - (Required) The ID of the Data Protection Backup Vault. Changing this forces a new resource to be created.

* `data_protection_resource_guard_id` - (Required) The ID of the Data Protection Resource Guard. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Backup Vault Resource Guard Association.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Backup Vault Resource Guard Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Vault Resource Guard Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the Backup Vault Resource Guard Association.

## Import

A Backup Vault Resource Guard Association can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_vault_resource_guard_association.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DataProtection/backupVaults/backupVault1/backupResourceGuardProxies/DppResourceGuardProxy
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DataProtection` - 2025-07-01
