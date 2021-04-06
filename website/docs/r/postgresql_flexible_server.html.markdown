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

resource "azurerm_postgresql_flexible_server" "example" {
  name                   = "example-psqlflexibleserver"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  version                = "12"
  delegated_subnet_id    = azurerm_subnet.example.id
  administrator_login    = "psqladminun"
  administrator_password = "H@Sh1CoR3!"

  storage_mb = 32768

  sku {
    name = "Standard_D4s_v3"
    tier = "GeneralPurpose"
  }
}
```

## Arguments Reference

The following arguments are supported:
* `name` - (Required) The name which should be used for this PostgreSQL Flexible Server. Changing this forces a new PostgreSQL Flexible Server to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the PostgreSQL Flexible Server should exist. Changing this forces a new PostgreSQL Flexible Server to be created.

* `location` - (Required) The Azure Region where the PostgreSQL Flexible Server should exist. Changing this forces a new PostgreSQL Flexible Server to be created.

* `administrator_login` - (Optional) The Administrator Login for the PostgreSQL Flexible Server. Required when `create_mode` is `Default`. Changing this forces a new PostgreSQL Flexible Server to be created.

* `administrator_password` - (Optional) The Password associated with the `administrator_login` for the PostgreSQL Flexible Server. Required when `create_mode` is `Default`.

* `availability_zone` - (Optional) The availability Zone of the PostgreSQL Flexible Server. Possible values are  `none`, `1`, `2` and `3`. Changing this forces a new PostgreSQL Flexible Server to be created.

* `backup_retention_days` - (Optional) The backup retention days for the PostgreSQL Flexible Server. Possible values are between `7` and `35` days.

* `create_mode` - (Optional) The creation mode which can be used to restore or replicate existing servers. Possible values are `Default` and `PointInTimeRestore`. Changing this forces a new PostgreSQL Flexible Server to be created.

* `delegated_subnet_id` - (Optional) The ID of the virtual network subnet to create the PostgreSQL Flexible Server. The provided subnet should not have any other resource deployed in it and this subnet will be delegated to the PostgreSQL Flexible Server, if not already delegated. Changing this forces a new PostgreSQL Flexible Server to be created.

* `ha_enabled` - (Optional) Should High availability for the PostgreSQL Flexible Server be enabled? If enabled the server will provisions a physically separate primary and standby PostgreSQL Flexible Server in different zones. Defaults to `false`.

* `identity` - (Optional) A `identity` block as defined below. Changing this forces a new PostgreSQL Flexible Server to be created.

* `maintenance_window` - (Optional) A `maintenance_window` block as defined below.

* `point_in_time_utc` - (Optional) The point in time to restore from `creation_source_server_id` when `create_mode` is `PointInTimeRestore`. Changing this forces a new PostgreSQL Flexible Server to be created.

* `sku` - (Optional) A `sku` block as defined below.

* `source_server_id` - (Optional) The resource ID of the source PostgreSQL Flexible Server to be restored. Required when `create_mode` is `PointInTimeRestore`. Changing this forces a new PostgreSQL Flexible Server to be created.

* `storage_mb` - (Optional) The max storage allowed for the PostgreSQL Flexible Server. Possible values are `32768`, `65536`, `131072`, `262144`, `524288`, `1048576`, `2097152`, `4194304`, `8388608`, `16777216`, and `33554432`.

* `version` - (Optional) The version of PostgreSQL Flexible Server to use. Possible values are `11` and `12`. Required when `create_mode` is `Default`. Changing this forces a new PostgreSQL Flexible Server to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the PostgreSQL Flexible Server.

---

A `identity` block supports the following:

* `type` - (Required) The Type of Identity which should be used for this PostgreSQL Flexible Server. Possible value is only `SystemAssigned`. Changing this forces a new PostgreSQL Flexible Server to be created.

---

A `maintenance_window` block supports the following:

* `day_of_week` - (Optional) The day of week for maintenance window. Defaults to `0`.

* `start_hour` - (Optional) The day of week for maintenance window. Defaults to `0`.

* `start_minute` - (Optional) The start minute for maintenance window. Defaults to `0`.

---

A `sku` block supports the following:

* `name` - (Required) The SKU Name for the PostgreSQL Flexible Server. Possible values are `Standard_B1ms`, `Standard_B2s`, `Standard_D2s_v3`, `Standard_D4s_v3`, `Standard_D8s_v3`, `Standard_D16s_v3`, `Standard_D32s_v3`, `Standard_D48s_v3`, `Standard_D64s_v3`, `Standard_E2s_v3`, `Standard_E4s_v3`, `Standard_E8s_v3`, `Standard_E16s_v3`, `Standard_E32s_v3`, `Standard_E48s_v3`, `Standard_E64s_v3`.

* `tier` - (Required) The SKU tier for the PostgreSQL Flexible Server. Possible values are `Burstable`, `GeneralPurpose`, or `MemoryOptimized`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the PostgreSQL Flexible Server.

* `cmk_enabled` - The status showing whether the data encryption is enabled with a customer-managed key.

* `fqdn` - The FQDN of the PostgreSQL Flexible Server.

* `ha_state` - The state of the High Availability server.

* `public_network_access` - Is public network access enabled?

* `standby_availability_zone` -  The standby availability Zone information of the server.

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
