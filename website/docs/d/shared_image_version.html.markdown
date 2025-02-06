---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_shared_image_version"
description: |-
  Gets information about an existing Version of a Shared Image within a Shared Image Gallery.

---

# Data Source: azurerm_shared_image_version

Use this data source to access information about an existing Version of a Shared Image within a Shared Image Gallery.

## Example Usage

```hcl
data "azurerm_shared_image_version" "example" {
  name                = "1.0.0"
  image_name          = "my-image"
  gallery_name        = "my-image-gallery"
  resource_group_name = "example-resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Image Version.

~> **Note:** You may specify `latest` to obtain the latest version or `recent` to obtain the most recently updated version.

~> **Note:** In 3.0, `latest` may return an image version with `exclude_from_latest` set to `true`. Starting from 4.0 onwards `latest` will not return image versions with `exlude_from_latest` set to `true`.

* `image_name` - The name of the Shared Image in which this Version exists.

* `gallery_name` - The name of the Shared Image Gallery in which the Shared Image exists.

* `resource_group_name` - The name of the Resource Group in which the Shared Image Gallery exists.

* `sort_versions_by_semver` - (Optional) Sort available versions taking SemVer versioning scheme into account. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Shared Image.

* `exclude_from_latest` - Is this Image Version excluded from the `latest` filter?

* `location` - The supported Azure location where the Shared Image Gallery exists.

* `managed_image_id` - The ID of the Managed Image which was the source of this Shared Image Version.

* `target_region` - One or more `target_region` blocks as documented below.

* `os_disk_snapshot_id` - The ID of the OS disk snapshot which was the source of this Shared Image Version.

* `os_disk_image_size_gb` - The size of the OS disk snapshot (in Gigabytes) which was the source of this Shared Image Version.

* `tags` - A mapping of tags assigned to the Shared Image.

---

The `target_region` block exports the following:

* `name` - The Azure Region in which this Image Version exists.

* `regional_replica_count` - The number of replicas of the Image Version to be created per region.

* `storage_account_type` - The storage account type for the image version.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Version of a Shared Image within a Shared Image Gallery.
