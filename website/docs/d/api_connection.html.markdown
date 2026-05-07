---
subcategory: "Connections"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_connection"
description: |-
  Gets information about an existing API Connection.
---

# Data Source: azurerm_api_connection

Use this data source to access information about an existing API Connection.

## Example Usage

```hcl
data "azurerm_api_connection" "example" {
  name                = "example-connection"
  resource_group_name = "example-resources"
}

output "connection_id" {
  value = data.azurerm_api_connection.example.id
}

output "connection_name" {
  value = data.azurerm_api_connection.example.name
}

output "managed_api_id" {
  value = data.azurerm_api_connection.example.managed_api_id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the API Connection.

* `resource_group_name` - (Required) The name of the Resource Group where the API Connection exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Connection.

* `location` - The Azure Region where the API Connection exists.

* `managed_api_id` - The ID of the Managed API that this connection is linked to.

* `display_name` - The display name of the API Connection.

* `parameter_values` - A mapping of parameter names to their values for the API Connection.

* `tags` - A mapping of tags assigned to the API Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the API Connection.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Web` - 2016-06-01
