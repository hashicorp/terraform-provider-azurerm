---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_vpn_gateway"
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

* `bgp_route_translation_for_nat_enabled` - (Optional) Is BGP route translation for NAT on this VPN Gateway enabled? Defaults to `false`.

* `bgp_settings` - (Optional) A `bgp_settings` block as defined below.

* `routing_preference` - (Optional) Azure routing preference lets you to choose how your traffic routes between Azure and the internet. You can choose to route traffic either via the Microsoft network (default value, `Microsoft Network`), or via the ISP network (public internet, set to `Internet`). More context of the configuration can be found in the [Microsoft Docs](https://docs.microsoft.com/azure/virtual-wan/virtual-wan-site-to-site-portal#gateway) to create a VPN Gateway. Defaults to `Microsoft Network`. Changing this forces a new resource to be created.

* `scale_unit` - (Optional) The Scale Unit for this VPN Gateway. Defaults to `1`.

* `tags` - (Optional) A mapping of tags to assign to the VPN Gateway.

---

A `bgp_settings` block supports the following:

* `asn` - (Required) The ASN of the BGP Speaker. Changing this forces a new resource to be created.

* `peer_weight` - (Required) The weight added to Routes learned from this BGP Speaker. Changing this forces a new resource to be created.

* `instance_0_bgp_peering_address` - (Optional) An `instance_bgp_peering_address` block as defined below.

* `instance_1_bgp_peering_address` - (Optional) An `instance_bgp_peering_address` block as defined below.

---

A `instance_bgp_peering_address` block supports the following:

* `custom_ips` - (Required) A list of custom BGP peering addresses to assign to this instance.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the VPN Gateway.

* `bgp_settings` - A `bgp_settings` block as defined below.

---

A `bgp_settings` block exports the following:

* `bgp_peering_address` - The Address which should be used for the BGP Peering.

* `instance_0_bgp_peering_address` - an `instance_bgp_peering_address` block as defined below.

* `instance_1_bgp_peering_address` - an `instance_bgp_peering_address` block as defined below.

---

A `instance_bgp_peering_address` block exports the following:

* `ip_configuration_id` - The pre-defined id of VPN Gateway IP Configuration.

* `default_ips` - The list of default BGP peering addresses which belong to the pre-defined VPN Gateway IP configuration.

* `tunnel_ips` - The list of tunnel public IP addresses which belong to the pre-defined VPN Gateway IP configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the VPN Gateway.
* `read` - (Defaults to 5 minutes) Used when retrieving the VPN Gateway.
* `update` - (Defaults to 90 minutes) Used when updating the VPN Gateway.
* `delete` - (Defaults to 90 minutes) Used when deleting the VPN Gateway.

## Import

VPN Gateways can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_vpn_gateway.gateway1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/vpnGateways/gateway1
```
