---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_dev_center_attached_network"
description: |-
  Gets information about an existing Dev Center Attached Network.
---

# Data Source: azurerm_dev_center_attached_network

Use this data source to access information about an existing Dev Center Attached Network.

## Example Usage

```hcl
data "azurerm_dev_center_attached_network" "example" {
  name          = azurerm_dev_center_attached_network.example.name
  dev_center_id = azurerm_dev_center_attached_network.example.dev_center_id
}

output "id" {
  value = data.azurerm_dev_center_attached_network.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Dev Center Attached Network.

* `dev_center_id` - (Required) The ID of the associated Dev Center.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center Attached Network.

* `network_connection_id` - The ID of the attached Dev Center Network Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Center Attached Network.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.DevCenter`: 2025-02-01
