---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_managed_instance_ad_admin"
description: |-
  Manages a Managed instance AAD admin details.
---

# azurerm_mssql_managed_instance_ad_admin

Manages a Managed instance AAD admin details.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "westus"
}

resource "azurerm_network_security_group" "example" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "sql_mi-network"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  virtual_network_name = azurerm_virtual_network.example.name
  resource_group_name  = azurerm_resource_group.example.name
  address_prefixes     = ["10.0.1.0/24"]
  delegation {
    name = "miDelegation"

    service_delegation {
      name = "Microsoft.Sql/managedInstances"
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "example" {
  subnet_id                 = azurerm_subnet.example.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_route_table" "example" {
  name                = "example-routetable"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  route {
    name                   = "example"
    address_prefix         = "10.100.0.0/14"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "10.10.1.1"
  }
}

resource "azurerm_subnet_route_table_association" "example" {
  subnet_id      = azurerm_subnet.example.id
  route_table_id = azurerm_route_table.example.id
}

resource "azurerm_mssql_managed_instance" "example" {
  name                         = "sql-mi"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  administrator_login          = "demoReadUser"
  administrator_login_password = "ReadUser@123456"
  subnet_id                    = azurerm_subnet.example.id
  identity {
    type = "SystemAssigned"
  }
  sku {
    capacity = 8
    family   = "Gen5"
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
  }
  license_type                 = "LicenseIncluded"
  collation                    = "SQL_Latin1_General_CP1_CI_AS"
  proxy_override               = "Redirect"
  storage_size_gb              = 64
  vcores                       = 8
  public_data_endpoint_enabled = false
  timezone_id                  = "Central America Standard Time"
  minimal_tls_version          = "1.2"
}

resource "azurerm_mssql_managed_instance_ad_admin" "example" {
  managed_instance_name = azurerm_mssql_managed_instance.example.name
  resource_group_name   = azurerm_mssql_managed_instance.example.resource_group_name
  login_user_name       = "user@example.com"
  object_id             = "00000000-0000-0000-0000-000000000000"
}

```

## Argument Reference

The following arguments are supported:

* `managed_instance_name` - (Required) The resource id of the MS SQL Managed Instance. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The resource group of the managed instance. Changing this forces a new resource to be created.

* `login_user_name` - (Required) The Active Directory User name. Changing this forces a new resource to be created.

* `object_id` - (Required) The active directory user service principal. Changing this forces a new resource to be created.

* `tenant_id` - (Optional) The active directory tenant id

~> **NOTE:** The managed instance once created should be granted Directory Read permissions. This operation can only be executed by Global/Company administrator or Privileged Role Administrators in Azure AD. See [managed instance admin permissions](https://docs.microsoft.com/en-us/azure/azure-sql/database/authentication-aad-configure?tabs=azure-cli#provision-azure-ad-admin-sql-managed-instance) for more details.

~> **NOTE:** Currently only a single AD Admin user can be set. 

~> **NOTE:** It is the responsibility of the user to specify proper object id of the user. The SQL Managed instance endpoints do not check for correctness of the AD object ids specified.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Managed instance AD Admin.

* `type` - The resource type.

* `admin_type` - The admin user type. Defaults always to `ActiveDirectory`.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 15 minutes) Used when creating Managed instance Active directory admins. 
* `update` - (Defaults to 15 minutes) Used when updating Managed instance Active directory admins.
* `read` - (Defaults to 5 minutes) Used when retrieving Managed instance Active directory admins.
* `delete` - (Defaults to 5 minutes) Used when deleting Managed instance Active directory admins.

## Import

SQL managed instance AD Admin details can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_managed_instance_ad_admin.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Sql/managedInstances/sql-mi/administrators/ActiveDirectory
```
