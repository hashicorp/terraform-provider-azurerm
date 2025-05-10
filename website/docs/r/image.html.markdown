---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_image"
description: |-
  Manages a custom virtual machine image that can be used to create virtual machines.
---

# azurerm_image

Manages a custom virtual machine image that can be used to create virtual machines.

## Example Usage

-> **Note:** For a more complete example, see [the `examples/image` directory](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/image) within the GitHub Repository.

```hcl
data "azurerm_virtual_machine" "example" {
  name                = "examplevm"
  resource_group_name = "example-resources"
}

resource "azurerm_image" "example" {
  name                      = "exampleimage"
  location                  = data.azurerm_virtual_machine.example.location
  resource_group_name       = data.azurerm_virtual_machine.example.name
  source_virtual_machine_id = data.azurerm_virtual_machine.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the image. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which to create the image. Changing this forces a new resource to be created.
* `location` - (Required) Specified the supported Azure location where the resource exists. Changing this forces a new resource to be created.
* `source_virtual_machine_id` - (Optional) The Virtual Machine ID from which to create the image.
* `os_disk` - (Optional) One or more `os_disk` blocks as defined below. Changing this forces a new resource to be created.
* `data_disk` - (Optional) One or more `data_disk` blocks as defined below.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `zone_resilient` - (Optional) Is zone resiliency enabled? Defaults to `false`. Changing this forces a new resource to be created.
* `hyper_v_generation` - (Optional) The HyperVGenerationType of the VirtualMachine created from the image as `V1`, `V2`. Defaults to `V1`. Changing this forces a new resource to be created.

~> **Note:** `zone_resilient` can only be set to `true` if the image is stored in a region that supports availability zones.

---

The `os_disk` block supports the following:

* `storage_type` - (Required) The type of Storage Disk to use. Possible values are `Premium_LRS`, `PremiumV2_LRS`, `Premium_ZRS`, `Standard_LRS`, `StandardSSD_LRS`, `StandardSSD_ZRS` and `UltraSSD_LRS`. Changing this forces a new resource to be created.
* `os_type` - (Optional) Specifies the type of operating system contained in the virtual machine image. Possible values are: `Windows` or `Linux`.
* `os_state` - (Optional) Specifies the state of the operating system contained in the blob. Currently, the only value is Generalized. Possible values are `Generalized` and `Specialized`.
* `managed_disk_id` - (Optional) Specifies the ID of the managed disk resource that you want to use to create the image.
* `blob_uri` - (Optional) Specifies the URI in Azure storage of the blob that you want to use to create the image. Changing this forces a new resource to be created.
* `caching` - (Optional) Specifies the caching mode as `ReadWrite`, `ReadOnly`, or `None`. The default is `None`.
* `size_gb` - (Optional) Specifies the size of the image to be created. Changing this forces a new resource to be created.
* `disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set which should be used to encrypt this disk. Changing this forces a new resource to be created.

---

The `data_disk` block supports the following:

* `storage_type` - (Required) The type of Storage Disk to use. Possible values are `Premium_LRS`, `PremiumV2_LRS`, `Premium_ZRS`, `Standard_LRS`, `StandardSSD_LRS`, `StandardSSD_ZRS` and `UltraSSD_LRS`. Changing this forces a new resource to be created.
* `lun` - (Optional) Specifies the logical unit number of the data disk.
* `managed_disk_id` - (Optional) Specifies the ID of the managed disk resource that you want to use to create the image. Changing this forces a new resource to be created.
* `blob_uri` - (Optional) Specifies the URI in Azure storage of the blob that you want to use to create the image.
* `caching` - (Optional) Specifies the caching mode as `ReadWrite`, `ReadOnly`, or `None`. Defaults to `None`.
* `size_gb` - (Optional) Specifies the size of the image to be created. The target size can't be smaller than the source size.
* `disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set which should be used to encrypt this disk. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Image.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Image.
* `read` - (Defaults to 5 minutes) Used when retrieving the Image.
* `update` - (Defaults to 90 minutes) Used when updating the Image.
* `delete` - (Defaults to 90 minutes) Used when deleting the Image.

## Import

Images can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_image.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/images/image1
```
