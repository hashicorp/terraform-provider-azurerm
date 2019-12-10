---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_express_route_circuit"
sidebar_current: "docs-azurerm-datasource-express-route-circuit"
description: |-
  Gets information about an existing ExpressRoute circuit.
---

# Data Source: azurerm_express_route_circuit

Use this data source to access information about an existing ExpressRoute circuit.

## Example Usage

```hcl
data "azurerm_express_route_circuit" "example" {
  resource_group_name = "${azurerm_resource_group.example.name}"
  name                = "${azurerm_express_route_circuit.example.name}"
}

output "express_route_circuit_id" {
  value = "${data.azurerm_express_route_circuit.example.id}"
}

output "service_key" {
  value = "${data.azurerm_express_route_circuit.example.service_key}"
}
```

## Argument Reference

* `name` - (Required) The name of the ExpressRoute circuit.
* `resource_group_name` - (Required) The Name of the Resource Group where the ExpressRoute circuit exists.

## Attributes Reference

* `id` - The ID of the ExpressRoute circuit.

* `location` - The Azure location where the ExpressRoute circuit exists

* `peerings` - A `peerings` block for the ExpressRoute circuit as documented below

* `service_provider_provisioning_state` - The ExpressRoute circuit provisioning state from your chosen service provider. Possible values are "NotProvisioned", "Provisioning", "Provisioned", and "Deprovisioning".

* `service_key` - The string needed by the service provider to provision the ExpressRoute circuit.

* `service_provider_properties` - A `service_provider_properties` block for the ExpressRoute circuit as documented below

* `sku` - A `sku` block for the ExpressRoute circuit as documented below.

---

`service_provider_properties` supports the following:

* `service_provider_name` - The name of the ExpressRoute Service Provider.
* `peering_location` - The name of the peering location and **not** the Azure resource location.
* `bandwidth_in_mbps` - The bandwidth in Mbps of the ExpressRoute circuit.

`peerings` supports the following:

* `peering_type` - The type of the ExpressRoute Circuit Peering. Acceptable values include `AzurePrivatePeering`, `AzurePublicPeering` and `MicrosoftPeering`. Changing this forces a new resource to be created.
~> **NOTE:** only one Peering of each Type can be created per ExpressRoute circuit.
* `primary_peer_address_prefix` - A `/30` subnet for the primary link.
* `secondary_peer_address_prefix` - A `/30` subnet for the secondary link.
* `vlan_id` - A valid VLAN ID to establish this peering on.
* `shared_key` - The shared key. Can be a maximum of 25 characters.
* `azure_asn` - The Either a 16-bit or a 32-bit ASN for Azure.
* `peer_asn` - The Either a 16-bit or a 32-bit ASN. Can either be public or private.

`sku` supports the following:

* `tier` - The service tier. Possible values are `Standard` or `Premium`.
* `family` - The billing mode for bandwidth. Possible values are `MeteredData` or `UnlimitedData`.
