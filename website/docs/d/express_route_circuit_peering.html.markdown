---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_express_route_circuit_peering"
description: |-
  Gets information about an existing ExpressRoute Circuit Peering.
---

# Data Source: azurerm_express_route_circuit_peering

Use this data source to access information about an existing ExpressRoute Circuit Peering.

## Example Usage

```hcl
data "azurerm_express_route_circuit_peering" "example" {
  peering_type               = "example-peering"
  express_route_circuit_name = "example-expressroute"
  resource_group_name        = "example-resources"
}
```

## Argument Reference

The following arguments are supported:

* `peering_type` - (Required) The type of the ExpressRoute Circuit Peering. Acceptable values include `AzurePrivatePeering`, `AzurePublicPeering` and `MicrosoftPeering`.

* `express_route_circuit_name` - (Required) The name of the ExpressRoute Circuit in which to create the Peering. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Express Route Circuit Peering. Changing this forces a new resource to be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `peering_type` - The type of the ExpressRoute Circuit Peering.

* `vlan_id` - The VLAN ID used for this peering.

* `primary_peer_address_prefix` - The primary peer address prefix.

* `secondary_peer_address_prefix` - The secondary peer address prefix.

* `ipv4_enabled` - Indicates if IPv4 is enabled.

* `ipv6` - A block describing the IPv6 configuration, if enabled.

* `microsoft_peering_config` - A block describing the Microsoft Peering configuration, if applicable.

* `azure_asn` - The ASN used by Azure for the peering.

* `primary_azure_port` - The primary port used by Azure for this peering.

* `secondary_azure_port` - The secondary port used by Azure for this peering.

* `id` - The ID of the ExpressRoute Circuit Peering.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the ExpressRoute Circuit Peering.
