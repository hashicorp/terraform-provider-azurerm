---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_database_extended_auditing_policy"
description: |-
  Manages a MS SQL Database Extended Auditing Policy.
---

# azurerm_mssql_database_extended_auditing_policy

Manages a MS SQL Database Extended Auditing Policy.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mssql_server" "example" {
  name                         = "example-sqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "AdminPassword123!"
}

resource "azurerm_mssql_database" "example" {
  name      = "example-db"
  server_id = azurerm_mssql_server.example.id
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database_extended_auditing_policy" "example" {
  database_id                             = azurerm_mssql_database.example.id
  storage_endpoint                        = azurerm_storage_account.example.primary_blob_endpoint
  storage_account_access_key              = azurerm_storage_account.example.primary_access_key
  storage_account_access_key_is_secondary = false
  retention_in_days                       = 6
}
```

## Arguments Reference

The following arguments are supported:

* `database_id` - (Required) The ID of the SQL database to set the extended auditing policy. Changing this forces a new resource to be created.

* `enabled` - (Optional) Whether to enable the extended auditing policy. Possible values are `true` and `false`. Defaults to `true`.

-> **Note:** If `enabled` is `true`, `storage_endpoint` or `log_monitoring_enabled` are required.

* `storage_endpoint` - (Optional) The blob storage endpoint (e.g. <https://example.blob.core.windows.net>). This blob storage will hold all extended auditing logs.

* `retention_in_days` - (Optional) The number of days to retain logs for in the storage account. Defaults to `0`.

* `storage_account_access_key` - (Optional) The access key to use for the auditing storage account.

* `storage_account_access_key_is_secondary` - (Optional) Is `storage_account_access_key` value the storage's secondary key?

* `log_monitoring_enabled` - (Optional) Enable audit events to Azure Monitor? Defaults to `true`.

~> **Note:** To enable sending audit events to Log Analytics, please refer to the example which can be found in [the `./examples/sql-azure/sql_auditing_log_analytics` directory within the GitHub Repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/sql-azure/sql_server_auditing_log_analytics).  To enable sending server audit events to Log Analytics, please enable the master database to send audit events to Log Analytics.
To enable audit events to Eventhub, please refer to the example which can be found in [the `./examples/sql-azure/sql_auditing_eventhub` directory within the GitHub Repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/sql-azure/sql_auditing_eventhub).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the MS SQL Database Extended Auditing Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the MS SQL Database Extended Auditing Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the MS SQL Database Extended Auditing Policy.
* `update` - (Defaults to 30 minutes) Used when updating the MS SQL Database Extended Auditing Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the MS SQL Database Extended Auditing Policy.

## Import

MS SQL Database Extended Auditing Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_database_extended_auditing_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Sql/servers/sqlServer1/databases/db1/extendedAuditingSettings/default
```
