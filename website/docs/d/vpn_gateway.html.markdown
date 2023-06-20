---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_vpn_gateway"
description: |-
    Manages a VPN Gateway within a Virtual Hub.
---

# Data Source: azurerm_vpn_gateway

Use this data source to access information about an existing VPN Gateway within a Virtual Hub.

## Example Usage

```hcl
data "azurerm_vpn_gateway" "example" {
  name                = "existing-local-vpn_gateway"
  resource_group_name = "existing-vpn_gateway"
}

output "azurerm_vpn_gateway_id" {
  value = data.azurerm_vpn_gateway.example.id
}
```

## Argument Reference

* `name` - (Required) The Name of the VPN Gateway.

* `resource_group_name` - (Required) The name of the Resource Group where the VPN Gateway exists.

## Attributes Reference

* `id` - The ID of the VPN Gateway.

* `location` - The Azure location where the VPN Gateway exists.

* `virtual_hub_id` -  The ID of the Virtual Hub within which this VPN Gateway has been created.

* `bgp_settings` - A `bgp_settings` block as defined below.

* `scale_unit` -  The Scale Unit of this VPN Gateway.

* `tags` - A mapping of tags assigned to the VPN Gateway.

---

A `bgp_settings` block exports the following:

* `asn` - The ASN of the BGP Speaker.

* `peer_weight` -  The weight added to Routes learned from this BGP Speaker.

* `instance_0_bgp_peering_address` -  An `instance_bgp_peering_address` block as defined below.

* `instance_1_bgp_peering_address` -  An `instance_bgp_peering_address` block as defined below.

---

A `instance_bgp_peering_address` block exports the following:

* `custom_ips` -  A list of custom BGP peering addresses to assigned to this instance.

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

* `read` - (Defaults to 5 minutes) Used when retrieving the VPN Gateway.
