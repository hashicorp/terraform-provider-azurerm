---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_virtual_network_rule"
description: |-
  Manages an Azure SQL Virtual Network Rule.
---

# azurerm_mssql_virtual_network_rule

Allows you to manage rules for allowing traffic between an Azure SQL server and a subnet of a virtual network.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-sql-server-vnet-rule"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.7.29.0/29"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.7.29.0/29"]
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_mssql_server" "example" {
  name                         = "uniqueazuresqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_virtual_network_rule" "example" {
  name      = "sql-vnet-rule"
  server_id = azurerm_mssql_server.example.id
  subnet_id = azurerm_subnet.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SQL virtual network rule. Changing this forces a new resource to be created.

* `server_id` - (Required) The resource ID of the SQL Server to which this SQL virtual network rule will be applied. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the subnet from which the SQL server will accept communications.

* `ignore_missing_vnet_service_endpoint` - (Optional) Create the virtual network rule before the subnet has the virtual network service endpoint enabled. Defaults to `false`.

~> **Note:** If `ignore_missing_vnet_service_endpoint` is false, and the target subnet does not contain the `Microsoft.SQL` endpoint in the `service_endpoints` array, the deployment will fail when it tries to create the SQL virtual network rule.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the SQL virtual network rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the SQL Virtual Network Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the SQL Virtual Network Rule.
* `update` - (Defaults to 30 minutes) Used when updating the SQL Virtual Network Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the SQL Virtual Network Rule.

## Import

SQL Virtual Network Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_virtual_network_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver/virtualNetworkRules/vnetrulename
```
