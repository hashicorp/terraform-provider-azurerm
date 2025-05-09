---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_route_map"
description: |-
  Manages a Route Map.
---

# azurerm_route_map

Manages a Route Map.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_wan" "example" {
  name                = "example-vwan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_hub" "example" {
  name                = "example-vhub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_wan_id      = azurerm_virtual_wan.example.id
  address_prefix      = "10.0.1.0/24"
}

resource "azurerm_route_map" "example" {
  name           = "example-rm"
  virtual_hub_id = azurerm_virtual_hub.example.id

  rule {
    name                 = "rule1"
    next_step_if_matched = "Continue"

    action {
      type = "Add"

      parameter {
        as_path = ["22334"]
      }
    }

    match_criterion {
      match_condition = "Contains"
      route_prefix    = ["10.0.0.0/8"]
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Route Map. Changing this forces a new resource to be created.

* `virtual_hub_id` - (Required) The resource ID of the Virtual Hub. Changing this forces a new resource to be created.

* `rule` - (Optional) A `rule` block as defined below.

---

A `rule` block supports the following:

* `name` - (Required) The unique name for the rule.

* `action` - (Optional) An `action` block as defined below.

* `match_criterion` - (Optional) A `match_criterion` block as defined below.

* `next_step_if_matched` - (Optional) The next step after the rule is evaluated. Possible values are `Continue`, `Terminate` and `Unknown`. Defaults to `Unknown`.

---

An `action` block supports the following:

* `parameter` - A `parameter` block as defined below. Required if `type` is anything other than `Drop`.

* `type` - (Required) The type of the action to be taken. Possible values are `Add`, `Drop`, `Remove`, `Replace` and `Unknown`.

---

A `parameter` block supports the following:

* `as_path` - (Optional) A list of AS paths.

* `community` - (Optional) A list of BGP communities.

* `route_prefix` - (Optional) A list of route prefixes.

---

A `match_criterion` block supports the following:

* `match_condition` - (Required) The match condition to apply the rule of the Route Map. Possible values are `Contains`, `Equals`, `NotContains`, `NotEquals` and `Unknown`.

* `as_path` - (Optional) A list of AS paths which this criterion matches.

* `community` - (Optional) A list of BGP communities which this criterion matches.

* `route_prefix` - (Optional) A list of route prefixes which this criterion matches.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Route Map.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Route Map.
* `read` - (Defaults to 5 minutes) Used when retrieving the Route Map.
* `update` - (Defaults to 1 hour) Used when updating the Route Map.
* `delete` - (Defaults to 1 hour) Used when deleting the Route Map.

## Import

Route Maps can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_route_map.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/virtualHubs/virtualHub1/routeMaps/routeMap1
```
