---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_database"
description: |-
  Gets information about an existing SQL database.
---

# Data Source: azurerm_mssql_database

Use this data source to access information about an existing SQL database.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mssql_server" "example" {
  name                         = "example"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

data "azurerm_mssql_database" "example" {
  name      = "example-mssql-db"
  server_id = azurerm_mssql_server.example.id
}

output "database_id" {
  value = data.azurerm_mssql_database.example.id
}
```

## Argument Reference

* `name` - The name of the MS SQL Database.

* `server_id` - The id of the MS SQL Server on which to read the database.

## Attributes Reference

* `id` - The ID of the database.

* `collation` - The collation of the database.

* `elastic_pool_id` - The id of the elastic pool containing this database.

* `enclave_type` - The type of enclave being used by the database.

* `license_type` - The license type to apply for this database.

* `max_size_gb` - The max size of the database in gigabytes.

* `read_replica_count` - The number of readonly secondary replicas associated with the database to which readonly application intent connections may be routed.

* `read_scale` - If enabled, connections that have application intent set to readonly in their connection string may be routed to a readonly secondary replica.

* `sku_name` - The name of the SKU of the database.

* `storage_account_type` - The storage account type used to store backups for this database.

* `zone_redundant` - Whether or not this database is zone redundant, which means the replicas of this database will be spread across multiple availability zones.

* `identity` - A `identity` block as defined below.

* `transparent_data_encryption_enabled` - Whether or not Transparent Data Encryption is enabled.

* `transparent_data_encryption_key_vault_key_id` - The Key Vault key URI to be used as the `Customer Managed Key`(CMK/BYOK) for the `Transparent Data Encryption`(TDE) layer.

* `transparent_data_encryption_key_automatic_rotation_enabled` - Whether or not TDE automatically rotates the encryption Key to latest version.

* `tags` -  A mapping of tags to assign to the resource.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Microsoft SQL Database.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Microsoft SQL Database.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the SQL database.
