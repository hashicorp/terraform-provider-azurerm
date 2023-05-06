---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_server"
description: |-
  Manages a Microsoft SQL Azure Database Server.

---

# azurerm_sql_server

Manages a Microsoft SQL Azure Database Server.

~> **Note:** The `azurerm_sql_server` resource is deprecated in version 3.0 of the AzureRM provider and will be removed in version 4.0. Please use the [`azurerm_mssql_server`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/mssql_server) resource instead.

~> **Note:** All arguments including the administrator login and password will be stored in the raw state as plain-text. [Read more about sensitive data in state](https://www.terraform.io/language/state/sensitive-data).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "database-rg"
  location = "West Europe"
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

* `administrator_login` - (Required) The administrator login name for the new server. Changing this forces a new resource to be created.

* `administrator_login_password` - (Required) The password associated with the `administrator_login` user. Needs to comply with Azure's [Password Policy](https://msdn.microsoft.com/library/ms161959.aspx)

* `connection_policy` - (Optional) The connection policy the server will use. Possible values are `Default`, `Proxy`, and `Redirect`. Defaults to `Default`.

* `identity` - (Optional) An `identity` block as defined below.

* `threat_detection_policy` - (Optional) Threat detection policy configuration. The `threat_detection_policy` block supports fields documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this SQL Server. The only possible value is `SystemAssigned`.

~> **NOTE:** The assigned `principal_id` and `tenant_id` can be retrieved after the identity `type` has been set to `SystemAssigned` and the Microsoft SQL Server has been created. More details are available below.

---

The `threat_detection_policy` block supports the following:

* `state` - (Optional) The State of the Policy. Possible values are `Disabled`, `Enabled` and `New`.
* `disabled_alerts` - (Optional) Specifies a list of alerts which should be disabled. Possible values include `Access_Anomaly`, `Data_Exfiltration`, `Sql_Injection`, `Sql_Injection_Vulnerability` and `Unsafe_Action"`,.
* `email_account_admins` - (Optional) Should the account administrators be emailed when this alert is triggered?
* `email_addresses` - (Optional) A list of email addresses which alerts should be sent to.
* `retention_days` - (Optional) Specifies the number of days to keep in the Threat Detection audit logs.
* `storage_account_access_key` - (Optional) Specifies the identifier key of the Threat Detection audit storage account. Required if `state` is `Enabled`.
* `storage_endpoint` - (Optional) Specifies the blob storage endpoint (e.g. <https://example.blob.core.windows.net>). This blob storage will hold all Threat Detection audit logs. Required if `state` is `Enabled`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Microsoft SQL Server ID.
* `fully_qualified_domain_name` - The fully qualified domain name of the Azure SQL Server (e.g. myServerName.database.windows.net)

---

An `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this SQL Server.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this SQL Server.

-> You can access the Principal ID via `${azurerm_mssql_server.example.identity.0.principal_id}` and the Tenant ID via `${azurerm_mssql_server.example.identity.0.tenant_id}`

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Microsoft SQL Server.
* `update` - (Defaults to 60 minutes) Used when updating the Microsoft SQL Server.
* `read` - (Defaults to 5 minutes) Used when retrieving the Microsoft SQL Server.
* `delete` - (Defaults to 60 minutes) Used when deleting the Microsoft SQL Server.

## Import

SQL Servers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sql_server.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver
```
