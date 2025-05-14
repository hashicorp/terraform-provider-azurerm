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
  depends_on            = [azurerm_subnet.example]
}

resource "azurerm_postgresql_flexible_server" "example" {
  name                          = "example-psqlflexibleserver"
  resource_group_name           = azurerm_resource_group.example.name
  location                      = azurerm_resource_group.example.location
  version                       = "12"
  delegated_subnet_id           = azurerm_subnet.example.id
  private_dns_zone_id           = azurerm_private_dns_zone.example.id
  public_network_access_enabled = false
  administrator_login           = "psqladmin"
  administrator_password        = "H@Sh1CoR3!"
  zone                          = "1"

  storage_mb   = 32768
  storage_tier = "P4"

  sku_name   = "B_Standard_B1ms"
  depends_on = [azurerm_private_dns_zone_virtual_network_link.example]

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this PostgreSQL Flexible Server. Changing this forces a new PostgreSQL Flexible Server to be created.

~> **Note:** This must be unique across the entire Azure service, not just within the resource group.

* `resource_group_name` - (Required) The name of the Resource Group where the PostgreSQL Flexible Server should exist. Changing this forces a new PostgreSQL Flexible Server to be created.

* `location` - (Required) The Azure Region where the PostgreSQL Flexible Server should exist. Changing this forces a new PostgreSQL Flexible Server to be created.

* `administrator_login` - (Optional) The Administrator login for the PostgreSQL Flexible Server. Required when `create_mode` is `Default` and `authentication.password_auth_enabled` is `true`.

-> **Note:** Once `administrator_login` is specified, changing this forces a new PostgreSQL Flexible Server to be created.

-> **Note:** To create with `administrator_login` specified or update with it first specified , `authentication.password_auth_enabled` must be set to `true`.

* `administrator_password` - (Optional) The Password associated with the `administrator_login` for the PostgreSQL Flexible Server.

* `administrator_password_wo` - (Optional) The Password associated with the `administrator_login` for the PostgreSQL Flexible Server.

~> **Note:** Either `administrator_password` or `administrator_password_wo` is required when `create_mode` is `Default` and `authentication.password_auth_enabled` is `true`.

* `administrator_password_wo_version` - (Optional) An integer value used to trigger an update for `administrator_password_wo`. This property should be incremented when updating `administrator_password_wo`.

* `authentication` - (Optional) An `authentication` block as defined below.

* `backup_retention_days` - (Optional) The backup retention days for the PostgreSQL Flexible Server. Possible values are between `7` and `35` days.

* `customer_managed_key` - (Optional) A `customer_managed_key` block as defined below. Changing this forces a new resource to be created.

* `geo_redundant_backup_enabled` - (Optional) Is Geo-Redundant backup enabled on the PostgreSQL Flexible Server. Defaults to `false`. Changing this forces a new PostgreSQL Flexible Server to be created.

* `create_mode` - (Optional) The creation mode which can be used to restore or replicate existing servers. Possible values are `Default`, `GeoRestore`, `PointInTimeRestore`, `Replica` and `Update`.

* `delegated_subnet_id` - (Optional) The ID of the virtual network subnet to create the PostgreSQL Flexible Server. The provided subnet should not have any other resource deployed in it and this subnet will be delegated to the PostgreSQL Flexible Server, if not already delegated. Changing this forces a new PostgreSQL Flexible Server to be created.

* `private_dns_zone_id` - (Optional) The ID of the private DNS zone to create the PostgreSQL Flexible Server.

~> **Note:** There will be a breaking change from upstream service at 15th July 2021, the `private_dns_zone_id` will be required when setting a `delegated_subnet_id`. For existing flexible servers who don't want to be recreated, you need to provide the `private_dns_zone_id` to the service team to manually migrate to the specified private DNS zone. The `azurerm_private_dns_zone` should end with suffix `.postgres.database.azure.com`.

* `public_network_access_enabled` - (Optional) Specifies whether this PostgreSQL Flexible Server is publicly accessible. Defaults to `true`.

-> **Note:** `public_network_access_enabled` must be set to `false` when `delegated_subnet_id` and `private_dns_zone_id` have a value.

* `high_availability` - (Optional) A `high_availability` block as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `maintenance_window` - (Optional) A `maintenance_window` block as defined below.

* `point_in_time_restore_time_in_utc` - (Optional) The point in time to restore from `source_server_id` when `create_mode` is `GeoRestore`, `PointInTimeRestore`. Changing this forces a new PostgreSQL Flexible Server to be created.

* `replication_role` - (Optional) The replication role for the PostgreSQL Flexible Server. Possible value is `None`.

~> **Note:** The `replication_role` cannot be set while creating and only can be updated to `None` for replica server.

* `sku_name` - (Optional) The SKU Name for the PostgreSQL Flexible Server. The name of the SKU, follows the `tier` + `name` pattern (e.g. `B_Standard_B1ms`, `GP_Standard_D2s_v3`, `MO_Standard_E4s_v3`).

* `source_server_id` - (Optional) The resource ID of the source PostgreSQL Flexible Server to be restored. Required when `create_mode` is `GeoRestore`, `PointInTimeRestore` or `Replica`. Changing this forces a new PostgreSQL Flexible Server to be created.

* `auto_grow_enabled` - (Optional) Is the storage auto grow for PostgreSQL Flexible Server enabled? Defaults to `false`.

* `storage_mb` - (Optional) The max storage allowed for the PostgreSQL Flexible Server. Possible values are `32768`, `65536`, `131072`, `262144`, `524288`, `1048576`, `2097152`, `4193280`, `4194304`, `8388608`, `16777216` and `33553408`.

~> **Note:** If the `storage_mb` field is undefined on the initial deployment of the PostgreSQL Flexible Server resource it will default to `32768`. If the `storage_mb` field has been defined and then removed, the `storage_mb` field will retain the previously defined value.

~> **Note:** The `storage_mb` can only be scaled up, for example, you can scale the `storage_mb` from `32768` to `65536`, but not from `65536` to `32768`. Scaling down `storage_mb` forces a new PostgreSQL Flexible Server to be created.

* `storage_tier` - (Optional) The name of storage performance tier for IOPS of the PostgreSQL Flexible Server. Possible values are `P4`, `P6`, `P10`, `P15`,`P20`, `P30`,`P40`, `P50`,`P60`, `P70` or `P80`. Default value is dependant on the `storage_mb` value. Please see the `storage_tier` defaults based on `storage_mb` table below.

~> **Note:** The `storage_tier` can be scaled once every 12 hours, this restriction is in place to ensure stability and performance after any changes to your PostgreSQL Flexible Server's configuration.

* `tags` - (Optional) A mapping of tags which should be assigned to the PostgreSQL Flexible Server.

* `version` - (Optional) The version of PostgreSQL Flexible Server to use. Possible values are `11`,`12`, `13`, `14`, `15` and `16`. Required when `create_mode` is `Default`.

-> **Note:** Downgrading `version` isn't supported and will force a new PostgreSQL Flexible Server to be created.

-> **Note:** In-place version updates are irreversible and may cause downtime for the PostgreSQL Flexible Server, determined by the size of the instance.

* `zone` - (Optional) Specifies the Availability Zone in which the PostgreSQL Flexible Server should be located.

-> **Note:** Azure will automatically assign an Availability Zone if one is not specified. If the PostgreSQL Flexible Server fails-over to the Standby Availability Zone, the `zone` will be updated to reflect the current Primary Availability Zone. You can use [Terraform's `ignore_changes` functionality](https://www.terraform.io/docs/language/meta-arguments/lifecycle.html#ignore_changes) to ignore changes to the `zone` and `high_availability[0].standby_availability_zone` fields should you wish for Terraform to not migrate the PostgreSQL Flexible Server back to it's primary Availability Zone after a fail-over.

-> **Note:** The Availability Zones available depend on the Azure Region that the PostgreSQL Flexible Server is being deployed into - see [the Azure Availability Zones documentation](https://azure.microsoft.com/global-infrastructure/geographies/#geographies) for more information on which Availability Zones are available in each Azure Region.

---

An `authentication` block supports the following:

* `active_directory_auth_enabled` - (Optional) Whether Active Directory authentication is allowed to access the PostgreSQL Flexible Server. Defaults to `false`.

* `password_auth_enabled` - (Optional) Whether password authentication is allowed to access the PostgreSQL Flexible Server. Defaults to `true`.

* `tenant_id` - (Optional) The Tenant ID of the Azure Active Directory which is used by the Active Directory authentication. `active_directory_auth_enabled` must be set to `true`.

-> **Note:** Setting `active_directory_auth_enabled` to `true` requires a Service Principal for the Postgres Flexible Server. For more details see [this document](https://learn.microsoft.com/en-us/azure/postgresql/flexible-server/how-to-configure-sign-in-azure-ad-authentication).

-> **Note:** `tenant_id` is required when `active_directory_auth_enabled` is set to `true`. And it should not be specified when `active_directory_auth_enabled` is set to `false`

---

A `customer_managed_key` block supports the following:

* `key_vault_key_id` - (Required) The versioned ID of the Key Vault Key.

* `primary_user_assigned_identity_id` - (Optional) Specifies the primary user managed identity id for a Customer Managed Key. Must be added to `identity.identity_ids`.

* `geo_backup_key_vault_key_id` - (Optional) The versioned ID of the geo backup Key Vault Key.

~> **Note:** The key vault in which this key exists must be in the same region as the geo-redundant backup.

* `geo_backup_user_assigned_identity_id` - (Optional) The geo backup user managed identity id for a Customer Managed Key. Must be added to `identity.identity_ids`.

~> **Note:** This managed identity cannot be the same as `primary_user_assigned_identity_id`, additionally this identity must be created in the same region as the geo-redundant backup.

~> **Note:** `primary_user_assigned_identity_id` or `geo_backup_user_assigned_identity_id` is required when `type` is set to `UserAssigned`.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this PostgreSQL Flexible Server. Possible values are `UserAssigned` and `SystemAssigned`.

* `identity_ids` - (Optional) A list of User Assigned Managed Identity IDs to be assigned to this PostgreSQL Flexible Server. Required if used together with `customer_managed_key` block.

~> **Note:** `identity_ids` is required when `type` is set to `UserAssigned`.

---

A `maintenance_window` block supports the following:

* `day_of_week` - (Optional) The day of week for maintenance window, where the week starts on a Sunday, i.e. Sunday = `0`, Monday = `1`. Defaults to `0`.

* `start_hour` - (Optional) The start hour for maintenance window. Defaults to `0`.

* `start_minute` - (Optional) The start minute for maintenance window. Defaults to `0`.

-> **Note:** The specified `maintenance_window` is always defined in UTC time. When unspecified, the maintenance window falls back to the default [system-managed](https://learn.microsoft.com/en-us/azure/postgresql/flexible-server/how-to-maintenance-portal#specify-maintenance-schedule-options).

---

A `high_availability` block supports the following:

* `mode` - (Required) The high availability mode for the PostgreSQL Flexible Server. Possible value are `SameZone` or `ZoneRedundant`.

* `standby_availability_zone` - (Optional) Specifies the Availability Zone in which the standby Flexible Server should be located.

-> **Note:** Azure will automatically assign an Availability Zone if one is not specified. If the PostgreSQL Flexible Server fails-over to the Standby Availability Zone, the `zone` will be updated to reflect the current Primary Availability Zone. You can use [Terraform's `ignore_changes` functionality](https://www.terraform.io/docs/language/meta-arguments/lifecycle.html#ignore_changes) to ignore changes to the `zone` and `high_availability[0].standby_availability_zone` fields should you wish for Terraform to not migrate the PostgreSQL Flexible Server back to it's primary Availability Zone after a fail-over.

-> **Note:** The Availability Zones available depend on the Azure Region that the PostgreSQL Flexible Server is being deployed into - see [the Azure Availability Zones documentation](https://azure.microsoft.com/global-infrastructure/geographies/#geographies) for more information on which Availability Zones are available in each Azure Region.

---

## `storage_tier` defaults based on `storage_mb`

| `storage_mb` | GiB   | TiB | Default | Supported `storage_tier`'s           | Provisioned `IOPS`  |
|:------------:|:-----:|:---:|:-------:|:------------------------------------:|:-------------------:|
| 32768        | 32    |  -  | P4      | P4, P6, P10, P15, P20, P30, P40, P50 | 120                 |
| 65536        | 64    |  -  | P6      | P6, P10, P15, P20, P30, P40, P50     | 240                 |
| 131072       | 128   |  -  | P10     | P10, P15, P20, P30, P40, P50         | 500                 |
| 262144       | 256   |  -  | P15     | P15, P20, P30, P40, P50              | 1,100               |
| 524288       | 512   |  -  | P20     | P20, P30, P40, P50                   | 2,300               |
| 1048576      | 1024  |  1  | P30     | P30, P40, P50                        | 5,000               |
| 2097152      | 2048  |  2  | P40     | P40, P50                             | 7,500               |
| 4193280      | 4095  |  4  | P50     | P50                                  | 7,500               |
| 4194304      | 4096  |  4  | P50     | P50                                  | 7,500               |
| 8388608      | 8192  |  8  | P60     | P60, P70                             | 16,000              |
| 16777216     | 16384 |  16 | P70     | P70, P80                             | 18,000              |
| 33553408     | 32767 |  32 | P80     | P80                                  | 20,000              |

-> **Note:** Host Caching (ReadOnly and Read/Write) is supported on disk sizes less than 4194304 MiB. This means any disk that is provisioned up to 4193280 MiB can take advantage of Host Caching. Host caching is not supported for disk sizes larger than 4193280 MiB. For example, a P50 premium disk provisioned at 4193280 GiB can take advantage of Host caching while a P50 disk provisioned at 4194304 MiB cannot. Moving from a smaller disk size to a larger disk size, greater than 4193280 MiB, will cause the disk to lose the disk caching ability.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the PostgreSQL Flexible Server.

* `fqdn` - The FQDN of the PostgreSQL Flexible Server.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the PostgreSQL Flexible Server.
* `read` - (Defaults to 5 minutes) Used when retrieving the PostgreSQL Flexible Server.
* `update` - (Defaults to 1 hour) Used when updating the PostgreSQL Flexible Server.
* `delete` - (Defaults to 1 hour) Used when deleting the PostgreSQL Flexible Server.

## Import

PostgreSQL Flexible Servers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_postgresql_flexible_server.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforPostgreSQL/flexibleServers/server1
```
