---
subcategory: "Database Migration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_database_migration_service"
description: |-
  Gets information about an existing Database Migration Service
---

# Data Source: azurerm_database_migration_service

Use this data source to access information about an existing Database Migration Service.


## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_database_migration_service" "example" {
  name                = "example-dms"
  resource_group_name = "example-rg"
}

output "azurerm_dms_id" {
  value = "${data.azurerm_database_migration_service.example.id}"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) Specify the name of the database migration service.

* `resource_group_name` - (Required) Specifies the Name of the Resource Group within which the database migration service exists

## Attributes Reference

The following attributes are exported:

* `id` - The ID of Database Migration Service.

* `location` - Azure location where the resource exists.

* `subnet_id` - The ID of the virtual subnet resource to which the database migration service exists.

* `sku_name` - The sku name of database migration service.

* `tags` - A mapping of tags to assigned to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the API Management API.
