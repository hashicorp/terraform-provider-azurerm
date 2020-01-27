---
subcategory: "Data Migration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_migration_service"
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

* `id` - The ID of Data Migration Service.

* `location` - Azure location where the resource exists.

* `subnet_id` - The ID of the virtual subnet resource to which the data migration service exists.

* `sku_name` - The sku name of data migration service.

* `tags` - A mapping of tags to assigned to the resource.
