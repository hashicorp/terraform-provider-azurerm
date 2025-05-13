---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_express_route_circuit"
description: |-
  Manages an ExpressRoute circuit.
---

# azurerm_express_route_circuit

Manages an ExpressRoute circuit.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "exprtTest"
  location = "West Europe"
}

resource "azurerm_express_route_circuit" "example" {
  name                  = "expressRoute1"
  resource_group_name   = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ExpressRoute circuit. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the ExpressRoute circuit. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block for the ExpressRoute circuit as documented below.

* `service_provider_name` - (Optional) The name of the ExpressRoute Service Provider. Changing this forces a new resource to be created.

* `peering_location` - (Optional) The name of the peering location and **not** the Azure resource location. Changing this forces a new resource to be created.

* `bandwidth_in_mbps` - (Optional) The bandwidth in Mbps of the circuit being created on the Service Provider.

~> **Note:** Once you increase your bandwidth, you will not be able to decrease it to its previous value.

~> **Note:** The `service_provider_name`, the `peering_location` and the `bandwidth_in_mbps` should be set together and they conflict with `express_route_port_id` and `bandwidth_in_gbps`.

* `allow_classic_operations` - (Optional) Allow the circuit to interact with classic (RDFE) resources. Defaults to `false`.

* `express_route_port_id` - (Optional) The ID of the Express Route Port this Express Route Circuit is based on. Changing this forces a new resource to be created.

* `bandwidth_in_gbps` - (Optional) The bandwidth in Gbps of the circuit being created on the Express Route Port.

~> **Note:** The `express_route_port_id` and the `bandwidth_in_gbps` should be set together and they conflict with `service_provider_name`, `peering_location` and `bandwidth_in_mbps`.

* `authorization_key` - (Optional) The authorization key. This can be used to set up an ExpressRoute Circuit with an ExpressRoute Port from another subscription.

* `rate_limiting_enabled` - (Optional) Enable [rate limiting](https://learn.microsoft.com/en-us/azure/expressroute/rate-limit) for the circuit. Only works with ExpressRoute Ports. Defaults to `false`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `sku` block supports the following:

* `tier` - (Required) The service tier. Possible values are `Basic`, `Local`, `Standard` or `Premium`.

* `family` - (Required) The billing mode for bandwidth. Possible values are `MeteredData` or `UnlimitedData`.

~> **Note:** You can migrate from `MeteredData` to `UnlimitedData`, but not the other way around.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the ExpressRoute circuit.
* `service_provider_provisioning_state` - The ExpressRoute circuit provisioning state from your chosen service provider. Possible values are `NotProvisioned`, `Provisioning`, `Provisioned`, and `Deprovisioning`.
* `service_key` - The string needed by the service provider to provision the ExpressRoute circuit.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the ExpressRoute Circuit.
* `read` - (Defaults to 5 minutes) Used when retrieving the ExpressRoute Circuit.
* `update` - (Defaults to 30 minutes) Used when updating the ExpressRoute Circuit.
* `delete` - (Defaults to 30 minutes) Used when deleting the ExpressRoute Circuit.

## Import

ExpressRoute circuits can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_express_route_circuit.myExpressRoute /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/expressRouteCircuits/myExpressRoute
```
