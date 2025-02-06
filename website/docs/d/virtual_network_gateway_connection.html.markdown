---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_network_gateway_connection"
description: |-
  Gets information about an existing Virtual Network Gateway Connection.
---

# Data Source: azurerm_virtual_network_gateway_connection

Use this data source to access information about an existing Virtual Network Gateway Connection.

## Example Usage

```hcl
data "azurerm_virtual_network_gateway_connection" "example" {
  name                = "production"
  resource_group_name = "networking"
}

output "virtual_network_gateway_connection_id" {
  value = data.azurerm_virtual_network_gateway_connection.example.id
}
```

## Argument Reference

* `name` - Specifies the name of the Virtual Network Gateway Connection.
* `resource_group_name` - Specifies the name of the resource group the Virtual Network Gateway Connection is located in.

## Attributes Reference

* `id` - The ID of the Virtual Network Gateway Connection.

* `location` - The location/region where the connection is
    located.

* `type` - The type of connection. Valid options are `IPsec`
    (Site-to-Site), `ExpressRoute` (ExpressRoute), and `Vnet2Vnet` (VNet-to-VNet).

* `virtual_network_gateway_id` - The ID of the Virtual Network Gateway
    in which the connection is created.

* `authorization_key` - The authorization key associated with the
    Express Route Circuit. This field is present only if the type is an
    ExpressRoute connection.

* `dpd_timeout_seconds` - The dead peer detection timeout of this connection in seconds.

* `express_route_circuit_id` - The ID of the Express Route Circuit
    (i.e. when `type` is `ExpressRoute`).

* `peer_virtual_network_gateway_id` - The ID of the peer virtual
    network gateway when a VNet-to-VNet connection (i.e. when `type`
    is `Vnet2Vnet`).

* `local_azure_ip_address_enabled` - Use private local Azure IP for the connection.

* `local_network_gateway_id` - The ID of the local network gateway
    when a Site-to-Site connection (i.e. when `type` is `IPsec`).

* `routing_weight` - The routing weight.

* `shared_key` - The shared IPSec key.

* `enable_bgp` - If `true`, BGP (Border Gateway Protocol) is enabled
    for this connection.

* `custom_bgp_addresses` - (Optional) A `custom_bgp_addresses` (Border Gateway Protocol custom IP Addresses) block which is documented below.
    The block can only be used on `IPSec` / `activeactive` connections,
    For details about see [the relevant section in the Azure documentation](https://docs.microsoft.com/en-us/azure/vpn-gateway/vpn-gateway-howto-aws-bgp).

* `express_route_gateway_bypass` - If `true`, data packets will bypass ExpressRoute Gateway for data forwarding. This is only valid for ExpressRoute connections.

* `private_link_fast_path_enabled` - If `true`, data packets will bypass the Express Route gateway when accessing private-links.
   This is only valid for ExpressRoute connections, on the conditions described in [the relevant section in the Azure documentation](https://learn.microsoft.com/en-us/azure/expressroute/expressroute-howto-linkvnet-arm#fastpath-virtual-network-peering-user-defined-routes-udrs-and-private-link-support-for-expressroute-direct-connections)

* `use_policy_based_traffic_selectors` - If `true`, policy-based traffic
    selectors are enabled for this connection. Enabling policy-based traffic
    selectors requires an `ipsec_policy` block.

* `ipsec_policy` (Optional) A `ipsec_policy` block which is documented below.
    Only a single policy can be defined for a connection. For details on
    custom policies refer to [the relevant section in the Azure documentation](https://docs.microsoft.com/azure/vpn-gateway/vpn-gateway-ipsecikepolicy-rm-powershell).

* `traffic_selector_policy` One or more `traffic_selector_policy` blocks which are documented below.
  A `traffic_selector_policy` allows to specify a traffic selector policy proposal to be used in a virtual network gateway connection.
  For details about traffic selectors refer to [the relevant section in the Azure documentation](https://docs.microsoft.com/azure/vpn-gateway/vpn-gateway-connect-multiple-policybased-rm-ps).

* `tags` - A mapping of tags to assign to the resource.

The `custom_bgp_addresses` block supports:

* `primary` (Required) single IP address that is part of the `azurerm_virtual_network_gateway` ip_configuration (first one)
* `secondary` (Required) single IP address that is part of the `azurerm_virtual_network_gateway` ip_configuration (second one)

The `ipsec_policy` block supports:

* `dh_group` - The DH group used in IKE phase 1 for initial SA. Valid
    options are `DHGroup1`, `DHGroup14`, `DHGroup2`, `DHGroup2048`, `DHGroup24`,
    `ECP256`, `ECP384`, or `None`.

* `ike_encryption` - The IKE encryption algorithm. Valid
    options are `AES128`, `AES192`, `AES256`, `DES`, or `DES3`.

* `ike_integrity` - The IKE integrity algorithm. Valid
    options are `MD5`, `SHA1`, `SHA256`, or `SHA384`.

* `ipsec_encryption` - The IPSec encryption algorithm. Valid
    options are `AES128`, `AES192`, `AES256`, `DES`, `DES3`, `GCMAES128`, `GCMAES192`, `GCMAES256`, or `None`.

* `ipsec_integrity` - The IPSec integrity algorithm. Valid
    options are `GCMAES128`, `GCMAES192`, `GCMAES256`, `MD5`, `SHA1`, or `SHA256`.

* `pfs_group` - The DH group used in IKE phase 2 for new child SA.
    Valid options are `ECP256`, `ECP384`, `PFS1`, `PFS2`, `PFS2048`, `PFS24`,
    or `None`.

* `sa_datasize` - The IPSec SA payload size in KB. Must be at least
    `1024` KB.

* `sa_lifetime` - The IPSec SA lifetime in seconds. Must be at least
    `300` seconds.

The `traffic_selector_policy` block supports:

* `local_address_cidrs` - List of local CIDRs.

* `remote_address_cidrs` - List of remote CIDRs.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Network Gateway Connection.
