---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_image"
sidebar_current: "docs-azurerm-resource-compute-image"
description: |-
  Manages a custom virtual machine image that can be used to create virtual machines.
---

# azurerm_image

Manages a custom virtual machine image that can be used to create virtual machines.

## Example Usage Creating from VHD

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_image" "example" {
  name                = "acctest"
  location            = "West US"
  resource_group_name = "${azurerm_resource_group.example.name}"

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = "{blob_uri}"
    size_gb  = 30
  }
}
```

## Example Usage Creating from Virtual Machine (VM must be generalized beforehand)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_image" "example" {
  name                      = "acctest"
  location                  = "West US"
  resource_group_name       = "${azurerm_resource_group.example.name}"
  source_virtual_machine_id = "{vm_id}"
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
* `os_disk` - (Optional) One or more `os_disk` elements as defined below.
* `data_disk` - (Optional) One or more `data_disk` elements as defined below.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `zone_resilient` - (Optional) Is zone resiliency enabled?  Defaults to `false`.  Changing this forces a new resource to be created.
* `hyper_v_generation` - (Optional) The HyperVGenerationType of the VirtualMachine created from the image as `V1`, `V2`. The default is `V1`.

~> **Note**: `zone_resilient` can only be set to `true` if the image is stored in a region that supports availability zones.

`os_disk` supports the following:

* `os_type` - (Required) Specifies the type of operating system contained in the virtual machine image. Possible values are: Windows or Linux.
* `os_state` - (Required) Specifies the state of the operating system contained in the blob. Currently, the only value is Generalized.
* `managed_disk_id` - (Optional) Specifies the ID of the managed disk resource that you want to use to create the image.
* `blob_uri` - (Optional) Specifies the URI in Azure storage of the blob that you want to use to create the image.
* `caching` - (Optional) Specifies the caching mode as `ReadWrite`, `ReadOnly`, or `None`. The default is `None`.
* `size_gb` - (Optional) Specifies the size of the image to be created. The target size can't be smaller than the source size.

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
terraform import azurerm_image.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.compute/images/image1
```
