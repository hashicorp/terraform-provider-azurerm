---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_instance_elastic_san_volume_group"
description: |-
  Manages a Data Protection Backup Instance for an Elastic SAN Volume Group.
---

# azurerm_data_protection_backup_instance_elastic_san_volume_group

Manages a Data Protection Backup Instance for an Elastic SAN Volume Group.

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

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_elastic_san" "example" {
  name                = "example-elastic-san"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  base_size_in_tib    = 1

  sku {
    name = "Premium_LRS"
  }
}

resource "azurerm_elastic_san_volume_group" "example" {
  name           = "example-elastic-san-volume-group"
  elastic_san_id = azurerm_elastic_san.example.id
}

resource "azurerm_elastic_san_volume" "example" {
  name            = "example-elastic-san-volume"
  volume_group_id = azurerm_elastic_san_volume_group.example.id
  size_in_gib     = 1
}

resource "azurerm_role_assignment" "snapshot_exporter" {
  scope                = azurerm_elastic_san.example.id
  role_definition_name = "Elastic SAN Snapshot Exporter"
  principal_id         = azurerm_data_protection_backup_vault.example.identity[0].principal_id
}

resource "azurerm_role_assignment" "snapshot_contributor" {
  scope                = azurerm_resource_group.example.id
  role_definition_name = "Disk Snapshot Contributor"
  principal_id         = azurerm_data_protection_backup_vault.example.identity[0].principal_id
}

resource "azurerm_data_protection_backup_policy_elastic_san_volume_group" "example" {
  name                            = "example-data-protection-backup-policy-elastic-san-volume-group"
  vault_id                        = azurerm_data_protection_backup_vault.example.id
  backup_repeating_time_intervals = ["R/2024-02-08T13:00:00+00:00/P1D"]

  default_retention_rule {
    life_cycle {
      data_store_type = "OperationalStore"
      duration        = "P7D"
    }
  }
}

resource "azurerm_data_protection_backup_instance_elastic_san_volume_group" "example" {
  name                         = "example-data-protection-backup-instance-elastic-san-volume-group"
  location                     = azurerm_resource_group.example.location
  vault_id                     = azurerm_data_protection_backup_vault.example.id
  backup_policy_id             = azurerm_data_protection_backup_policy_elastic_san_volume_group.example.id
  elastic_san_volume_group_id  = azurerm_elastic_san_volume_group.example.id
  snapshot_resource_group_name = azurerm_resource_group.example.name
  volume_name                  = azurerm_elastic_san_volume.example.name

  depends_on = [
    azurerm_role_assignment.snapshot_exporter,
    azurerm_role_assignment.snapshot_contributor,
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Data Protection Backup Instance for the Elastic SAN Volume Group. Changing this forces a new resource to be created.

* `vault_id` - (Required) The ID of the Data Protection Backup Vault where the Backup Instance should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region of the Elastic SAN Volume Group. Changing this forces a new resource to be created.

* `backup_policy_id` - (Required) The ID of the Data Protection Backup Policy for Elastic SAN Volume Groups.

* `elastic_san_volume_group_id` - (Required) The ID of the Elastic SAN Volume Group to protect. Changing this forces a new resource to be created.

* `snapshot_resource_group_name` - (Required) The name of the Resource Group where volume snapshots are stored. Changing this forces a new resource to be created.

* `volume_name` - (Required) The name of the Elastic SAN Volume within the Volume Group to protect.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Protection Backup Instance for the Elastic SAN Volume Group.

* `protection_state` - The protection state of the Elastic SAN Volume Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Data Protection Backup Instance for the Elastic SAN Volume Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Protection Backup Instance for the Elastic SAN Volume Group.
* `update` - (Defaults to 1 hour) Used when updating the Data Protection Backup Instance for the Elastic SAN Volume Group.
* `delete` - (Defaults to 1 hour) Used when deleting the Data Protection Backup Instance for the Elastic SAN Volume Group.

## Import

A Data Protection Backup Instance for an Elastic SAN Volume Group can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_instance_elastic_san_volume_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DataProtection/backupVaults/backupVault1/backupInstances/backupInstance1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DataProtection` - 2026-06-01
