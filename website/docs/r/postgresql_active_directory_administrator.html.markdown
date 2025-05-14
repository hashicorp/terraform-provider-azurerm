---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_active_directory_administrator"
description: |-
  Manages an Active Directory administrator on a PostgreSQL server
---

# azurerm_postgresql_active_directory_administrator

Allows you to set a user or group as the AD administrator for an PostgreSQL server in Azure

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_postgresql_server" "example" {
  name                         = "example-psqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "9.6"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
  sku_name                     = "GP_Gen5_2"
  ssl_enforcement_enabled      = true
}

resource "azurerm_postgresql_active_directory_administrator" "example" {
  server_name         = azurerm_postgresql_server.example.name
  resource_group_name = azurerm_resource_group.example.name
  login               = "sqladmin"
  tenant_id           = data.azurerm_client_config.current.tenant_id
  object_id           = data.azurerm_client_config.current.object_id
}
```

## Argument Reference

The following arguments are supported:

* `server_name` - (Required) The name of the PostgreSQL Server on which to set the administrator. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group for the PostgreSQL server. Changing this forces a new resource to be created.

* `login` - (Required) The login name of the principal to set as the server administrator

* `object_id` - (Required) The ID of the principal to set as the server administrator. For a managed identity this should be the Client ID of the identity.

* `tenant_id` - (Required) The Azure Tenant ID

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the PostgreSQL Active Directory Administrator.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the PostgreSQL Active Directory Administrator.
* `read` - (Defaults to 5 minutes) Used when retrieving the PostgreSQL Active Directory Administrator.
* `update` - (Defaults to 30 minutes) Used when updating the PostgreSQL Active Directory Administrator.
* `delete` - (Defaults to 30 minutes) Used when deleting the PostgreSQL Active Directory Administrator.

## Import

A PostgreSQL Active Directory Administrator can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_postgresql_active_directory_administrator.administrator /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.DBforPostgreSQL/servers/myserver
```
