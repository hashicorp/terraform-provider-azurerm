---
subcategory: "Network"
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
  name                = "example"
  resource_group_name = "example"
  location            = "East US"

  rule {
    name        = "rule"
    access      = "Allow"
    rule_type   = "Community"
    communities = ["12076:52004"]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Route Filter should exist. Changing this forces a new Route Filter to be created.

* `name` - (Required) The Name which should be used for this Route Filter.

* `resource_group_name` - (Required) The name of the Resource Group where the Route Filter should exist. Changing this forces a new Route Filter to be created.

---

* `rule` - (Optional) A `rules` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Route Filter.

---

A `rule` block supports the following:

* `access` - (Required) The access type of the rule. The only possible value is `Allow`.

* `communities` - (Required) The collection for bgp community values to filter on. e.g. ['12076:5010','12076:5020'].

* `name` - (Required) The name of the route filter rule.

* `rule_type` - (Required) The rule type of the rule. The only possible value is `Community`.

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
