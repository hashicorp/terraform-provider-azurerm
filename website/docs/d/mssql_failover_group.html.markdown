---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_mssql_failover_group"
description: |-
  Gets information about an existing Microsoft Azure SQL Failover Group.

---

# Data Source: azurerm_mssql_failover_group

Use this data source to access information about an existing Microsoft Azure SQL Failover Group.

## Example Usage

```hcl

data "azurerm_mssql_failover_group" "example" {
  name      = "example"
  server_id = "example-sql-server"
}

output "mssql_failover_group_id" {
  value = data.azurerm_mssql_failover_group.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Failover Group.

* `server_id` - (Required) The ID of the primary SQL Server where the Failover Group exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Failover Group.

* `databases` - A set of database names in the Failover Group.

* `partner_server` - A `partner_server` block as defined below.

* `readonly_endpoint_failover_policy_enabled` - Whether failover is enabled for the readonly endpoint.

* `read_write_endpoint_failover_policy` - A `read_write_endpoint_failover_policy` block as defined below.
 
* `tags` - A mapping of tags which are assigned to the resource.

---

A `partner_server` block exports the following:

* `id` - The ID of the partner SQL server.

* `location` - The location of the partner server.

* `role` - The replication role of the partner server.

---

The `read_write_endpoint_failover_policy` block exports the following:

* `mode` - The failover policy of the read-write endpoint for the Failover Group.

* `grace_minutes` - The grace period in minutes, before failover with data loss is attempted for the read-write endpoint.


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Failover Group.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Sql` - 2023-08-01-preview
