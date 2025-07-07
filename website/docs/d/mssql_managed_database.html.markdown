---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_mssql_managed_database"
description: |-
  Gets information about an existing Azure SQL Azure Managed Database.
---

# Data Source: azurerm_mssql_managed_database

Use this data source to access information about an existing Azure SQL Azure Managed Database.

## Example Usage

```hcl
data "azurerm_mssql_managed_database" "example" {
  name                  = "example"
  resource_group_name   = azurerm_resource_group.example.name
  managed_instance_name = azurerm_mssql_managed_instance.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Azure SQL Azure Managed Database.

* `managed_instance_id` - (Required) The SQL Managed Instance ID.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The Azure SQL Managed Database ID.

* `long_term_retention_policy` - A `long_term_retention_policy` block as defined below.

* `resource_group_name` - The name of the Resource Group where the Azure SQL Azure Managed Instance exists.

* `managed_instance_name` - The name of the Managed Instance.

* `point_in_time_restore` - A `point_in_time_restore` block as defined below.

* `short_term_retention_days` -  The backup retention period in days. This is how many days Point-in-Time Restore will be supported.

---

A `long_term_retention_policy` block exports the following:

* `immutable_backups_enabled` - Specifies if the backups are immutable.

* `monthly_retention` - The monthly retention policy for an LTR backup in an ISO 8601 format.

* `week_of_year` - The week of year to take the yearly backup.

* `weekly_retention` - The weekly retention policy for an LTR backup in an ISO 8601 format.

* `yearly_retention` - The yearly retention policy for an LTR backup in an ISO 8601 format.

---

A `point_in_time_restore` block exports the following:

* `restore_point_in_time` - The point in time for the restore from `source_database_id`.

* `source_database_id` - The source database ID that is used to restore from.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Azure SQL Azure Managed Database.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Sql`: 2023-08-01-preview
