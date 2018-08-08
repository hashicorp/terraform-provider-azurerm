---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_image"
sidebar_current: "docs-azurerm-resource-compute-image"
description: |-
  Manages a custom virtual machine image that can be used to create virtual machines.
---

# azurerm_image

Create a custom virtual machine image that can be used to create virtual machines.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_virtual_machine" "example" {
  # ...
}

resource "azurerm_image" "example" {
  name                      = "example-image"
  location                  = "${azurerm_resource_group.example.location}"
  resource_group_name       = "${azurerm_resource_group.example.name}"
  source_virtual_machine_id = "${azurerm_virtual_machine.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the image. Changing this forces a
    new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which to create
    the image. Changing this forces a new resource to be created.
* `location` - (Required) Specified the supported Azure location where the resource exists.
    Changing this forces a new resource to be created.
* `source_virtual_machine_id` - (Optional) The Virtual Machine ID from which to create the image.

-> **NOTE:** The Virtual Machine must be Generalized prior to an image being taken - more information can be found for both [Linux](https://docs.microsoft.com/en-us/azure/virtual-machines/linux/capture-image) and [Windows](https://docs.microsoft.com/en-gb/azure/virtual-machines/windows/sa-copy-generalized) respectively.

* `os_disk` - (Optional) One or more `os_disk` elements as defined below.
* `data_disk` - (Optional) One or more `data_disk` elements as defined below.
* `tags` - (Optional) A mapping of tags to assign to the resource.

`os_disk` supports the following:

* `os_type` - (Required) Specifies the type of operating system contained in the the virtual machine image. Possible values are: Windows or Linux.
* `os_state` - (Required) Specifies the state of the operating system contained in the blob. Currently, the only value is Generalized.
* `managed_disk_id` - (Optional) Specifies the ID of the managed disk resource that you want to use to create the image.
* `blob_uri` - (Optional) Specifies the URI in Azure storage of the blob that you want to use to create the image.
* `caching` - (Optional) Specifies the caching mode as `ReadWrite`, `ReadOnly`, or `None`. The default is `None`.

`data_disk` supports the following:

* `lun` - (Required) Specifies the logical unit number of the data disk.
* `managed_disk_id` - (Optional) Specifies the ID of the managed disk resource that you want to use to create the image.
* `blob_uri` - (Optional) Specifies the URI in Azure storage of the blob that you want to use to create the image.
* `caching` - (Optional) Specifies the caching mode as `ReadWrite`, `ReadOnly`, or `None`. The default is `None`.
* `size_gb` - (Optional) Specifies the size of the image to be created. The target size can't be smaller than the source size.

## Attributes Reference

The following attributes are exported:

* `id` - The managed image ID.

## Import

Image can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_image.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.compute/images/image1
```
