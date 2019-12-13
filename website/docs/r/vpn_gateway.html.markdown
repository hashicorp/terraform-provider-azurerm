---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_vpn_gateway"
sidebar_current: "docs-azurerm-network-vpn-gateway"
description: |-
    Manages a VPN Gateway within a Virtual Hub.
---

# azurerm_vpn_gateway

Manages a VPN Gateway within a Virtual Hub, which enables Site-to-Site communication.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_virtual_wan" "example" {
  name                = "example-vwan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_hub" "example" {
  name                = "example-hub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_wan_id      = azurerm_virtual_wan.example.id
  address_prefix      = "10.0.1.0/24"
}

resource "azurerm_vpn_gateway" "example" {
  name                = "example-vpng"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  virtual_hub_id      = azurerm_virtual_hub.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Name which should be used for this VPN Gateway. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The Name of the Resource Group in which this VPN Gateway should be created. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where this VPN Gateway should be created. Changing this forces a new resource to be created.

* `virtual_hub_id` - (Required) The ID of the Virtual Hub within which this VPN Gateway should be created. Changing this forces a new resource to be created.

---

* `bgp_settings` - (Optional) A `bgp_settings` block as defined below.

* `scale_unit` - (Optional) The Scale Unit for this VPN Gateway. Defaults to `1`.

* `tags` - (Optional) A mapping of tags to assign to the VPN Gateway.

---

A `bgp_settings` block supports the following:

* `asn` - (Required) The ASN of the BGP Speaker. Changing this forces a new resource to be created.

* `peer_weight` - (Required) The weight added to Routes learned from this BGP Speaker. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the VPN Gateway.

* `bgp_settings` - A `bgp_settings` block as defined below.

---

A `bgp_settings` block exports the following:

* `bgp_peering_address` - The Address which should be used for the BGP Peering.

## Import

VPN Gateways can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_vpn_gateway.gateway1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/vpnGateways/gateway1
```
