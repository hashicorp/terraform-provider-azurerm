---
subcategory: "Data Migration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_migration_service"
sidebar_current: "docs-azurerm-datasource-data-migration-service"
description: |-
  Gets information about an existing Data Migration Service
---

# Data Source: azurerm_data_migration_service

Use this data source to access information about an existing Data Migration Service.


## Example Usage

```hcl
data "azurerm_data_migration_service" "example" {
  name                  = "example-dms"
  resource_group_name   = "example-rg"
}

output "azurerm_dms_id" {
  value = "${data.azurerm_data_migration_service.example.id}"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) Specify the name of the data migration service.

* `resource_group_name` - (Required) Specifies the Name of the Resource Group within which the data migration service exists

## Attributes Reference

The following attributes are exported:

* `id` - Resource ID of data migration service.

* `location` - Azure location where the resource exists.

* `virtual_subnet_id` - The ID of the virtual subnet resource to which the data migration service exists.

* `sku_name` - The sku name of data migration service.

* `type` - The resource type chain (e.g. virtualMachines/extensions)

* `kind` - The resource kind.

* `tags` - A mapping of tags to assigned to the resource.