---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_server"
description: |-
  Manages a Microsoft SQL Azure Database Server.

---

# azurerm_mssql_server

Manages a Microsoft SQL Azure Database Server.

~> **Note:** All arguments including the administrator login and password will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "database-rg"
  location = "West Europe"
}

resource "azurerm_mssql_server" "example" {
  name                         = "mssqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"
  minimum_tls_version          = "1.2"

  azuread_administrator {
    login_username = "AzureAD Admin"
    object_id      = "00000000-0000-0000-0000-000000000000"
  }

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Microsoft SQL Server. This needs to be globally unique within Azure. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Microsoft SQL Server. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `version` - (Required) The version for the new server. Valid values are: 2.0 (for v11 server) and 12.0 (for v12 server). Changing this forces a new resource to be created.

* `administrator_login` - (Optional) The administrator login name for the new server. Required unless `azuread_authentication_only` in the `azuread_administrator` block is `true`. When omitted, Azure will generate a default username which cannot be subsequently changed. Changing this forces a new resource to be created.

* `administrator_login_password` - (Optional) The password associated with the `administrator_login` user. Needs to comply with Azure's [Password Policy](https://msdn.microsoft.com/library/ms161959.aspx). Required unless `azuread_authentication_only` in the `azuread_administrator` block is `true`.

* `azuread_administrator` - (Optional) An `azuread_administrator` block as defined below.

* `connection_policy` - (Optional) The connection policy the server will use. Possible values are `Default`, `Proxy`, and `Redirect`. Defaults to `Default`.

* `identity` - (Optional) An `identity` block as defined below.

* `minimum_tls_version` - (Optional) The Minimum TLS Version for all SQL Database and SQL Data Warehouse databases associated with the server. Valid values are: `1.0`, `1.1` , `1.2` and `Disabled`. Defaults to `1.2`.

~> **NOTE:** The `minimum_tls_version` is set to `Disabled` means all TLS versions are allowed. After you enforce a version of `minimum_tls_version`, it's not possible to revert to `Disabled`.

* `public_network_access_enabled` - (Optional) Whether public network access is allowed for this server. Defaults to `true`.

* `outbound_network_restriction_enabled` - (Optional) Whether outbound network traffic is restricted for this server. Defaults to `false`.

* `primary_user_assigned_identity_id` - (Optional) Specifies the primary user managed identity id. Required if `type` is `UserAssigned` and should be combined with `identity_ids`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this SQL Server. Possible values are `SystemAssigned`, `UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this SQL Server.

~> **NOTE:** This is required when `type` is set to `UserAssigned`

~> **NOTE:** When `type` is set to `SystemAssigned`, the assigned `principal_id` and `tenant_id` can be retrieved after the Microsoft SQL Server has been created. More details are available below.

---

An `azuread_administrator` block supports the following:

* `login_username` - (Required) The login username of the Azure AD Administrator of this SQL Server.

* `object_id` - (Required) The object id of the Azure AD Administrator of this SQL Server.

* `tenant_id` - (Optional) The tenant id of the Azure AD Administrator of this SQL Server.

* `azuread_authentication_only` - (Optional) Specifies whether only AD Users and administrators (like `azuread_administrator.0.login_username`) can be used to login, or also local database users (like `administrator_login`). When `true`, the `administrator_login` and `administrator_login_password` properties can be omitted.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - the Microsoft SQL Server ID.

* `fully_qualified_domain_name` - The fully qualified domain name of the Azure SQL Server (e.g. myServerName.database.windows.net)

* `restorable_dropped_database_ids` - A list of dropped restorable database IDs on the server.

---

A `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this SQL Server.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this SQL Server.

-> You can access the Principal ID via `azurerm_mssql_server.example.identity.0.principal_id` and the Tenant ID via `azurerm_mssql_server.example.identity.0.tenant_id`

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Microsoft SQL Server.
* `update` - (Defaults to 60 minutes) Used when updating the Microsoft SQL Server.
* `read` - (Defaults to 5 minutes) Used when retrieving the Microsoft SQL Server.
* `delete` - (Defaults to 60 minutes) Used when deleting the Microsoft SQL Server.

## Import

SQL Servers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_server.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver
```
