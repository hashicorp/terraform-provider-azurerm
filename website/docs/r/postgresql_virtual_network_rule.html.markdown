---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_virtual_network_rule"
description: |-
  Manages a PostgreSQL Virtual Network Rule.
---

# azurerm_postgresql_virtual_network_rule

Manages a PostgreSQL Virtual Network Rule.

-> **NOTE:** PostgreSQL Virtual Network Rules [can only be used with SKU Tiers of `GeneralPurpose` or `MemoryOptimized`](https://docs.microsoft.com/en-us/azure/postgresql/concepts-data-access-and-security-vnet)

## Example Usage

```hcl
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

resource "azurerm_postgresql_server" "example" {
  name                = "postgresql-server-1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 5120
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "psqladminun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement              = "Enabled"
}

resource "azurerm_postgresql_virtual_network_rule" "example" {
  name                                 = "postgresql-vnet-rule"
  resource_group_name                  = azurerm_resource_group.example.name
  server_name                          = azurerm_postgresql_server.example.name
  subnet_id                            = azurerm_subnet.internal.id
  ignore_missing_vnet_service_endpoint = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the PostgreSQL virtual network rule. Cannot be empty and must only contain alphanumeric characters and hyphens. Cannot start with a number, and cannot start or end with a hyphen. Changing this forces a new resource to be created.

~> **NOTE:** `name` must be between 1-128 characters long and must satisfy all of the requirements below:
1. Contains only alphanumeric and hyphen characters
2. Cannot start with a number or hyphen
3. Cannot end with a hyphen

* `resource_group_name` - (Required) The name of the resource group where the PostgreSQL server resides. Changing this forces a new resource to be created.

* `server_name` - (Required) The name of the SQL Server to which this PostgreSQL virtual network rule will be applied to. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the subnet that the PostgreSQL server will be connected to.

* `ignore_missing_vnet_service_endpoint` - (Optional) Should the Virtual Network Rule be created before the Subnet has the Virtual Network Service Endpoint enabled? Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the PostgreSQL Virtual Network Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the PostgreSQL Virtual Network Rule.
* `update` - (Defaults to 30 minutes) Used when updating the PostgreSQL Virtual Network Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the PostgreSQL Virtual Network Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the PostgreSQL Virtual Network Rule.

## Import

PostgreSQL Virtual Network Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_postgresql_virtual_network_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.DBforPostgreSQL/servers/myserver/virtualNetworkRules/vnetrulename
```
