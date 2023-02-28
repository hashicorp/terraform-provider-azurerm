---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_watchlist"
description: |-
  Manages a Sentinel Watchlist.
---

# azurerm_sentinel_watchlist

Manages a Sentinel Watchlist.

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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Sentinel Watchlist. Changing this forces a new Sentinel Watchlist to be created.

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace where this Sentinel Watchlist resides in. Changing this forces a new Sentinel Watchlist to be created.

* `display_name` - (Required) The display name of this Sentinel Watchlist. Changing this forces a new Sentinel Watchlist to be created.

* `item_search_key` - (Required) The key used to optimize query performance when using Watchlist for joins with other data. Changing this forces a new Sentinel Watchlist to be created.

---

* `default_duration` - (Optional) The default duration in ISO8601 duration form of this Sentinel Watchlist. Changing this forces a new Sentinel Watchlist to be created.

* `description` - (Optional) The description of this Sentinel Watchlist. Changing this forces a new Sentinel Watchlist to be created.

* `labels` - (Optional) Specifies a list of labels related to this Sentinel Watchlist. Changing this forces a new Sentinel Watchlist to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Sentinel Watchlist.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Sentinel Watchlist.
* `read` - (Defaults to 5 minutes) Used when retrieving the Sentinel Watchlist.
* `delete` - (Defaults to 30 minutes) Used when deleting the Sentinel Watchlist.

## Import

Sentinel Watchlists can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_watchlist.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/watchlists/list1
```
