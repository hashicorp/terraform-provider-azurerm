---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_shared_image_version"
sidebar_current: "docs-azurerm-resource-compute-shared-image-version"
description: |-
  Manages a Version of a Shared Image within a Shared Image Gallery.

---

# azurerm_shared_image_version

Manages a Version of a Shared Image within a Shared Image Gallery.

## Example Usage

```hcl
data "azurerm_image" "existing" {
  name                = "search-api"
  resource_group_name = "packerimages"
}

data "azurerm_shared_image" "existing" {
  name                = "existing-image"
  gallery_name        = "existing_gallery"
  resource_group_name = "existing-resources"
}

resource "azurerm_shared_image_version" "example" {
  name                = "0.0.1"
  gallery_name        = "${data.azurerm_shared_image.existing.gallery_name}"
  image_name          = "${data.azurerm_shared_image.existing.name}"
  resource_group_name = "${data.azurerm_shared_image.existing.resource_group_name}"
  location            = "${data.azurerm_shared_image.existing.location}"
  managed_image_id    = "${data.azurerm_image.existing.id}"

  target_region {
    name                   = "${data.azurerm_shared_image.existing.location}"
    regional_replica_count = "5"
    storage_account_type   = "Standard_LRS"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The version number for this Image Version, such as `1.0.0`. Changing this forces a new resource to be created.

* `gallery_name` - (Required) The name of the Shared Image Gallery in which the Shared Image exists. Changing this forces a new resource to be created.

* `image_name` - (Required) The name of the Shared Image within the Shared Image Gallery in which this Version should be created. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region in which the Shared Image Gallery exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Shared Image Gallery exists. Changing this forces a new resource to be created.

* `managed_image_id` - (Required) The ID of the Managed Image which should be used for this Shared Image Version. Changing this forces a new resource to be created.

-> **NOTE:** The ID can be sourced from the `azurerm_image` [Data Source](https://www.terraform.io/docs/providers/azurerm/d/image.html) or [Resource](https://www.terraform.io/docs/providers/azurerm/r/image.html).

* `target_region` - (Required) One or more `target_region` blocks as documented below.

* `exclude_from_latest` - (Optional) Should this Image Version be excluded from the `latest` filter? If set to `true` this Image Version won't be returned for the `latest` version. Defaults to `false`.

* `tags` - (Optional) A collection of tags which should be applied to this resource.

---

The `target_region` block exports the following:

* `name` - (Required) The Azure Region in which this Image Version should exist.

* `regional_replica_count` - (Required) The number of replicas of the Image Version to be created per region.

* `storage_account_type` - (Optional) The storage account type for the image version, which defaults to `Standard_LRS`. You can store all of your image version replicas in Zone Redundant Storage by specifying `Standard_ZRS`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Shared Image Version.

## Import

Shared Image Versions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_shared_image_version.version /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/galleries/gallery1/images/image1/versions/1.2.3
```
