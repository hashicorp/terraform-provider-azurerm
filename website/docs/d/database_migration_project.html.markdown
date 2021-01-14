---
subcategory: "Database Migration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_database_migration_project"
description: |-
  Gets information about an existing Database Migration Project
---

# Data Source: azurerm_database_migration_project

Use this data source to access information about an existing Database Migration Project.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_database_migration_project" "example" {
  name                = "example-dbms-project"
  resource_group_name = "example-rg"
  service_name        = "example-dbms"
}

output "name" {
  value = "${data.azurerm_database_migration_project.example.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the database migration project.

* `resource_group_name` - (Required) Name of the resource group where resource belongs to.

* `service_name` - (Required) Name of the database migration service where resource belongs to.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of Database Migration Project.

* `location` - Azure location where the resource exists.

* `source_platform` - The platform type of the migration source.

* `target_platform` - The platform type of the migration target.

* `tags` - A mapping of tags to assigned to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the API Management API.
