---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_center_gallery"
description: |-
  Manages a Dev Center Gallery.
---

# azurerm_dev_center_gallery

Manages a Dev Center Gallery.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dev_center" "test" {
  name                = "example-devcenter"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "example-uai"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_shared_image_gallery" "example" {
  name                = "example-image-gallery"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_dev_center_gallery" "example" {
  dev_center_id     = azurerm_dev_center.example.id
  shared_gallery_id = azurerm_shared_image_gallery.example.id
  name              = "example"
}
```

## Arguments Reference

The following arguments are supported:

* `dev_center_id` - (Required) Specifies the ID of the Dev Center within which this Dev Center Gallery should exist. Changing this forces a new Dev Center Gallery to be created.

* `shared_gallery_id` - (Required) The ID of the Shared Gallery which should be connected to the Dev Center Gallery. Changing this forces a new Dev Center Gallery to be created.

* `name` - (Required) Specifies the name of this Dev Center Gallery. Changing this forces a new Dev Center Gallery to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center Gallery.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dev Center Gallery.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Center Gallery.
* `delete` - (Defaults to 30 minutes) Used when deleting the Dev Center Gallery.

## Import

An existing Dev Center Gallery can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_dev_center_gallery.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevCenter/devCenters/{devCenterName}/galleries/{galleryName}
```

* Where `{subscriptionId}` is the ID of the Azure Subscription where the Dev Center Gallery exists. For example `12345678-1234-9876-4563-123456789012`.
* Where `{resourceGroupName}` is the name of Resource Group where this Dev Center Gallery exists. For example `example-resource-group`.
* Where `{devCenterName}` is the name of the Dev Center. For example `devCenterValue`.
* Where `{galleryName}` is the name of the Gallery. For example `galleryValue`.
