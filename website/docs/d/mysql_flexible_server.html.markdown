---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_mysql_flexible_server"
description: |-
  Gets information about an existing MySQL Flexible Server.
---

# azurerm_mysql_flexible_server

Use this data source to access information about an existing MySQL Flexible Server.

## Example Usage

```hcl
data "azurerm_mysql_flexible_server" "example" {
  name                = "existingMySqlFlexibleServer"
  resource_group_name = "existingResGroup"
}

output "id" {
  value = data.azurerm_mysql_flexible_server.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MySQL Flexible Server.

* `resource_group_name` - (Required) The name of the resource group for the MySQL Flexible Server.

## Attributes Reference

* `id` - The ID of the MySQL Flexible Server.

* `fqdn` -  The fully qualified domain name of the MySQL Flexible Server.

* `public_network_access_enabled` - Is the public network access enabled?

* `replica_capacity` - The maximum number of replicas that a primary MySQL Flexible Server can have.

* `location` - The Azure Region of the MySQL Flexible Server.

* `administrator_login` - The Administrator login of the MySQL Flexible Server.

* `backup_retention_days` - The backup retention days of the MySQL Flexible Server.

* `delegated_subnet_id` - The ID of the virtual network subnet the MySQL Flexible Server is created in.

* `geo_redundant_backup_enabled` - Is geo redundant backup enabled?

* `high_availability` - A `high_availability` block for this MySQL Flexible Server as defined below.

* `maintenance_window` - A `maintenance_window` block for this MySQL Flexible Server as defined below.

* `private_dns_zone_id` - The ID of the Private DNS zone of the MySQL Flexible Server.

* `replication_role` - The replication role of the MySQL Flexible Server.

* `sku_name` - The SKU Name of the MySQL Flexible Server.

* `storage` - A `storage` block for this MySQL Flexible Server as defined below.

* `version` - The version of the MySQL Flexible Server.

* `zone` - The Availability Zones where this MySQL Flexible Server is located.

* `tags` - A mapping of tags which are assigned to the MySQL Flexible Server.

---

A `high_availability` block exports the following:

* `mode` - The high availability mode of the MySQL Flexible Server.

* `standby_availability_zone` - The availability zone of the standby Flexible Server.

---

A `maintenance_window` block exports the following:

* `day_of_week` - The day of week of the maintenance window.

* `start_hour` - The start hour of the maintenance window.

* `start_minute` - The start minute of the maintenance window.

---

A `storage` block exports the following:

* `auto_grow_enabled` - Is Storage Auto Grow enabled?

* `iops` - The storage IOPS of the MySQL Flexible Server.

* `size_gb` - The max storage allowed for the MySQL Flexible Server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the MySQL Flexible Server.
