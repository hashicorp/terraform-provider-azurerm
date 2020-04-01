---
subcategory: "TODO - pick from: Load Balancer|Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_route_filter"
description: |-
  Manages a Route Filter.
---

# azurerm_route_filter

Manages a Route Filter.

## Example Usage

```hcl
resource "azurerm_route_filter" "example" {
  name = "example"
  resource_group_name = "example"
  location = "West Europe"
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Route Filter should exist. Changing this forces a new Route Filter to be created.

* `name` - (Required) The Name which should be used for this Route Filter. Changing this forces a new Route Filter to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Route Filter should exist. Changing this forces a new Route Filter to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Route Filter.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Route Filter.
* `read` - (Defaults to 5 minutes) Used when retrieving the Route Filter.
* `update` - (Defaults to 30 minutes) Used when updating the Route Filter.
* `delete` - (Defaults to 30 minutes) Used when deleting the Route Filter.

## Import

Route Filters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_route_filter.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/routeFilters/routeFilter1
```