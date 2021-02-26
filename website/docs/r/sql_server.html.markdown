---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_server"
description: |-
  Manages a Microsoft SQL Azure Database Server.

---

# azurerm_sql_server

Manages a Microsoft SQL Azure Database Server.

~> **Note:** All arguments including the administrator login and password will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "database-rg"
  location = "West US"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_sql_server" "example" {
  name                         = "mssqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"

  extended_auditing_policy {
    storage_endpoint                        = azurerm_storage_account.example.primary_blob_endpoint
    storage_account_access_key              = azurerm_storage_account.example.primary_access_key
    storage_account_access_key_is_secondary = true
    retention_in_days                       = 6
  }

  tags = {
    environment = "production"
  }
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Microsoft SQL Server. This needs to be globally unique within Azure.

* `resource_group_name` - (Required) The name of the resource group in which to create the Microsoft SQL Server.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `version` - (Required) The version for the new server. Valid values are: 2.0 (for v11 server) and 12.0 (for v12 server).

* `administrator_login` - (Required) The administrator login name for the new server. Changing this forces a new resource to be created.

* `administrator_login_password` - (Required) The password associated with the `administrator_login` user. Needs to comply with Azure's [Password Policy](https://msdn.microsoft.com/library/ms161959.aspx)

* `connection_policy` - (Optional) The connection policy the server will use. Possible values are `Default`, `Proxy`, and `Redirect`. Defaults to `Default`.

* `identity` - (Optional) An `identity` block as defined below.

* `extended_auditing_policy` - (Optional) A `extended_auditing_policy` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the Microsoft SQL Server. At this time the only allowed value is `SystemAssigned`.

~> **NOTE:** The assigned `principal_id` and `tenant_id` can be retrieved after the identity `type` has been set to `SystemAssigned` and the Microsoft SQL Server has been created. More details are available below.

## Attributes Reference

The following attributes are exported:

* `id` - The Microsoft SQL Server ID.
* `fully_qualified_domain_name` - The fully qualified domain name of the Azure SQL Server (e.g. myServerName.database.windows.net)

---

`identity` exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this SQL Server.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this SQL Server.

-> You can access the Principal ID via `${azurerm_mssql_server.example.identity.0.principal_id}` and the Tenant ID via `${azurerm_mssql_server.example.identity.0.tenant_id}`

---

A `extended_auditing_policy` block supports the following:

* `storage_account_access_key` - (Optional)  Specifies the access key to use for the auditing storage account.

* `storage_endpoint` - (Optional) Specifies the blob storage endpoint (e.g. https://MyAccount.blob.core.windows.net).

* `storage_account_access_key_is_secondary` - (Optional) Specifies whether `storage_account_access_key` value is the storage's secondary key.

* `retention_in_days` - (Optional) Specifies the number of days to retain logs for in the storage account.

* `monitor_enabled` - (Optional) Enable audit events to Azure Monitor? To enable server audit events to Azure Monitor, please enable its master database audit events to Azure Monitor.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Microsoft SQL Server.
* `update` - (Defaults to 60 minutes) Used when updating the Microsoft SQL Server.
* `read` - (Defaults to 5 minutes) Used when retrieving the Microsoft SQL Server.
* `delete` - (Defaults to 60 minutes) Used when deleting the Microsoft SQL Server.

## Import

SQL Servers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sql_server.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver
```
