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
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_managed_api" "example" {
  name     = "servicebus"
  location = azurerm_resource_group.example.location
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "acctestsbn-conn-example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Basic"
}

resource "azurerm_api_connection" "example" {
  name                = "example-connection"
  resource_group_name = azurerm_resource_group.example.name
  managed_api_id      = data.azurerm_managed_api.example.id
  display_name        = "Example 1"

  parameter_values = {
    connectionString = azurerm_servicebus_namespace.example.default_primary_connection_string
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

## Arguments Reference

The following arguments are supported:

* `managed_api_id` - (Required) The ID of the Managed API which this API Connection is linked to. Changing this forces a new API Connection to be created.

* `name` - (Required) The Name which should be used for this API Connection. Changing this forces a new API Connection to be created.

* `resource_group_name` - (Required) The name of the Resource Group where this API Connection should exist. Changing this forces a new API Connection to be created.

---

* `display_name` - (Optional) A display name for this API Connection. Defaults to `Service Bus`. Changing this forces a new API Connection to be created.

* `parameter_values` - (Optional) A map of parameter values associated with this API Connection. Changing this forces a new API Connection to be created.

-> **Note:** The Azure API doesn't return sensitive parameters in the API response which can lead to a diff, as such you may need to use Terraform's `ignore_changes` functionality on this field as shown in the Example Usage above.

* `tags` - (Optional) A mapping of tags which should be assigned to the API Connection.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Connection.
* `update` - (Defaults to 30 minutes) Used when updating the API Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Connection.

## Import

API Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Web/connections/example-connection
```
