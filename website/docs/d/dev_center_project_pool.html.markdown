---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_dev_center_project_pool"
description: |-
  Gets information about an existing Dev Center Project Pool.
---

# Data Source: azurerm_dev_center_project_pool

Use this data source to access information about an existing Dev Center Project Pool.

## Example Usage

```hcl
data "azurerm_dev_center_project_pool" "example" {
  name                  = azurerm_dev_center_project_pool.example.name
  dev_center_project_id = azurerm_dev_center_project_pool.example.dev_center_project_id
}

output "id" {
  value = data.azurerm_dev_center_project_pool.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Dev Center Project Pool.

* `dev_center_project_id` - (Required) The ID of the associated Dev Center Project.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center Project Pool.

* `dev_box_definition_name` - The name of the Dev Center Dev Box Definition.

* `local_administrator_enabled` - Specifies whether owners of Dev Boxes in the Dev Center Project Pool are added as local administrators on the Dev Box.

* `dev_center_attached_network_name` - The name of the Dev Center Attached Network in parent Project of the Dev Center Project Pool.

* `stop_on_disconnect_grace_period_minutes` - The specified time in minutes to wait before stopping a Dev Center Dev Box once disconnect is detected.

* `location` - The Azure Region where the Dev Center Project Pool exists.

* `tags` - A mapping of tags assigned to the Dev Center Project Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Center Project Pool.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.DevCenter`: 2025-02-01
