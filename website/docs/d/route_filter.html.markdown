---
subcategory: "TODO - pick from: Load Balancer|Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_route_filter"
description: |-
  Gets information about an existing Route Filter.
---

# Data Source: azurerm_route_filter

Use this data source to access information about an existing Route Filter.

## Example Usage

```hcl
data "azurerm_route_filter" "example" {
  name = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_route_filter.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The Name of this Route Filter.

* `resource_group_name` - (Required) The name of the Resource Group where the Route Filter exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Route Filter.

* `location` - The Azure Region where the Route Filter exists.

* `tags` - A mapping of tags assigned to the Route Filter.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Route Filter.