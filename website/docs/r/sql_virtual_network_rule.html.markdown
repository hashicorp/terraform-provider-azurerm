---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_virtual_network_rule"
sidebar_current: "docs-azurerm-resource-database-sql-virtual-network-rule"
description: |-
  Create a SQL Virtual Network Rule.
---

# azurerm_sql_virtual_network_rule

Allows you to add, update, or remove an Azure SQL server to a subnet of a virtual network.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-sql-server-vnet-rule"
  location = "West US"
}

resource "azurerm_virtual_network" "vnet" {
  name                = "example-vnet"
  address_space       = ["10.7.29.0/29"]
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_subnet" "subnet" {
  name                 = "example-subnet"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  virtual_network_name = "${azurerm_virtual_network.vnet.name}"
  address_prefix       = "10.7.29.0/29"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_sql_server" "sqlserver" {
    name                         = "unqiueazuresqlserver"
    resource_group_name          = "${azurerm_resource_group.example.name}"
    location                     = "${azurerm_resource_group.example.location}"
    version                      = "12.0"
    administrator_login          = "4dm1n157r470r"
    administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_sql_virtual_network_rule" "sqlvnetrule" {
  name                = "sql-vnet-rule"
  resource_group_name = "${azurerm_resource_group.example.name}"
  server_name         = "${azurerm_sql_server.sqlserver.name}"
  subnet_id           = "${azurerm_subnet.subnet.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SQL virtual network rule. Changing this forces a new resource to be created. Cannot be empty and must only contain alphanumeric characters and hyphens. Cannot start with a number, and cannot start or end with a hyphen.

~> **NOTE:** `name` must be between 1-128 characters long and must satisfy all of the requirements below:
1. Contains only alphanumeric and hyphen characters
2. Cannot start with a number or hyphen
3. Cannot end with a hyphen

* `resource_group_name` - (Required) The name of the resource group where the SQL server resides. Changing this forces a new resource to be created.

* `server_name` - (Required) The name of the SQL Server to which this SQL virtual network rule will be applied to. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the subnet that the SQL server will be connected to.

* `ignore_missing_vnet_service_endpoint` - (Optional) Create the virtual network rule before the subnet has the virtual network service endpoint enabled. The default value is false.

~> **NOTE:** If `ignore_missing_vnet_service_endpoint` is false, and the target subnet does not contain the `Microsoft.SQL` endpoint in the `service_endpoints` array, the deployment will fail when it tries to create the SQL virtual network rule.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the SQL virtual network rule.

## Import

SQL Virtual Network Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sql_virtual_network_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver/virtualNetworkRules/vnetrulename
```
