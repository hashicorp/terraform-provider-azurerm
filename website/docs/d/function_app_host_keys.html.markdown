---
subcategory: 'App Service (Web Apps)'
layout: 'azurerm'
page_title: 'Azure Resource Manager: azurerm_function_app_host_keys'
description: |-
  Gets the Host Keys of an existing Function App.
---

# Data Source: azurerm_function_app_host_keys

Use this data source to fetch the Host Keys of an existing Function App

## Example Usage

```hcl
data "azurerm_function_app_host_keys" "example" {
  name                = "example-function"
  resource_group_name = azurerm_resource_group.example.name
}
```

~> **Note:** All arguments including the secret value will be stored in the raw state as plain-text, including `default_function_key` and `primary_key`. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Argument Reference

The following arguments are supported:

- `name` - The name of the Function App.

- `resource_group_name` - The name of the Resource Group where the Function App exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `default_function_key` - Function App resource's default function key.

- `primary_key` - Function App resource's secret key

- `event_grid_extension_config_key` - Function App resource's Event Grid Extension Config system key.

- `signalr_extension_key` - Function App resource's SignalR Extension system key.

- `durabletask_extension_key` - Function App resource's Durable Task Extension system key.

- `webpubsub_extension_key` - Function App resource's Web PubSub Extension system key.

- `blobs_extension_key` - Function App resource's Blobs Extension system key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Function App Host Keys