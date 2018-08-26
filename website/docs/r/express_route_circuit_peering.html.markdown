---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_express_route_circuit_peering"
sidebar_current: "docs-azurerm-resource-network-express-route-circuit-peering"
description: |-
  Manages an ExpressRoute Circuit Peering.
---

# azurerm_express_route_circuit_peering

Manages an ExpressRoute Circuit Peering.

## Example Usage (Creating a Microsoft Peering)

```hcl
resource "azurerm_resource_group" "test" {
  name     = "exprtTest"
  location = "West US"
}

resource "azurerm_express_route_circuit" "test" {
  name                     = "expressRoute1"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  service_provider_name    = "Equinix"
  peering_location         = "Silicon Valley"
  bandwidth_in_mbps        = 50
  sku {
    tier   = "Standard"
    family = "MeteredData"
  }
  allow_classic_operations = false

  tags {
    environment = "Production"
  }
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "MicrosoftPeering"
  express_route_circuit_name    = "${azurerm_express_route_circuit.test.name}"
  resource_group_name           = "${azurerm_resource_group.test.name}"
  peer_asn                      = 100
  primary_peer_address_prefix   = "123.0.0.0/30"
  secondary_peer_address_prefix = "123.0.0.4/30"
  vlan_id                       = 300

  microsoft_peering_config {
    advertised_public_prefixes = ["123.1.0.0/24"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `peering_type` - (Required) The type of the ExpressRoute Circuit Peering. Acceptable values include `AzurePrivatePeering`, `AzurePublicPeering` and `MicrosoftPeering`. Changing this forces a new resource to be created.

~> **NOTE:** only one Peering of each Type can be created. Attempting to create multiple peerings of the same type will overwrite the original peering.

* `express_route_circuit_name` - (Required) The name of the ExpressRoute Circuit in which to create the Peering.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the Express Route Circuit Peering. Changing this forces a new resource to be created.

* `primary_peer_address_prefix` - (Optional) A `/30` subnet for the primary link.
* `secondary_peer_address_prefix` - (Optional) A `/30` subnet for the secondary link.
* `vlan_id` - (Optional) A valid VLAN ID to establish this peering on.
* `shared_key` - (Optional) The shared key. Can be a maximum of 25 characters.
* `peer_asn` - (Optional) The Either a 16-bit or a 32-bit ASN. Can either be public or private..
* `microsoft_peering_config` - (Optional) A `microsoft_peering_config` block as defined below. Required when `peering_type` is set to `MicrosoftPeering`.

---

A `microsoft_peering_config` block contains:

* `advertised_public_prefixes` - (Required) A list of Advertised Public Prefixes

## Attributes Reference

The following attributes are exported:

* `id` - The Resource ID of the ExpressRoute Circuit Peering.

* `azure_asn` - The ASN used by Azure.

* `primary_azure_port` - The Primary Port used by Azure for this Peering.

* `secondary_azure_port` - The Secondary Port used by Azure for this Peering.

## Import

ExpressRoute Circuit Peerings can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_express_route_circuit_peering.peering1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/expressRouteCircuits/myExpressRoute/peerings/peering1
```
