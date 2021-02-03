---
subcategory: "Databox Edge"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_databox_edge_device"
description: |-
  Gets information about an existing Databox Edge Device.
---

# Data Source: azurerm_databox_edge_device

Use this data source to access information about an existing Databox Edge Device.

## Example Usage

```hcl
data "azurerm_databox_edge_device" "example" {
  name                = "example-device"
  resource_group_name = "example-rg"
}

output "id" {
  value = data.azurerm_databox_edge_device.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Databox Edge Device.

* `resource_group_name` - (Required) The name of the Resource Group where the Databox Edge Device exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Databox Edge Device.

* `location` - The Azure Region where the Databox Edge Device exists.

* `tags` - A mapping of tags assigned to the Databox Edge Device.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Databox Edge Device.
