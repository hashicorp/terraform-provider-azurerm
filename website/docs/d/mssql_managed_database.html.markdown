---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_managed_database"
description: |-
  Gets information about an existing managed database.
---

# Data Source: azurerm_mssql_managed_database

Use this data source to access information about an existing managed database.

## Example Usage

```hcl
data "azurerm_mssql_managed_database" "example" {
  name                  = "managed databas ename"
  resource_group_name   = "example-resources"
  managed_instance_name = "managed instance name"
}

output "manageddatabase" {
  value = data.azurerm_mssql_managed_database.example
}
```

## Argument Reference

* `name` - The name of the managed database.

* `managed_instance_name` - The name of the managed instance to which the database belongs to.

* `resource_group_name` - The name of the resource group which contains the managed database.


## Attributes Reference

* `managed_instance_id` - The managed instance id of the managed database.

* `location` - Specifies the supported Azure location where the resource exists.

* `collation` - The managed database SQL collation.
 
* `type` - The resource type of the managed database.

* `status` - The status of the managed database.

* `creation_date` - The creation date of the managed database.

* `earliest_restore_point` - The earliest restore point in time of the managed database.

* `default_secondary_location` - The default secondary location of the managed database.

* `tags` - The resource tags


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the managed database.
