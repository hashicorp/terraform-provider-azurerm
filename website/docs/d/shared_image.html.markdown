---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_shared_image"
sidebar_current: "docs-azurerm-datasource-shared-image-x"
description: |-
  Gets information about the specified Shared Image within the Shared Image Gallery.

---

# Data Source: azurerm_shared_image

Gets information about the specified Shared Image within the Shared Image Gallery.

## Example Usage

```hcl
data "azurerm_shared_image" "test" {
  name                = "my-image"
  gallery_name        = "my-image-gallery"
  resource_group_name = "example-resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Shared Image.

* `gallery_name` - (Required) The name of the Shared Image Gallery in which the Shared Image exists.

* `resource_group_name` - (Required) The name of the Resource Group in which the Shared Image Gallery exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Resource ID of the Shared Image.

* `description` - The description of this Shared Image.

* `eula` - The End User Licence Agreement for the Shared Image.

* `location` - The supported Azure location where the Shared Image Gallery exists.

* `identity` - An `identity` block as defined below.

* `os_type` - The type of Operating System present in this Shared Image.

* `privacy_statement_uri` - The URI containing the Privacy Statement for this Shared Image.

* `release_note_uri` - The URI containing the Release Notes for this Shared Image.

* `tags` - A mapping of tags assigned to the Shared Image.

---

A `identity` block exports the following:

* `offer` - The Offer Name for this Shared Image.

* `publisher` - The Publisher Name for this Gallery Image.

* `sku` - The Name of the SKU for this Gallery Image.