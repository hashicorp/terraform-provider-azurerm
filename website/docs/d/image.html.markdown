---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_image"
sidebar_current: "docs-azurerm-datasource-image"
description: |-
  Get information about an Image

---

# azurerm_image

Use this data source to access information about an Image.

## Example Usage

```hcl
data "azurerm_image" "search" {
  name                = "search-api"
  resource_group_name = "packerimages"
}

output "image_id" {
  value = "${data.azurerm_image.search.id}"
}
```

## Argument Reference

* `name` - (Required) The name of the Image.
* `resource_group_name` - (Required) The Name of the Resource Group where this Image exists.

## Attributes Reference

* `location` - the Azure Location where this Image exists.
* `os_disk` - a `os_disk` block as defined below.
* `data_disk` - a collection of `data_disk` blocks as defined below.
* `tags` - a mapping of tags to assigned to the resource.

`os_disk` supports the following:

* `blob_uri` - the URI in Azure storage of the blob used to create the image.
* `caching` - the caching mode for the OS Disk, such as `ReadWrite`, `ReadOnly`, or `None`.
* `managed_disk_id` - the ID of the Managed Disk used as the OS Disk Image.
* `os_state` - the State of the OS used in the Image, such as `Generalized`.
* `os_type` - the type of Operating System used on the OS Disk. such as `Linux` or `Windows`.
* `size_gb` - the size of the OS Disk in GB.

`data_disk` supports the following:

* `blob_uri` - the URI in Azure storage of the blob used to create the image.
* `caching` - the caching mode for the Data Disk, such as `ReadWrite`, `ReadOnly`, or `None`.
* `lun` - the logical unit number of the data disk.
* `managed_disk_id` - the ID of the Managed Disk used as the Data Disk Image.
* `size_gb` - the size of this Data Disk in GB.
