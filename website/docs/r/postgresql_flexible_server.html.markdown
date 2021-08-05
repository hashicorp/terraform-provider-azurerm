---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_flexible_server"
description: |-
  Manages a PostgreSQL Flexible Server.
---

# azurerm_postgresql_flexible_server

Manages a PostgreSQL Flexible Server.

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
  name                = "example-vn"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "example-sn"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
  service_endpoints    = ["Microsoft.Storage"]
  delegation {
    name = "fs"
    service_delegation {
      name = "Microsoft.DBforPostgreSQL/flexibleServers"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}
resource "azurerm_private_dns_zone" "example" {
  name                = "example.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "example" {
  name                  = "exampleVnetZone.com"
  private_dns_zone_name = azurerm_private_dns_zone.example.name
  virtual_network_id    = azurerm_virtual_network.example.id
  resource_group_name   = azurerm_resource_group.example.name
}

resource "azurerm_postgresql_flexible_server" "example" {
  name                   = "example-psqlflexibleserver"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  version                = "12"
  delegated_subnet_id    = azurerm_subnet.example.id
  private_dns_zone_id    = azurerm_private_dns_zone.example.id
  administrator_login    = "psqladminun"
  administrator_password = "H@Sh1CoR3!"

  storage_mb = 32768

  sku_name   = "GP_Standard_D4s_v3"
  depends_on = [azurerm_private_dns_zone_virtual_network_link.example]

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this PostgreSQL Flexible Server. Changing this forces a new PostgreSQL Flexible Server to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the PostgreSQL Flexible Server should exist. Changing this forces a new PostgreSQL Flexible Server to be created.

* `location` - (Required) The Azure Region where the PostgreSQL Flexible Server should exist. Changing this forces a new PostgreSQL Flexible Server to be created.

* `administrator_login` - (Optional) The Administrator Login for the PostgreSQL Flexible Server. Required when `create_mode` is `Default`. Changing this forces a new PostgreSQL Flexible Server to be created.

* `administrator_password` - (Optional) The Password associated with the `administrator_login` for the PostgreSQL Flexible Server. Required when `create_mode` is `Default`.

* `zone` - (Optional) The availability zone of the PostgreSQL Flexible Server. Possible values are `1`, `2` and `3`. Changing this forces a new PostgreSQL Flexible Server to be created.

* `backup_retention_days` - (Optional) The backup retention days for the PostgreSQL Flexible Server. Possible values are between `7` and `35` days.

* `create_mode` - (Optional) The creation mode which can be used to restore or replicate existing servers. Possible values are `Default` and `PointInTimeRestore`. Changing this forces a new PostgreSQL Flexible Server to be created.

* `delegated_subnet_id` - (Optional) The ID of the virtual network subnet to create the PostgreSQL Flexible Server. The provided subnet should not have any other resource deployed in it and this subnet will be delegated to the PostgreSQL Flexible Server, if not already delegated. Changing this forces a new PostgreSQL Flexible Server to be created.

* `private_dns_zone_id` - (Optional) The ID of the private dns zone to create the PostgreSQL Flexible Server. Changing this forces a new PostgreSQL Flexible Server to be created.

~> **NOTE:** There will be a breaking change from upstream service at 15th July 2021, the `private_dns_zone_id` will be required when setting a `delegated_subnet_id`. For existing flexible servers who don't want to be recreated, you need to provide the `private_dns_zone_id` to the service team to manually migrate to the specified private dns zone. The `azurerm_private_dns_zone` should end with suffix `.postgres.database.azure.com`.

* `high_availability` - (Optional) A `high_availability` block as defined below.

* `maintenance_window` - (Optional) A `maintenance_window` block as defined below.

* `point_in_time_restore_time_in_utc` - (Optional) The point in time to restore from `creation_source_server_id` when `create_mode` is `PointInTimeRestore`. Changing this forces a new PostgreSQL Flexible Server to be created.

* `sku_name` - (Optional) The SKU Name for the PostgreSQL Flexible Server. The name of the SKU, follows the `tier` + `name` pattern (e.g. `B_Standard_B1ms`, `GP_Standard_D2s_v3`, `MO_Standard_E4s_v3`).

* `source_server_id` - (Optional) The resource ID of the source PostgreSQL Flexible Server to be restored. Required when `create_mode` is `PointInTimeRestore`. Changing this forces a new PostgreSQL Flexible Server to be created.

* `storage_mb` - (Optional) The max storage allowed for the PostgreSQL Flexible Server. Possible values are `32768`, `65536`, `131072`, `262144`, `524288`, `1048576`, `2097152`, `4194304`, `8388608`, `16777216`, and `33554432`.

* `version` - (Optional) The version of PostgreSQL Flexible Server to use. Possible values are `11`,`12` and `13`. Required when `create_mode` is `Default`. Changing this forces a new PostgreSQL Flexible Server to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the PostgreSQL Flexible Server.

---

A `maintenance_window` block supports the following:

* `day_of_week` - (Optional) The day of week for maintenance window. Defaults to `0`.

* `start_hour` - (Optional) The day of week for maintenance window. Defaults to `0`.

* `start_minute` - (Optional) The start minute for maintenance window. Defaults to `0`.

---

A `high_availability` block supports the following:

* `mode` - (Required) The high availability mode for the PostgreSQL Flexible Server. The only possible value is `ZoneRedundant`.

* `standby_availability_zone` - (Optional) The availability zone of the standby Flexible Server. Possible values are `1`, `2` and `3`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the PostgreSQL Flexible Server.

* `cmk_enabled` - The status showing whether the data encryption is enabled with a customer-managed key.

~> **Note:** Attribute `cmk_enabled` has been deprecated and will be removed in version 3.0 of the provider.

* `fqdn` - The FQDN of the PostgreSQL Flexible Server.

* `public_network_access_enabled` - Is public network access enabled?

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the PostgreSQL Flexible Server.
* `read` - (Defaults to 5 minutes) Used when retrieving the PostgreSQL Flexible Server.
* `update` - (Defaults to 1 hour) Used when updating the PostgreSQL Flexible Server.
* `delete` - (Defaults to 1 hour) Used when deleting the PostgreSQL Flexible Server.

## Import

PostgreSQL Flexible Servers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_postgresql_flexible_server.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforPostgreSQL/flexibleServers/server1
```
