---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mysql_flexible_server"
description: |-
  Manages a MySQL Flexible Server.
---

# azurerm_mysql_flexible_server

Manages a MySQL Flexible Server.

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
      name = "Microsoft.DBforMySQL/flexibleServers"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_private_dns_zone" "example" {
  name                = "example.mysql.database.azure.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "example" {
  name                  = "exampleVnetZone.com"
  private_dns_zone_name = azurerm_private_dns_zone.example.name
  virtual_network_id    = azurerm_virtual_network.example.id
  resource_group_name   = azurerm_resource_group.example.name
}

resource "azurerm_mysql_flexible_server" "example" {
  name                   = "example-fs"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  administrator_login    = "psqladmin"
  administrator_password = "H@Sh1CoR3!"
  backup_retention_days  = 7
  delegated_subnet_id    = azurerm_subnet.example.id
  private_dns_zone_id    = azurerm_private_dns_zone.example.id
  sku_name               = "GP_Standard_D2ds_v4"

  depends_on = [azurerm_private_dns_zone_virtual_network_link.example]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this MySQL Flexible Server. Changing this forces a new MySQL Flexible Server to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the MySQL Flexible Server should exist. Changing this forces a new MySQL Flexible Server to be created.

* `location` - (Required) The Azure Region where the MySQL Flexible Server should exist. Changing this forces a new MySQL Flexible Server to be created.

* `administrator_login` - (Optional) The Administrator login for the MySQL Flexible Server. Required when `create_mode` is `Default`. Changing this forces a new MySQL Flexible Server to be created.

* `administrator_password` - (Optional) The Password associated with the `administrator_login` for the MySQL Flexible Server.

* `administrator_password_wo` - (Optional, Write-Only) The Password associated with the `administrator_login` for the MySQL Flexible Server.

* `administrator_password_wo_version` - (Optional) An integer value used to trigger an update for `administrator_password_wo`. This property should be incremented when updating `administrator_password_wo`.

~> **Note:** Either `administrator_password` or `administrator_password_wo` is required when `create_mode` is `Default`.

* `backup_retention_days` - (Optional) The backup retention days for the MySQL Flexible Server. Possible values are between `1` and `35` days. Defaults to `7`.

* `create_mode` - (Optional)The creation mode which can be used to restore or replicate existing servers. Possible values are `Default`, `PointInTimeRestore`, `GeoRestore`, and `Replica`. Changing this forces a new MySQL Flexible Server to be created.

~> **Note:** Creating a `GeoRestore` server requires the source server with `geo_redundant_backup_enabled` enabled.

~> **Note:** When a server is first created it may not be immediately available for `geo restore` or `replica`. It may take a few minutes to several hours for the necessary metadata to be populated. Please see the [Geo Restore](https://learn.microsoft.com/azure/mysql/single-server/how-to-restore-server-portal#geo-restore) and the [Replica](https://learn.microsoft.com/azure/mysql/flexible-server/concepts-read-replicas#create-a-replica) for more information.

* `customer_managed_key` - (Optional) A `customer_managed_key` block as defined below.

~> **Note:** `identity` is required when `customer_managed_key` is specified.

* `delegated_subnet_id` - (Optional) The ID of the virtual network subnet to create the MySQL Flexible Server. Changing this forces a new MySQL Flexible Server to be created.

* `geo_redundant_backup_enabled` - (Optional) Should geo redundant backup enabled? Defaults to `false`. Changing this forces a new MySQL Flexible Server to be created.

* `high_availability` - (Optional) A `high_availability` block as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `maintenance_window` - (Optional) A `maintenance_window` block as defined below.

* `point_in_time_restore_time_in_utc` - (Optional) The point in time to restore from `creation_source_server_id` when `create_mode` is `PointInTimeRestore`. Changing this forces a new MySQL Flexible Server to be created.

* `private_dns_zone_id` - (Optional) The ID of the private DNS zone to create the MySQL Flexible Server. Changing this forces a new MySQL Flexible Server to be created.

~> **Note:** The `private_dns_zone_id` is required when setting a `delegated_subnet_id`. The `azurerm_private_dns_zone` should end with suffix `.mysql.database.azure.com`.

* `public_network_access` - (Optional) Whether approved public traffic is allowed through the firewall to this server. Possible values are `Enabled` and `Disabled`.

~> **Note:** `public_network_access` is automatically set to `Disabled` if the server is created with VNet Integration (i.e. values are provided for `delegated_subnet_id` and `private_dns_zone_id`").

* `replication_role` - (Optional) The replication role. Possible value is `None`.

~> **Note:** The `replication_role` cannot be set while creating and only can be updated from `Replica` to `None`.

* `sku_name` - (Optional) The SKU Name for the MySQL Flexible Server.

-> **Note:** `sku_name` should start with SKU tier `B (Burstable)`, `GP (General Purpose)`, `MO (Memory Optimized)` like `B_Standard_B1ms`.

* `source_server_id` - (Optional) The resource ID of the source MySQL Flexible Server to be restored. Required when `create_mode` is `PointInTimeRestore`, `GeoRestore`, and `Replica`. Changing this forces a new MySQL Flexible Server to be created.

-> **Note:** The replica server is always created in the same resource group and subscription as the source server.

* `storage` - (Optional) A `storage` block as defined below.

* `version` - (Optional) The version of the MySQL Flexible Server to use. Possible values are `5.7`, and `8.0.21`. Changing this forces a new MySQL Flexible Server to be created.

* `zone` - (Optional) Specifies the Availability Zone in which this MySQL Flexible Server should be located. Possible values are `1`, `2` and `3`.

-> **Note:** Azure will automatically assign an Availability Zone if one is not specified. If the MySQL Flexible Server fails-over to the Standby Availability Zone, the `zone` will be updated to reflect the current Primary Availability Zone. You can use [Terraform's `ignore_changes` functionality](https://www.terraform.io/docs/language/meta-arguments/lifecycle.html#ignore_changes) to ignore changes to the `zone` and `high_availability[0].standby_availability_zone` fields should you wish for Terraform to not migrate the MySQL Flexible Server back to it's primary Availability Zone after a fail-over.

-> **Note:** The Availability Zones available depend on the Azure Region that the MySQL Flexible Server is being deployed into - see [the Azure Availability Zones documentation](https://azure.microsoft.com/global-infrastructure/geographies/#geographies) for more information on which Availability Zones are available in each Azure Region.

* `tags` - (Optional) A mapping of tags which should be assigned to the MySQL Flexible Server.

---

A `customer_managed_key` block supports the following:

* `key_vault_key_id` - (Optional) The ID of the Key Vault Key.

* `primary_user_assigned_identity_id` - (Optional) Specifies the primary user managed identity id for a Customer Managed Key. Should be added with `identity_ids`.

* `geo_backup_key_vault_key_id` - (Optional) The ID of the geo backup Key Vault Key. It can't cross region and need Customer Managed Key in same region as geo backup.

* `geo_backup_user_assigned_identity_id` - (Optional) The geo backup user managed identity id for a Customer Managed Key. Should be added with `identity_ids`. It can't cross region and need identity in same region as geo backup.

~> **Note:** `primary_user_assigned_identity_id` or `geo_backup_user_assigned_identity_id` is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this MySQL Flexible Server. The only possible value is `UserAssigned`.

* `identity_ids` - (Required) A list of User Assigned Managed Identity IDs to be assigned to this MySQL Flexible Server.

---

A `high_availability` block supports the following:

* `mode` - (Required) The high availability mode for the MySQL Flexible Server. Possibles values are `SameZone` and `ZoneRedundant`.

~> **Note:** `storage[0].auto_grow_enabled` must be enabled when `high_availability` is enabled. To change the `high_availability` for a MySQL Flexible Server created with `high_availability` disabled during creation, the resource has to be recreated.

* `standby_availability_zone` - (Optional) Specifies the Availability Zone in which the standby Flexible Server should be located. Possible values are `1`, `2` and `3`.

-> **Note:** Azure will automatically assign an Availability Zone if one is not specified. If the MySQL Flexible Server fails-over to the Standby Availability Zone, the `zone` will be updated to reflect the current Primary Availability Zone. You can use [Terraform's `ignore_changes` functionality](https://www.terraform.io/docs/language/meta-arguments/lifecycle.html#ignore_changes) to ignore changes to the `zone` and `high_availability[0].standby_availability_zone` fields should you wish for Terraform to not migrate the MySQL Flexible Server back to it's primary Availability Zone after a fail-over.

-> **Note:** The Availability Zones available depend on the Azure Region that the MySQL Flexible Server is being deployed into - see [the Azure Availability Zones documentation](https://azure.microsoft.com/global-infrastructure/geographies/#geographies) for more information on which Availability Zones are available in each Azure Region.

---

A `maintenance_window` block supports the following:

* `day_of_week` - (Optional) The day of week for maintenance window. Defaults to `0`.

* `start_hour` - (Optional) The start hour for maintenance window. Defaults to `0`.

* `start_minute` - (Optional) The start minute for maintenance window. Defaults to `0`.

---

A `storage` block supports the following:

* `auto_grow_enabled` - (Optional) Should Storage Auto Grow be enabled? Defaults to `true`.

* `io_scaling_enabled` - (Optional) Should IOPS be scaled automatically? If `true`, `iops` can not be set. Defaults to `false`.

* `iops` - (Optional) The storage IOPS for the MySQL Flexible Server. Possible values are between `360` and `20000`.

* `log_on_disk_enabled` - (Optional) Should Storage Log On Disk be enabled? Defaults to `false`.

* `size_gb` - (Optional) The max storage allowed for the MySQL Flexible Server. Possible values are between `20` and `16384`.

-> **Note:** Decreasing `size_gb` forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the MySQL Flexible Server.

* `fqdn` - The fully qualified domain name of the MySQL Flexible Server.

* `replica_capacity` - The maximum number of replicas that a primary MySQL Flexible Server can have.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 hours) Used when creating the MySQL Flexible Server.
* `read` - (Defaults to 5 minutes) Used when retrieving the MySQL Flexible Server.
* `update` - (Defaults to 2 hours) Used when updating the MySQL Flexible Server.
* `delete` - (Defaults to 1 hour) Used when deleting the MySQL Flexible Server.

## Import

MySQL Flexible Servers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mysql_flexible_server.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DBforMySQL/flexibleServers/flexibleServer1
```
