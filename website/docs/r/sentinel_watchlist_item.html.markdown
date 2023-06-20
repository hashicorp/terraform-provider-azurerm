---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_watchlist_item"
description: |-
  Manages a Sentinel Watchlist Item.
---

# azurerm_sentinel_watchlist_item

Manages a Sentinel Watchlist Item.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "example" {
  workspace_id = azurerm_log_analytics_workspace.example.id
}

resource "azurerm_sentinel_watchlist" "example" {
  name                       = "example-watchlist"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.example.workspace_id
  display_name               = "example-wl"
  item_search_key            = "Key"
}

resource "azurerm_sentinel_watchlist_item" "example" {
  name         = "0aac6fa5-223e-49cf-9bfd-3554dc9d2b76"
  watchlist_id = azurerm_sentinel_watchlist.example.id
  properties = {
    k1 = "v1"
    k2 = "v2"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `watchlist_id` - (Required) The ID of the Sentinel Watchlist that this Item resides in. Changing this forces a new Sentinel Watchlist Item to be created.

* `properties` - (Required) The key value pairs of the Sentinel Watchlist Item.

---

* `name` - (Optional) The name in UUID format which should be used for this Sentinel Watchlist Item. Changing this forces a new Sentinel Watchlist Item to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Sentinel Watchlist Item.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Sentinel Watchlist Item.
* `read` - (Defaults to 5 minutes) Used when retrieving the Sentinel Watchlist Item.
* `update` - (Defaults to 30 minutes) Used when updating the Sentinel Watchlist Item.
* `delete` - (Defaults to 30 minutes) Used when deleting the Sentinel Watchlist Item.

## Import

Sentinel Watchlist Items can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_watchlist_item.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/watchlists/list1/watchlistItems/item1
```
