---
subcategory: "Media"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_media_asset"
description: |-
  Manages a Media Asset.
---

# azurerm_media_asset

Manages a Media Asset.

## Example Usage

```hcl

resource "azurerm_resource_group" "example" {
  name     = "media-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoracc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "example" {
  name                = "examplemediaacc"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  storage_account {
    id         = azurerm_storage_account.example.id
    is_primary = true
  }
}

resource "azurerm_media_asset" "example" {
  name                        = "Asset1"
  resource_group_name         = azurerm_resource_group.example.name
  media_services_account_name = azurerm_media_services_account.example.name
  description                 = "Asset description"
}
```

## Arguments Reference

The following arguments are supported:

* `media_services_account_name` - (Required) Specifies the name of the Media Services Account. Changing this forces a new Media Asset to be created.

* `name` - (Required) The name which should be used for this Media Asset. Changing this forces a new Media Asset to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Media Asset should exist. Changing this forces a new Media Asset to be created.

---

* `alternate_id` - (Optional) The alternate ID of the Asset.

* `container` - (Optional) The name of the asset blob container. Changing this forces a new Media Asset to be created.

* `description` - (Optional) The Asset description.

* `storage_account_name` - (Optional) The name of the storage account where to store the media asset. Changing this forces a new Media Asset to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Media Asset.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Media Asset.
* `read` - (Defaults to 5 minutes) Used when retrieving the Media Asset.
* `update` - (Defaults to 30 minutes) Used when updating the Media Asset.
* `delete` - (Defaults to 30 minutes) Used when deleting the Media Asset.

## Import

Media Assets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_media_asset.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Media/mediaservices/account1/assets/asset1
```
