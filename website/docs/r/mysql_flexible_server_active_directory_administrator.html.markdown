---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mysql_flexible_server_active_directory_administrator"
description: |-
  Manages an Active Directory administrator on a MySQL Flexible Server
---

# azurerm_mysql_flexible_server_active_directory_administrator

Manages an Active Directory administrator on a MySQL Flexible Server

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mysql_flexible_server" "example" {
  name                   = "example-mysqlfs"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  administrator_login    = "_admin_Terraform_892123456789312"
  administrator_password = "QAZwsx123"
  sku_name               = "B_Standard_B1ms"
  zone                   = "2"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "exampleUAI"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_mysql_flexible_server_active_directory_administrator" "example" {
  server_id   = azurerm_mysql_flexible_server.example.id
  identity_id = azurerm_user_assigned_identity.example.id
  login       = "sqladmin"
  object_id   = data.azurerm_client_config.current.client_id
  tenant_id   = data.azurerm_client_config.current.tenant_id
}
```

## Arguments Reference

The following arguments are supported:

* `server_id` - (Required) The resource ID of the MySQL Flexible Server. Changing this forces a new resource to be created.

* `identity_id` - (Required) The resource ID of the identity used for AAD Authentication.

* `login` - (Required) The login name of the principal to set as the server administrator

* `object_id` - (Required) The ID of the principal to set as the server administrator. For a managed identity this should be the Client ID of the identity.

* `tenant_id` - (Required) The Azure Tenant ID.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the MySQL Flexible Server Active Directory Administrator.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the MySQL Flexible Server Active Directory Administrator.
* `read` - (Defaults to 5 minutes) Used when retrieving the MySQL Flexible Server Active Directory Administrator.
* `update` - (Defaults to 30 minutes) Used when updating the MySQL Flexible Server Active Directory Administrator.
* `delete` - (Defaults to 30 minutes) Used when deleting the MySQL Flexible Server Active Directory Administrator.

## Import

A MySQL Flexible Server Active Directory Administrator can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mysql_flexible_server_active_directory_administrator.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DBforMySQL/flexibleServers/server1/administrators/ActiveDirectory
```
