---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_shared_image"
description: |-
  Gets information about an existing Shared Image a Shared Image Gallery.

---

# Data Source: azurerm_shared_image

Use this data source to access information about an existing Shared Image within a Shared Image Gallery.

## Example Usage

```hcl
data "azurerm_shared_image" "example" {
  name                = "my-image"
  gallery_name        = "my-image-gallery"
  resource_group_name = "example-resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Shared Image.

* `gallery_name` - The name of the Shared Image Gallery in which the Shared Image exists.

* `resource_group_name` - The name of the Resource Group in which the Shared Image Gallery exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Shared Image.

* `description` - The description of this Shared Image.

* `eula` - The End User Licence Agreement for the Shared Image.

* `location` - The supported Azure location where the Shared Image Gallery exists.

* `specialized` - Specifies that the Operating System used inside this Image has not been Generalized (for example, `sysprep` on Windows has not been run).

* `identifier` - An `identifier` block as defined below.

* `os_type` - The type of Operating System present in this Shared Image.

* `hyper_v_generation` - The generation of HyperV that the Virtual Machine used to create the Shared Image is based on.

* `privacy_statement_uri` - The URI containing the Privacy Statement for this Shared Image.

* `release_note_uri` - The URI containing the Release Notes for this Shared Image.

* `tags` - A mapping of tags assigned to the Shared Image.

---

A `identifier` block exports the following:

* `offer` - The Offer Name for this Shared Image.

* `publisher` - The Publisher Name for this Gallery Image.

* `sku` - The Name of the SKU for this Gallery Image.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Shared Image a Shared Image Gallery.
