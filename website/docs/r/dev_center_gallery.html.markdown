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

resource "azurem_user_assigned_identity" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_dev_center" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  identity {
    type = "UserAssigned"
    identity_ids = [
      azurem_user_assigned_identity.example.id,
    ]
  }
}

resource "azurerm_shared_image_gallery" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_shared_image_gallery.example.id
  role_definition_name = "Owner"
  principal_id         = azurem_user_assigned_identity.example.principal_id
}

resource "azurerm_dev_center_gallery" "example" {
  resource_group_name = azurerm_resource_group.example.name
  name                = "example"
  dev_center_name     = azurerm_dev_center.example.name
  gallery_resource_id = azurerm_shared_image_gallery.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `dev_center_name` - (Required) Resource name of an associated DevCenter. Changing this forces a new Dev Center Gallery to be created.

* `name` - (Required) Specifies the name of this Dev Center Gallery. Changing this forces a new Dev Center Gallery to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group within which this Dev Center Gallery should exist. Changing this forces a new Dev Center Gallery to be created.

* `gallery_resource_id` - (Required) Specifies the resource ID of the Shared Image Gallery to be associated with this Dev Center Gallery. Changing this forces a new Dev Center Gallery to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center Gallery.

---


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating this Dev Center Gallery.
* `delete` - (Defaults to 30 minutes) Used when deleting this Dev Center Gallery.
* `read` - (Defaults to 5 minutes) Used when retrieving this Dev Center Gallery.
* `update` - (Defaults to 30 minutes) Used when updating this Dev Center Gallery.

## Import

An existing Dev Center Gallery can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_dev_center_gallery.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevCenter/devcenters/{devcentersName}/galleries/{galleryName}
```
```

* Where `{subscriptionId}` is the ID of the Azure Subscription where the Dev Center Gallery exists. For example `12345678-1234-9876-4563-123456789012`.
* Where `{resourceGroupName}` is the name of Resource Group where this Dev Center Gallery exists. For example `example-resource-group`.
* Where `{galleryName}` is the name of the Gallery. For example `galleryValue`.
