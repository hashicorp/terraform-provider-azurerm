---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mysql_virtual_network_rule"
description: |-
  Manages a MySQL Virtual Network Rule.
---

# azurerm_mysql_virtual_network_rule

Manages a MySQL Virtual Network Rule.

-> **NOTE:** MySQL Virtual Network Rules [can only be used with SKU Tiers of `GeneralPurpose` or `MemoryOptimized`](https://docs.microsoft.com/en-us/azure/mysql/concepts-data-access-and-security-vnet)

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.7.29.0/29"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "internal" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.7.29.0/29"]
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_mysql_server" "example" {
  name                = "example-mysqlserver"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  administrator_login          = "mysqladminun"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name   = "B_Gen5_2"
  storage_mb = 5120
  version    = "5.7"

  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  ssl_enforcement_enabled      = true
}

resource "azurerm_mysql_virtual_network_rule" "example" {
  name                = "mysql-vnet-rule"
  resource_group_name = azurerm_resource_group.example.name
  server_name         = azurerm_mysql_server.example.name
  subnet_id           = azurerm_subnet.internal.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the MySQL Virtual Network Rule. Cannot be empty and must only contain alphanumeric characters and hyphens. Cannot start with a number, and cannot start or end with a hyphen. Changing this forces a new resource to be created.

~> **NOTE:** `name` must be between 1-128 characters long and must satisfy all of the requirements below:
1. Contains only alphanumeric and hyphen characters
2. Cannot start with a number or hyphen
3. Cannot end with a hyphen

* `resource_group_name` - (Required) The name of the resource group where the MySQL server resides. Changing this forces a new resource to be created.

* `server_name` - (Required) The name of the SQL Server to which this MySQL virtual network rule will be applied to. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the subnet that the MySQL server will be connected to.

~> **NOTE:** Due to [a bug in the Azure API](https://github.com/Azure/azure-rest-api-specs/issues/3719) this resource currently doesn't expose the `ignore_missing_vnet_service_endpoint` field and defaults this to `false`. Terraform will check during the provisioning of the Virtual Network Rule that the Subnet contains the Service Rule to verify that the Virtual Network Rule can be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MySQL Virtual Network Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the MySQL Virtual Network Rule.
* `update` - (Defaults to 30 minutes) Used when updating the MySQL Virtual Network Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the MySQL Virtual Network Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the MySQL Virtual Network Rule.

## Import

MySQL Virtual Network Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mysql_virtual_network_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.DBforMySQL/servers/myserver/virtualNetworkRules/vnetrulename
```
