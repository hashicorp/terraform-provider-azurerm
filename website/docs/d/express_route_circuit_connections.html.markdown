---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_express_route_circuit_connection"
description: |-
  Gets information about an existing network ExpressRouteCircuitConnection.
---

# Data Source: azurerm_express_route_circuit_connection

Use this data source to access information about an existing network ExpressRouteCircuitConnection.

## Example Usage

```hcl
data "azurerm_express_route_circuit_connection" "example" {
  name = "example-expressroutecircuitconnection"
  resource_group_name = "existing"
  circuit_name = "existing"
}

output "id" {
  value = data.azurerm_express_route_circuit_connection.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this network ExpressRouteCircuitConnection.

* `resource_group_name` - (Required) The name of the Resource Group where the network ExpressRouteCircuitConnection exists.

* `circuit_name` - (Required) The name of the express route circuit.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the network ExpressRouteCircuitConnection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the network ExpressRouteCircuitConnection.