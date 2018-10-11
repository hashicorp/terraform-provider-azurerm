---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_shared_image_version"
sidebar_current: "docs-azurerm-datasource-shared-image-version"
description: |-
  Gets information about an existing Version of a Shared Image within a Shared Image Gallery.

---

# Data Source: azurerm_shared_image_version

Use this data source to access information about an existing Version of a Shared Image within a Shared Image Gallery.

-> **NOTE** Shared Image Galleries are currently in Public Preview. You can find more information, including [how to register for the Public Preview here](https://azure.microsoft.com/en-gb/blog/announcing-the-public-preview-of-shared-image-gallery/).

## Example Usage

```hcl
data "azurerm_shared_image_version" "test" {
  name                = "1.0.0"
  image_name          = "my-image"
  gallery_name        = "my-image-gallery"
  resource_group_name = "example-resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Image Version.

* `image_name` - (Required) The name of the Shared Image in which this Version exists.

* `gallery_name` - (Required) The name of the Shared Image in which the Shared Image exists.

* `resource_group_name` - (Required) The name of the Resource Group in which the Shared Image Gallery exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Resource ID of the Shared Image.

* `exclude_from_latest` - Is this Image Version excluded from the `latest` filter?

* `location` - The supported Azure location where the Shared Image Gallery exists.

* `managed_image_id` - The ID of the Managed Image which was the source of this Shared Image Version.

* `regions` - The list of Azure Regions in which this Image Version exists.

* `tags` - A mapping of tags assigned to the Shared Image.
