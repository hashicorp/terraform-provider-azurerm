---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mysql_active_directory_administrator"
description: |-
  Manages an Active Directory administrator on a MySQL server
---

# azurerm_mysql_active_directory_administrator

Allows you to set a user or group as the AD administrator for an MySQL server in Azure

~> **Note:** Azure Database for MySQL Single Server and its sub resources are scheduled for retirement by 2024-09-16 and will migrate to using Azure Database for MySQL Flexible Server: https://go.microsoft.com/fwlink/?linkid=2216041. The `azurerm_mysql_active_directory_administrator` resource is deprecated and will be removed in v4.0 of the AzureRM Provider. Please use the `azurerm_mysql_flexible_server_active_directory_administrator` resource instead.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mysql_server" "example" {
  name                = "example-mysqlserver"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  administrator_login          = "mysqladminun"
  administrator_login_password = "H@Sh1CoR3!"
  ssl_enforcement_enabled      = true

  sku_name   = "B_Gen5_2"
  storage_mb = 5120
  version    = "5.7"
}

resource "azurerm_mysql_active_directory_administrator" "example" {
  server_name         = azurerm_mysql_server.example.name
  resource_group_name = azurerm_resource_group.example.name
  login               = "sqladmin"
  tenant_id           = data.azurerm_client_config.current.tenant_id
  object_id           = data.azurerm_client_config.current.object_id
}
```

## Argument Reference

The following arguments are supported:

* `server_name` - (Required) The name of the MySQL Server on which to set the administrator. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group for the MySQL server. Changing this forces a new resource to be created.

* `login` - (Required) The login name of the principal to set as the server administrator

* `object_id` - (Required) The ID of the principal to set as the server administrator. For a managed identity this should be the Client ID of the identity.

* `tenant_id` - (Required) The Azure Tenant ID

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the MySQL Active Directory Administrator.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the MySQL Active Directory Administrator.
* `update` - (Defaults to 30 minutes) Used when updating the MySQL Active Directory Administrator.
* `read` - (Defaults to 5 minutes) Used when retrieving the MySQL Active Directory Administrator.
* `delete` - (Defaults to 30 minutes) Used when deleting the MySQL Active Directory Administrator.

## Import

A MySQL Active Directory Administrator can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mysql_active_directory_administrator.administrator /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.DBforMySQL/servers/myserver/administrators/activeDirectory
```
