---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_dev_center_dev_box_definition"
description: |-
  Gets information about an existing Dev Center Dev Box Definition.
---

# Data Source: azurerm_dev_center_dev_box_definition

Use this data source to access information about an existing Dev Center Dev Box Definition.

## Example Usage

```hcl
data "azurerm_dev_center_dev_box_definition" "example" {
  name          = azurerm_dev_center_dev_box_definition.example.name
  dev_center_id = azurerm_dev_center_dev_box_definition.example.dev_center_id
}

output "id" {
  value = data.azurerm_dev_center_dev_box_definition.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Dev Center Dev Box Definition.

* `dev_center_id` - (Required) The ID of the associated Dev Center.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center Dev Box Definition.

* `image_reference_id` - The ID of the image for the Dev Center Dev Box Definition.

* `location` - The Azure Region where the Dev Center Dev Box Definition exists.

* `sku_name` - The name of the SKU for the Dev Center Dev Box Definition.

* `tags` - A mapping of tags assigned to the Dev Center Dev Box Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Center Dev Box Definition.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.DevCenter`: 2025-02-01
