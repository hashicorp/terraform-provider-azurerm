---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_shared_image"
description: |-
  Manages a Shared Image within a Shared Image Gallery.

---

# azurerm_shared_image

Manages a Shared Image within a Shared Image Gallery.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_shared_image_gallery" "example" {
  name                = "example_image_gallery"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  description         = "Shared images and things."

  tags = {
    Hello = "There"
    World = "Example"
  }
}

resource "azurerm_shared_image" "example" {
  name                = "my-image"
  gallery_name        = azurerm_shared_image_gallery.example.name
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  os_type             = "Linux"

  identifier {
    publisher = "PublisherName"
    offer     = "OfferName"
    sku       = "ExampleSku"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Shared Image. Changing this forces a new resource to be created.

* `gallery_name` - (Required) Specifies the name of the Shared Image Gallery in which this Shared Image should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Shared Image Gallery exists. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the Shared Image Gallery exists. Changing this forces a new resource to be created.

* `identifier` - (Required) An `identifier` block as defined below.

* `os_type` - (Required) The type of Operating System present in this Shared Image. Possible values are `Linux` and `Windows`. Changing this forces a new resource to be created.

* `purchase_plan` - (Optional) A `purchase_plan` block as defined below.

---

* `description` - (Optional) A description of this Shared Image.

* `disk_types_not_allowed` - (Optional) One or more Disk Types not allowed for the Image. Possible values include `Standard_LRS` and `Premium_LRS`.

* `end_of_life_date` - (Optional) The end of life date in RFC3339 format of the Image.

* `eula` - (Optional) The End User Licence Agreement for the Shared Image. Changing this forces a new resource to be created.

* `specialized` - (Optional) Specifies that the Operating System used inside this Image has not been Generalized (for example, `sysprep` on Windows has not been run). Changing this forces a new resource to be created.

!> **Note:** It's recommended to Generalize images where possible - Specialized Images reuse the same UUID internally within each Virtual Machine, which can have unintended side-effects.

* `architecture` - (Optional) CPU architecture supported by an OS. Possible values are `x64` and `Arm64`. Defaults to `x64`. Changing this forces a new resource to be created.

* `hyper_v_generation` - (Optional) The generation of HyperV that the Virtual Machine used to create the Shared Image is based on. Possible values are `V1` and `V2`. Defaults to `V1`. Changing this forces a new resource to be created.

* `max_recommended_vcpu_count` - (Optional) Maximum count of vCPUs recommended for the Image.

* `min_recommended_vcpu_count` - (Optional) Minimum count of vCPUs recommended for the Image.

* `max_recommended_memory_in_gb` - (Optional) Maximum memory in GB recommended for the Image.

* `min_recommended_memory_in_gb` - (Optional) Minimum memory in GB recommended for the Image.

* `privacy_statement_uri` - (Optional) The URI containing the Privacy Statement associated with this Shared Image. Changing this forces a new resource to be created.

* `release_note_uri` - (Optional) The URI containing the Release Notes associated with this Shared Image.

* `trusted_launch_supported` - (Optional) Specifies if supports creation of both Trusted Launch virtual machines and Gen2 virtual machines with standard security created from the Shared Image. Changing this forces a new resource to be created.

* `trusted_launch_enabled` - (Optional) Specifies if Trusted Launch has to be enabled for the Virtual Machine created from the Shared Image. Changing this forces a new resource to be created.

* `confidential_vm_supported` - (Optional) Specifies if supports creation of both Confidential virtual machines and Gen2 virtual machines with standard security from a compatible Gen2 OS disk VHD or Gen2 Managed image. Changing this forces a new resource to be created.

* `confidential_vm_enabled` - (Optional) Specifies if Confidential Virtual Machines enabled. It will enable all the features of trusted, with higher confidentiality features for isolate machines or encrypted data. Available for Gen2 machines. Changing this forces a new resource to be created.

-> **Note:**: Only one of `trusted_launch_supported`, `trusted_launch_enabled`, `confidential_vm_supported` and `confidential_vm_enabled` can be specified.

* `accelerated_network_support_enabled` - (Optional) Specifies if the Shared Image supports Accelerated Network. Changing this forces a new resource to be created.

* `hibernation_enabled` - (Optional) Specifies if the Shared Image supports hibernation. Changing this forces a new resource to be created.

* `disk_controller_type_nvme_enabled` - (Optional) Specifies if the Shared Image supports NVMe disks. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the Shared Image.

---

A `identifier` block supports the following:

* `offer` - (Required) The Offer Name for this Shared Image. Changing this forces a new resource to be created.

* `publisher` - (Required) The Publisher Name for this Gallery Image. Changing this forces a new resource to be created.

* `sku` - (Required) The Name of the SKU for this Gallery Image. Changing this forces a new resource to be created.

---

A `purchase_plan` block supports the following:

* `name` - (Required) The Purchase Plan Name for this Shared Image. Changing this forces a new resource to be created.

* `publisher` - (Optional) The Purchase Plan Publisher for this Gallery Image. Changing this forces a new resource to be created.

* `product` - (Optional) The Purchase Plan Product for this Gallery Image. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Shared Image.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Shared Image.
* `update` - (Defaults to 30 minutes) Used when updating the Shared Image.
* `read` - (Defaults to 5 minutes) Used when retrieving the Shared Image.
* `delete` - (Defaults to 30 minutes) Used when deleting the Shared Image.

## Import

Shared Images can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_shared_image.image1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/galleries/gallery1/images/image1
```
