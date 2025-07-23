---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_shared_image_gallery"
description: |-
  Manages a Shared Image Gallery.

---

# azurerm_shared_image_gallery

Manages a Shared Image Gallery.

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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Shared Image Gallery. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Shared Image Gallery. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `description` - (Optional) A description for this Shared Image Gallery.

* `sharing` - (Optional) A `sharing` block as defined below. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the Shared Image Gallery.

---

A `sharing` block supports the following:

* `permission` - (Required) The permission of the Shared Image Gallery when sharing. Possible values are `Community`, `Groups` and `Private`. Changing this forces a new resource to be created.

-> **Note:** This requires that the Preview Feature `Microsoft.Compute/CommunityGalleries` is enabled, see [the documentation](https://learn.microsoft.com/azure/virtual-machines/share-gallery-community?tabs=cli) for more information.

* `community_gallery` - (Optional) A `community_gallery` block as defined below. Changing this forces a new resource to be created.

~> **Note:** `community_gallery` must be set when `permission` is set to `Community`.

---

A `community_gallery` block supports the following:

* `eula` - (Required) The End User Licence Agreement for the Shared Image Gallery. Changing this forces a new resource to be created.

* `prefix` - (Required) Prefix of the community public name for the Shared Image Gallery. Changing this forces a new resource to be created.

* `publisher_email` - (Required) Email of the publisher for the Shared Image Gallery. Changing this forces a new resource to be created.

* `publisher_uri` - (Required) URI of the publisher for the Shared Image Gallery. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Shared Image Gallery.

* `unique_name` - The Unique Name for this Shared Image Gallery.

---

A `community_gallery` block exports the following:

* `name` - The community public name of the Shared Image Gallery.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Shared Image Gallery.
* `read` - (Defaults to 5 minutes) Used when retrieving the Shared Image Gallery.
* `update` - (Defaults to 30 minutes) Used when updating the Shared Image Gallery.
* `delete` - (Defaults to 30 minutes) Used when deleting the Shared Image Gallery.

## Import

Shared Image Galleries can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_shared_image_gallery.gallery1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/galleries/gallery1
```
