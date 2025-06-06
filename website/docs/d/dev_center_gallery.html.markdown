---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_dev_center_gallery"
description: |-
  Gets information about an existing Dev Center Gallery.
---

# Data Source: azurerm_dev_center_gallery

Use this data source to access information about an existing Dev Center Gallery.

## Example Usage

```hcl
data "azurerm_dev_center_gallery" "example" {
  name          = azurerm_dev_center_gallery.example.name
  dev_center_id = azurerm_dev_center_gallery.example.dev_center_id
}

output "id" {
  value = data.azurerm_dev_center_gallery.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Dev Center Gallery.

* `dev_center_id` - (Required) The ID of the Dev Center within which this Dev Center Gallery should exist.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center Gallery.

* `shared_gallery_id` - The ID of the Shared Gallery which should be connected to the Dev Center Gallery.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Center Gallery.
