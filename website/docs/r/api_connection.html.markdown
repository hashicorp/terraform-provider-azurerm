---
subcategory: "Connections"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_connection"
description: |-
  Manages an API Connection.
---

# azurerm_api_connection

Manages an API Connection.

## Example Usage

```hcl
resource "azurerm_api_connection" "test" {
  name                = "example-connection"
  resource_group_name = azurerm_resource_group.test.name
  managed_api_id      = data.azurerm_managed_api.test.id
  display_name        = "Example 1"

  parameter_values = {
    connectionString = azurerm_servicebus_namespace.test.default_primary_connection_string
  }

  tags = {
    Hello = "World"
  }

  lifecycle {
    # NOTE: since the connectionString is a secure value it's not returned from the API
    ignore_changes = ["parameter_values"]
  }
}
```