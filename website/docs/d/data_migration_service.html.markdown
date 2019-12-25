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

* `resource_group_name` - (Required) Name of the resource group

* `name` - (Required) Name of the service


## Attributes Reference

The following attributes are exported:

* `id` - Resource ID.

* `location` - Resource location.

* `virtual_subnet_id` - The ID of the virtual subnet resource to which the service should be joined

* `sku_name` - The resource's sku name

* `type` - The resource type chain (e.g. virtualMachines/extensions)

* `kind` - The resource kind.

* `provisioning_state` - The resource's provisioning state

* `tags` - Resource tags.