---
subcategory: "Databox Edge"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_databox_edge_order"
description: |-
  Gets information about an existing Databox Edge Order.
---

# Data Source: azurerm_databox_edge_order

Use this data source to access information about an existing Databox Edge Order.

## Example Usage

```hcl
data "azurerm_databox_edge_order" "example" {
  resource_group_name = "existing"
  device_name = "existing"
}

output "id" {
  value = data.azurerm_databox_edge_order.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the Resource Group where the Databox Edge Order exists.

* `device_name` - (Required) The device name.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Databox Edge Order.

* `name` - The Name of this Databox Edge Order.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Databox Edge Order.