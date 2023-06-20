---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_instance_disk"
description: |-
  Manages a Backup Instance to back up Disk.
---

# azurerm_data_protection_backup_instance_disk

Manages a Backup Instance to back up Disk.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_managed_disk" "example" {
  name                 = "example-disk"
  location             = azurerm_resource_group.example.location
  resource_group_name  = azurerm_resource_group.example.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"
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

resource "azurerm_role_assignment" "example1" {
  scope                = azurerm_resource_group.example.id
  role_definition_name = "Disk Snapshot Contributor"
  principal_id         = azurerm_data_protection_backup_vault.example.identity[0].principal_id
}

resource "azurerm_role_assignment" "example2" {
  scope                = azurerm_managed_disk.example.id
  role_definition_name = "Disk Backup Reader"
  principal_id         = azurerm_data_protection_backup_vault.example.identity[0].principal_id
}


resource "azurerm_data_protection_backup_policy_disk" "example" {
  name     = "example-backup-policy"
  vault_id = azurerm_data_protection_backup_vault.example.id

  backup_repeating_time_intervals = ["R/2021-05-19T06:33:16+00:00/PT4H"]
  default_retention_duration      = "P7D"
}

resource "azurerm_data_protection_backup_instance_disk" "example" {
  name                         = "example-backup-instance"
  location                     = azurerm_data_protection_backup_vault.example.location
  vault_id                     = azurerm_data_protection_backup_vault.example.id
  disk_id                      = azurerm_managed_disk.example.id
  snapshot_resource_group_name = azurerm_resource_group.example.name
  backup_policy_id             = azurerm_data_protection_backup_policy_disk.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Backup Instance Disk. Changing this forces a new Backup Instance Disk to be created.

* `location` - (Required) The Azure Region where the Backup Instance Disk should exist. Changing this forces a new Backup Instance Disk to be created.

* `vault_id` - (Required) The ID of the Backup Vault within which the Backup Instance Disk should exist. Changing this forces a new Backup Instance Disk to be created.

* `disk_id` - (Required) The ID of the source Disk. Changing this forces a new Backup Instance Disk to be created.

* `snapshot_resource_group_name` - (Required) The name of the Resource Group where snapshots are stored. Changing this forces a new Backup Instance Disk to be created.

* `backup_policy_id` - (Required) The ID of the Backup Policy.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Backup Instance Disk.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Backup Instance Disk.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Instance Disk.
* `update` - (Defaults to 30 minutes) Used when updating the Backup Instance Disk.
* `delete` - (Defaults to 30 minutes) Used when deleting the Backup Instance Disk.

## Import

Backup Instance Disks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_instance_disk.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataProtection/backupVaults/vault1/backupInstances/backupInstance1
```
