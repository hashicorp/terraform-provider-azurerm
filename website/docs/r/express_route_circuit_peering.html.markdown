---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_express_route_circuit_peering"
description: |-
  Manages an ExpressRoute Circuit Peering.
---

# azurerm_express_route_circuit_peering

Manages an ExpressRoute Circuit Peering.

## Example Usage (Creating a Microsoft Peering)

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

  allow_classic_operations = false

  tags = {
    environment = "Production"
  }
}

resource "azurerm_express_route_circuit_peering" "example" {
  peering_type                  = "MicrosoftPeering"
  express_route_circuit_name    = azurerm_express_route_circuit.example.name
  resource_group_name           = azurerm_resource_group.example.name
  peer_asn                      = 100
  primary_peer_address_prefix   = "123.0.0.0/30"
  secondary_peer_address_prefix = "123.0.0.4/30"
  ipv4_enabled                  = true
  vlan_id                       = 300

  microsoft_peering_config {
    advertised_public_prefixes = ["123.1.0.0/24"]
  }

  ipv6 {
    primary_peer_address_prefix   = "2002:db01::/126"
    secondary_peer_address_prefix = "2003:db01::/126"
    enabled                       = true

    microsoft_peering {
      advertised_public_prefixes = ["2002:db01::/126"]
    }
  }
}
```

## Example Usage (Creating Azure Private Peering)

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

  allow_classic_operations = false

  tags = {
    environment = "Production"
  }
}

resource "azurerm_express_route_circuit_peering" "example" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = azurerm_express_route_circuit.example.name
  resource_group_name           = azurerm_resource_group.example.name
  peer_asn                      = 100
  primary_peer_address_prefix   = "123.0.0.0/30"
  secondary_peer_address_prefix = "123.0.0.4/30"
  ipv4_enabled                  = true
  vlan_id                       = 300

  ipv6 {
    primary_peer_address_prefix   = "2002:db01::/126"
    secondary_peer_address_prefix = "2003:db01::/126"
    enabled                       = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `peering_type` - (Required) The type of the ExpressRoute Circuit Peering. Acceptable values include `AzurePrivatePeering`, `AzurePublicPeering` and `MicrosoftPeering`. 

~> **Note:** only one Peering of each Type can be created. Attempting to create multiple peerings of the same type will overwrite the original peering.

* `express_route_circuit_name` - (Required) The name of the ExpressRoute Circuit in which to create the Peering. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Express Route Circuit Peering. Changing this forces a new resource to be created.

* `vlan_id` - (Required) A valid VLAN ID to establish this peering on.

* `primary_peer_address_prefix` - (Optional) A `/30` subnet for the primary link. Required when config for IPv4.

* `secondary_peer_address_prefix` - (Optional) A `/30` subnet for the secondary link. Required when config for IPv4.

* `ipv4_enabled` - (Optional) A boolean value indicating whether the IPv4 peering is enabled. Defaults to `true`.

* `shared_key` - (Optional) The shared key. Can be a maximum of 25 characters.

* `peer_asn` - (Optional) The Either a 16-bit or a 32-bit ASN. Can either be public or private.

* `microsoft_peering_config` - (Optional) A `microsoft_peering_config` block as defined below. Required when `peering_type` is set to `MicrosoftPeering` and config for IPv4.

* `ipv6` - (Optional) A `ipv6` block as defined below.

* `route_filter_id` - (Optional) The ID of the Route Filter. Only available when `peering_type` is set to `MicrosoftPeering`.

---

A `microsoft_peering_config` block contains:

* `advertised_public_prefixes` - (Required) A list of Advertised Public Prefixes.

* `customer_asn` - (Optional) The CustomerASN of the peering. Defaults to `0`.

* `routing_registry_name` - (Optional) The Routing Registry against which the AS number and prefixes are registered. For example: `ARIN`, `RIPE`, `AFRINIC` etc. Defaults to `NONE`.

* `advertised_communities` - (Optional) The communities of Bgp Peering specified for microsoft peering.

---

A `ipv6` block contains:

* `primary_peer_address_prefix` - (Required) A subnet for the primary link.

* `secondary_peer_address_prefix` - (Required) A subnet for the secondary link.

* `enabled` - (Optional) A boolean value indicating whether the IPv6 peering is enabled. Defaults to `true`.

* `microsoft_peering` - (Optional) A `microsoft_peering` block as defined below. 

* `route_filter_id` - (Optional) The ID of the Route Filter. Only available when `peering_type` is set to `MicrosoftPeering`.

~> **Note:** `ipv6` can be specified when `peering_type` is `MicrosoftPeering` or `AzurePrivatePeering`

---

A `microsoft_peering` block contains:

* `advertised_public_prefixes` - (Optional) A list of Advertised Public Prefixes.

* `customer_asn` - (Optional) The CustomerASN of the peering. Defaults to `0`.

* `routing_registry_name` - (Optional) The Routing Registry against which the AS number and prefixes are registered. For example: `ARIN`, `RIPE`, `AFRINIC` etc. Defaults to `NONE`.

* `advertised_communities` - (Optional) The communities of Bgp Peering specified for microsoft peering.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the ExpressRoute Circuit Peering.

* `azure_asn` - The ASN used by Azure.

* `primary_azure_port` - The Primary Port used by Azure for this Peering.

* `secondary_azure_port` - The Secondary Port used by Azure for this Peering.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the ExpressRoute Circuit Peering.

* `read` - (Defaults to 5 minutes) Used when retrieving the ExpressRoute Circuit Peering.

* `update` - (Defaults to 30 minutes) Used when updating the ExpressRoute Circuit Peering.

* `delete` - (Defaults to 30 minutes) Used when deleting the ExpressRoute Circuit Peering.

## Import

ExpressRoute Circuit Peerings can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_express_route_circuit_peering.peering1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/expressRouteCircuits/myExpressRoute/peerings/peering1
```
