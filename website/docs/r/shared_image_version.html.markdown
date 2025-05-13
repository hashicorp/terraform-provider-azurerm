---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_shared_image_version"
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
  gallery_name        = data.azurerm_shared_image.existing.gallery_name
  image_name          = data.azurerm_shared_image.existing.name
  resource_group_name = data.azurerm_shared_image.existing.resource_group_name
  location            = data.azurerm_shared_image.existing.location
  managed_image_id    = data.azurerm_image.existing.id

  target_region {
    name                   = data.azurerm_shared_image.existing.location
    regional_replica_count = 5
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

* `target_region` - (Required) One or more `target_region` blocks as documented below.

* `blob_uri` - (Optional) URI of the Azure Storage Blob used to create the Image Version. Changing this forces a new resource to be created.

-> **Note:** You must specify exact one of `blob_uri`, `managed_image_id` and `os_disk_snapshot_id`.

-> **Note:** `blob_uri` and `storage_account_id` must be specified together

* `end_of_life_date` - (Optional) The end of life date in RFC3339 format of the Image Version.

* `exclude_from_latest` - (Optional) Should this Image Version be excluded from the `latest` filter? If set to `true` this Image Version won't be returned for the `latest` version. Defaults to `false`.

* `managed_image_id` - (Optional) The ID of the Managed Image or Virtual Machine ID which should be used for this Shared Image Version. Changing this forces a new resource to be created.

-> **Note:** The ID can be sourced from the `azurerm_image` [Data Source](https://www.terraform.io/docs/providers/azurerm/d/image.html) or [Resource](https://www.terraform.io/docs/providers/azurerm/r/image.html).

-> **Note:** You must specify exact one of `blob_uri`, `managed_image_id` and `os_disk_snapshot_id`.

* `os_disk_snapshot_id` - (Optional) The ID of the OS disk snapshot which should be used for this Shared Image Version. Changing this forces a new resource to be created.

-> **Note:** You must specify exact one of `blob_uri`, `managed_image_id` and `os_disk_snapshot_id`.

* `deletion_of_replicated_locations_enabled` - (Optional) Specifies whether this Shared Image Version can be deleted from the Azure Regions this is replicated to. Defaults to `false`. Changing this forces a new resource to be created.

* `replication_mode` - (Optional) Mode to be used for replication. Possible values are `Full` and `Shallow`. Defaults to `Full`. Changing this forces a new resource to be created.

* `storage_account_id` - (Optional) The ID of the Storage Account where the Blob exists. Changing this forces a new resource to be created.

-> **Note:** `blob_uri` and `storage_account_id` must be specified together

* `tags` - (Optional) A collection of tags which should be applied to this resource.

---

The `target_region` block supports the following:

* `name` - (Required) The Azure Region in which this Image Version should exist.

* `regional_replica_count` - (Required) The number of replicas of the Image Version to be created per region.

* `disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set to encrypt the Image Version in the target region. Changing this forces a new resource to be created.

* `exclude_from_latest_enabled` - (Optional) Specifies whether this Shared Image Version should be excluded when querying for the `latest` version. Defaults to `false`.

* `storage_account_type` - (Optional) The storage account type for the image version. Possible values are `Standard_LRS`, `Premium_LRS` and `Standard_ZRS`. Defaults to `Standard_LRS`. You can store all of your image version replicas in Zone Redundant Storage by specifying `Standard_ZRS`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Shared Image Version.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Shared Image Version.
* `read` - (Defaults to 5 minutes) Used when retrieving the Shared Image Version.
* `update` - (Defaults to 30 minutes) Used when updating the Shared Image Version.
* `delete` - (Defaults to 30 minutes) Used when deleting the Shared Image Version.

## Import

Shared Image Versions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_shared_image_version.version /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/galleries/gallery1/images/image1/versions/1.2.3
```
