---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_instance_mysql_flexible_server"
description: |-
  Manages a Backup Instance to back up MySQL Flexible Server.
---

# azurerm_data_protection_backup_instance_mysql_flexible_server

Manages a Backup Instance to back up MySQL Flexible Server.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mysql_flexible_server" "example" {
  name                   = "example-mysqlfs"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  version                = "8.0.21"
  sku_name               = "B_Standard_B1ms"
  zone                   = "1"
}

resource "azurerm_data_protection_backup_vault" "example" {
  name                = "example-backupvault"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
  soft_delete         = "Off"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_resource_group.example.id
  role_definition_name = "Reader"
  principal_id         = azurerm_data_protection_backup_vault.example.identity.0.principal_id
}

resource "azurerm_role_assignment" "example2" {
  scope                = azurerm_mysql_flexible_server.example.id
  role_definition_name = "MySQL Backup And Export Operator"
  principal_id         = azurerm_data_protection_backup_vault.example.identity.0.principal_id
}

resource "azurerm_data_protection_backup_policy_mysql_flexible_server" "example" {
  name                            = "example-dp"
  vault_id                        = azurerm_data_protection_backup_vault.example.id
  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]

  default_retention_rule {
    life_cycle {
      duration        = "P4M"
      data_store_type = "VaultStore"
    }
  }

  depends_on = [azurerm_role_assignment.example, azurerm_role_assignment.example2]
}

resource "azurerm_data_protection_backup_instance_mysql_flexible_server" "example" {
  name             = "example-dbi"
  location         = azurerm_resource_group.example.location
  vault_id         = azurerm_data_protection_backup_vault.example.id
  server_id        = azurerm_mysql_flexible_server.example.id
  backup_policy_id = azurerm_data_protection_backup_policy_mysql_flexible_server.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Backup Instance for the MySQL Flexible Server. Changing this forces a new resource to be created.

* `location` - (Required) The location of the source database. Changing this forces a new resource to be created.

* `backup_policy_id` - (Required) The ID of the Backup Policy.

* `server_id` - (Required) The ID of the source server. Changing this forces a new resource to be created.

* `vault_id` - (Required) The ID of the Backup Vault within which the MySQL Flexible Server Backup Instance should exist. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Backup Instance MySQL Flexible Server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Backup Instance MySQL Flexible Server.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Instance MySQL Flexible Server.
* `update` - (Defaults to 1 hour) Used when updating the Backup Instance MySQL Flexible Server.
* `delete` - (Defaults to 1 hour) Used when deleting the Backup Instance MySQL Flexible Server.

## Import

Backup Instance MySQL Flexible Servers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_instance_mysql_flexible_server.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataProtection/backupVaults/vault1/backupInstances/backupInstance1
```
