---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_managed_instance_active_directory_administrator"
description: |-
  Manages an Active Directory Administrator on a Microsoft Azure SQL Managed Instance
---

# azurerm_mssql_managed_instance_active_directory_administrator

Allows you to set a user, group or service principal as the AAD Administrator for an Azure SQL Managed Instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "rg-example"
  location = "West Europe"
}

data "azurerm_client_config" "current" {
}

resource "azurerm_virtual_network" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "example"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_mssql_managed_instance" "example" {
  name                = "managedsqlinstance"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  license_type       = "BasePrice"
  sku_name           = "GP_Gen5"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.example.id
  vcores             = 4

  administrator_login          = "msadministrator"
  administrator_login_password = "thisIsDog11"

  identity {
    type = "SystemAssigned"
  }
}

resource "azuread_directory_role" "reader" {
  display_name = "Directory Readers"
}

resource "azuread_directory_role_member" "example" {
  role_object_id   = azuread_directory_role.reader.object_id
  member_object_id = azurerm_mssql_managed_instance.example.identity[0].principal_id
}

resource "azuread_user" "admin" {
  user_principal_name = "ms.admin@hashicorp.com"
  display_name        = "Ms Admin"
  mail_nickname       = "ms.admin"
  password            = "SecretP@sswd99!"
}

resource "azurerm_mssql_managed_instance_active_directory_administrator" "example" {
  managed_instance_id = azurerm_mssql_managed_instance.example.id
  login_username      = "msadmin"
  object_id           = azuread_user.admin.object_id
  tenant_id           = data.azurerm_client_config.current.tenant_id
}
```

## Argument Reference

The following arguments are supported:

* `managed_instance_id` - (Required) The ID of the Azure SQL Managed Instance for which to set the administrator. Changing this forces a new resource to be created.

* `login_username` - (Required) The login name of the principal to set as the Managed Instance Administrator.

* `object_id` - (Required) The Object ID of the principal to set as the Managed Instance Administrator.

* `tenant_id` - (Required) The Azure Active Directory Tenant ID.

* `azuread_authentication_only` - (Optional) When `true`, only permit logins from AAD users and administrators. When `false`, also allow local database users.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the SQL Managed Instance Active Directory Administrator.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the SQL Active Directory Administrator.
* `read` - (Defaults to 5 minutes) Used when retrieving the SQL Active Directory Administrator.
* `update` - (Defaults to 30 minutes) Used when updating the SQL Active Directory Administrator.
* `delete` - (Defaults to 3 hours) Used when deleting the SQL Active Directory Administrator.

## Import

An Azure SQL Active Directory Administrator can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_managed_instance_active_directory_administrator.administrator /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/managedInstances/mymanagedinstance/administrators/activeDirectory
```
