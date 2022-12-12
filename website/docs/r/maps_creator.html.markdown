---
subcategory: "Maps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_maps_creator"
description: |-
  Manages an Azure Maps Creator.
---

# azurerm_maps_creator

Manages an Azure Maps Creator.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_maps_account" "example" {
  name                = "example-maps-account"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "G2"

  tags = {
    environment = "Test"
  }
}

resource "azurerm_maps_creator" "example" {
  name            = "example-maps-creator"
  maps_account_id = azurerm_maps_account.example.id
  location        = azurerm_resource_group.example.location
  storage_units   = 1

  tags = {
    environment = "Test"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Azure Maps Creator. Changing this forces a new resource to be created.

* `maps_account_id` - (Required) The ID of the Azure Maps Creator. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Azure Maps Creator should exist. Changing this forces a new resource to be created.

* `storage_units` - (Required) The storage units to be allocated. Integer values from 1 to 100, inclusive.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Maps Creator.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Maps Creator.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Maps Creator.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Maps Creator.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Maps Creator.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Maps Creator.

## Import

An Azure Maps Creators can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_maps_creator.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Maps/accounts/account1/creators/creator1
```
