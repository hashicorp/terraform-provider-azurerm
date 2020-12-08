---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_shared_image_versions"
description: |-
  Gets information about existing Versions of a Shared Image within a Shared Image Gallery.

---

# Data Source: azurerm_shared_image_versions

Use this data source to access information about existing Versions of a Shared Image within a Shared Image Gallery.

## Example Usage

```hcl
data "azurerm_shared_image_versions" "example" {
  image_name          = "my-image"
  gallery_name        = "my-image-gallery"
  resource_group_name = "example-resources"
}
```

## Argument Reference

The following arguments are supported:

* `image_name` - The name of the Shared Image in which this Version exists.

* `gallery_name` - The name of the Shared Image in which the Shared Image exists.

* `resource_group_name` - The name of the Resource Group in which the Shared Image Gallery exists.

* `tags_filter` - A mapping of tags to filter the list of images against.

## Attributes Reference

The following attributes are exported:

* `images` - An `images` block as defined below:

---

A `images` block exports the following:

* `exclude_from_latest` - Is this Image Version excluded from the `latest` filter?

* `location` - The supported Azure location where the Shared Image Gallery exists.

* `managed_image_id` - The ID of the Managed Image which was the source of this Shared Image Version.

* `target_region` - One or more `target_region` blocks as documented below.

* `tags` - A mapping of tags assigned to the Shared Image.

---

The `target_region` block exports the following:

* `name` - The Azure Region in which this Image Version exists.

* `regional_replica_count` - The number of replicas of the Image Version to be created per region.

* `storage_account_type` - The storage account type for the image version.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Versions of a Shared Image within a Shared Image Gallery.
