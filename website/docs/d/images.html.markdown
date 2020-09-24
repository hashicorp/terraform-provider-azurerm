---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_images"
description: |-
  Gets information about existing Images within a Resource Group.

---

# Data Source: azurerm_images

Use this data source to access information about existing Images within a Resource Group.

## Example Usage

```hcl
data "azurerm_images" "example" {
  resource_group_name = "example-resources"
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - The name of the Resource Group in which the Image exists.

* `tags_filter` - A mapping of tags to filter the list of images against.

## Attributes Reference

The following attributes are exported:

* `images` - One or more `images` blocks as defined below:

---

A `images` block exports the following:

* `name` - The name of the Image.

* `location` - The supported Azure location where the Image exists.

* `zone_resilient` - Is zone resiliency enabled?

* `os_disk` - An `os_disk` block as defined below.

* `data_disk` - One or more `data_disk` blocks as defined below.

* `tags` - A mapping of tags assigned to the Image.

---

The `os_disk` block exports the following:

* `blob_uri` - the URI in Azure storage of the blob used to create the image.

* `caching` - the caching mode for the OS Disk.

* `managed_disk_id` - the ID of the Managed Disk used as the OS Disk Image.

* `os_state` - the State of the OS used in the Image.

* `os_type` - the type of Operating System used on the OS Disk.

* `size_gb` - the size of the OS Disk in GB.

---

The `data_disk` block exports the following:

* `blob_uri` - the URI in Azure storage of the blob used to create the image.

* `caching` - the caching mode for the Data Disk.

* `lun` - the logical unit number of the data disk.

* `managed_disk_id` - the ID of the Managed Disk used as the Data Disk Image.

* `size_gb` - the size of this Data Disk in GB.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Images within a Resource Group.
