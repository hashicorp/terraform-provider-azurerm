---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_elasticpool"
sidebar_current: "docs-azurerm-datasource-mssql-elasticpool"
description: |-
  Gets information about an existing SQL elastic pool.
---

# Data Source: azurerm_mssql_elasticpool

Use this data source to access information about an existing SQL elastic pool.

## Example Usage

```hcl
data "azurerm_mssql_elasticpool" "example" {
  name                = "mssqlelasticpoolname"
  resource_group_name = "example-resources"
  server_name         = "example-sql-server"
}

output "elasticpool_id" {
  value = "${data.azurerm_mssql_elasticpool.example.id}"
}
```

## Argument Reference

* `name` - (Required) The name of the elastic pool.

* `resource_group_name` - (Required) The name of the resource group which contains the elastic pool.

* `server_name` - (Required) The name of the SQL Server which contains the elastic pool.

## Attributes Reference

* `location` - Specifies the supported Azure location where the resource exists.

* `max_size_gb` - The max data size of the elastic pool in gigabytes.
 
* `max_size_bytes` - The max data size of the elastic pool in bytes.

* `per_db_min_capacity` - The minimum capacity all databases are guaranteed.

* `per_db_max_capacity` - The maximum capacity any one database can consume.

* `tags` - A mapping of tags to assign to the resource.

* `zone_redundant` - Whether or not this elastic pool is zone redundant.
