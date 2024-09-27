---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_instance_postgresql"
description: |-
  Manages a Backup Instance to back up PostgreSQL.
---

# azurerm_data_protection_backup_instance_postgresql

Manages a Backup Instance to back up PostgreSQL.

-> **Note:** Before using this resource, there are some prerequisite permissions for configure backup and restore. See more details from <https://docs.microsoft.com/azure/backup/backup-azure-database-postgresql#prerequisite-permissions-for-configure-backup-and-restore>.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_postgresql_server" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku_name = "B_Gen5_2"

  storage_mb                   = 5120
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  auto_grow_enabled            = true

  administrator_login          = "psqladmin"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true
}

resource "azurerm_postgresql_firewall_rule" "example" {
  name                = "AllowAllWindowsAzureIps"
  resource_group_name = azurerm_resource_group.example.name
  server_name         = azurerm_postgresql_server.example.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "0.0.0.0"
}

resource "azurerm_postgresql_database" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  server_name         = azurerm_postgresql_server.example.name
  charset             = "UTF8"
  collation           = "English_United States.1252"
}

resource "azurerm_data_protection_backup_vault" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault" "example" {
  name                       = "example"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = ["Create", "Get"]

    secret_permissions = [
      "Set",
      "Get",
      "Delete",
      "Purge",
      "Recover"
    ]
  }

  access_policy {
    tenant_id = azurerm_data_protection_backup_vault.example.identity[0].tenant_id
    object_id = azurerm_data_protection_backup_vault.example.identity[0].principal_id

    key_permissions = ["Create", "Get"]

    secret_permissions = [
      "Set",
      "Get",
      "Delete",
      "Purge",
      "Recover"
    ]
  }
}

resource "azurerm_key_vault_secret" "example" {
  name         = "example"
  value        = "Server=${azurerm_postgresql_server.example.name}.postgres.database.azure.com;Database=${azurerm_postgresql_database.example.name};Port=5432;User Id=psqladmin@${azurerm_postgresql_server.example.name};Password=H@Sh1CoR3!;Ssl Mode=Require;"
  key_vault_id = azurerm_key_vault.example.id
}

resource "azurerm_data_protection_backup_policy_postgresql" "example" {
  name                            = "example"
  resource_group_name             = azurerm_resource_group.example.name
  vault_name                      = azurerm_data_protection_backup_vault.example.name
  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]
  default_retention_duration      = "P4M"
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_postgresql_server.example.id
  role_definition_name = "Reader"
  principal_id         = azurerm_data_protection_backup_vault.example.identity[0].principal_id
}

resource "azurerm_data_protection_backup_instance_postgresql" "example" {
  name                                    = "example"
  location                                = azurerm_resource_group.example.location
  vault_id                                = azurerm_data_protection_backup_vault.example.id
  database_id                             = azurerm_postgresql_database.example.id
  backup_policy_id                        = azurerm_data_protection_backup_policy_postgresql.example.id
  database_credential_key_vault_secret_id = azurerm_key_vault_secret.example.versionless_id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Backup Instance PostgreSQL. Changing this forces a new Backup Instance PostgreSQL to be created.

* `location` - (Required) The location of the source database. Changing this forces a new Backup Instance PostgreSQL to be created.

* `vault_id` - (Required) The ID of the Backup Vault within which the PostgreSQL Backup Instance should exist. Changing this forces a new Backup Instance PostgreSQL to be created.

* `database_id` - (Required) The ID of the source database. Changing this forces a new Backup Instance PostgreSQL to be created.

* `backup_policy_id` - (Required) The ID of the Backup Policy.

* `database_credential_key_vault_secret_id` - (Optional) The ID or versionless ID of the key vault secret which stores the connection string of the database.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Backup Instance PostgreSQL.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Backup Instance PostgreSQL.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Instance PostgreSQL.
* `update` - (Defaults to 30 minutes) Used when updating the Backup Instance PostgreSQL.
* `delete` - (Defaults to 30 minutes) Used when deleting the Backup Instance PostgreSQL.

## Import

Backup Instance PostgreSQL can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_instance_postgresql.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataProtection/backupVaults/vault1/backupInstances/backupInstance1
```
