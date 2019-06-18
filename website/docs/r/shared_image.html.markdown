---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_shared_image"
sidebar_current: "docs-azurerm-resource-compute-shared-image-x"
description: |-
  Manages a Shared Image within a Shared Image Gallery.

---

# azurerm_shared_image

Manages a Shared Image within a Shared Image Gallery.

-> **NOTE** Shared Image Galleries are currently in Public Preview. You can find more information, including [how to register for the Public Preview here](https://azure.microsoft.com/en-gb/blog/announcing-the-public-preview-of-shared-image-gallery/).

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "example_image_gallery"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  description         = "Shared images and things."

  tags = {
    Hello = "There"
    World = "Example"
  }
}

resource "azurerm_shared_image" "test" {
  name                = "my-image"
  gallery_name        = "${azurerm_shared_image_gallery.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
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

* `os_type` - (Required) The type of Operating System present in this Shared Image. Possible values are `Linux` and `Windows`.

---

* `description` - (Optional) A description of this Shared Image.

* `eula` - (Optional) The End User Licence Agreement for the Shared Image.

* `privacy_statement_uri` - (Optional) The URI containing the Privacy Statement associated with this Shared Image.

* `release_note_uri` - (Optional) The URI containing the Release Notes associated with this Shared Image.

* `tags` - (Optional) A mapping of tags to assign to the Shared Image.

---

A `identifier` block supports the following:

* `offer` - (Required) The Offer Name for this Shared Image.

* `publisher` - (Required) The Publisher Name for this Gallery Image.

* `sku` - (Required) The Name of the SKU for this Gallery Image.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Shared Image.

## Import

Shared Images can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_shared_image.image1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/galleries/gallery1/images/image1
```
