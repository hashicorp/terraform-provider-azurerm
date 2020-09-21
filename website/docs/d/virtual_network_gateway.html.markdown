---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_network_gateway"
description: |-
  Gets information about an existing Virtual Network Gateway.
---

# Data Source: azurerm_virtual_network_gateway

Use this data source to access information about an existing Virtual Network Gateway.

## Example Usage

```hcl
data "azurerm_virtual_network_gateway" "example" {
  name                = "production"
  resource_group_name = "networking"
}

output "virtual_network_gateway_id" {
  value = data.azurerm_virtual_network_gateway.example.id
}
```

## Argument Reference

* `name` - Specifies the name of the Virtual Network Gateway.
* `resource_group_name` - Specifies the name of the resource group the Virtual Network Gateway is located in.

## Attributes Reference

* `id` - The ID of the Virtual Network Gateway.

* `location` - The location/region where the Virtual Network Gateway is located.

* `type` - The type of the Virtual Network Gateway.

* `vpn_type` - The routing type of the Virtual Network Gateway.

* `enable_bgp` - Will BGP (Border Gateway Protocol) will be enabled
    for this Virtual Network Gateway.

* `active_active` - Is this an Active-Active Gateway?

* `default_local_network_gateway_id` -  The ID of the local network gateway
    through which outbound Internet traffic from the virtual network in which the
    gateway is created will be routed (*forced tunneling*). Refer to the
    [Azure documentation on forced tunneling](https://docs.microsoft.com/en-us/azure/vpn-gateway/vpn-gateway-forced-tunneling-rm).

* `sku` - Configuration of the size and capacity of the Virtual Network Gateway.

* `generation` - The Generation of the Virtual Network Gateway.

* `ip_configuration` - One or two `ip_configuration` blocks documented below.

* `vpn_client_configuration` - A `vpn_client_configuration` block which is documented below.

* `tags` - A mapping of tags assigned to the resource.

The `ip_configuration` block supports:

* `name` - A user-defined name of the IP configuration.

* `private_ip_address_allocation` - Defines how the private IP address
    of the gateways virtual interface is assigned.

* `subnet_id` - The ID of the gateway subnet of a virtual network in
    which the virtual network gateway will be created. It is mandatory that
    the associated subnet is named `GatewaySubnet`. Therefore, each virtual
    network can contain at most a single Virtual Network Gateway.

* `public_ip_address_id` - The ID of the Public IP Address associated
    with the Virtual Network Gateway.

The `vpn_client_configuration` block supports:

* `address_space` - The address space out of which ip addresses for
    vpn clients will be taken. You can provide more than one address space, e.g.
    in CIDR notation.

* `aad_tenant` - AzureAD Tenant URL
    This setting is incompatible with the use of
    `root_certificate` and `revoked_certificate`, `radius_server_address`, and `radius_server_secret`.

* `aad_audience` - The client id of the Azure VPN application.
    See [Create an Active Directory (AD) tenant for P2S OpenVPN protocol connections](https://docs.microsoft.com/en-gb/azure/vpn-gateway/openvpn-azure-ad-tenant-multi-app) for values
    This setting is incompatible with the use of
    `root_certificate` and `revoked_certificate`, `radius_server_address`, and `radius_server_secret`.

* `aad_issuer` - The STS url for your tenant
    This setting is incompatible with the use of
    `root_certificate` and `revoked_certificate`, `radius_server_address`, and `radius_server_secret`.

* `root_certificate` - One or more `root_certificate` blocks which are
    defined below. These root certificates are used to sign the client certificate
    used by the VPN clients to connect to the gateway.
    This setting is incompatible with the use of
    `aad_tenant`, `aad_audience`, `aad_issuer`, `radius_server_address`, and `radius_server_secret`.

* `revoked_certificate` - One or more `revoked_certificate` blocks which
    are defined below.
    This setting is incompatible with the use of
    `aad_tenant`, `aad_audience`, `aad_issuer`, `radius_server_address`, and `radius_server_secret`.

* `radius_server_address` - The address of the Radius server.
    This setting is incompatible with the use of
    `aad_tenant`, `aad_audience`, `aad_issuer`, `root_certificate` and `revoked_certificate`.

* `radius_server_secret` - The secret used by the Radius server.
    This setting is incompatible with the use of
    `aad_tenant`, `aad_audience`, `aad_issuer`, `root_certificate` and `revoked_certificate`.

* `vpn_client_protocols` - List of the protocols supported by the vpn client.
    The supported values are `SSTP`, `IkeV2` and `OpenVPN`.

The `bgp_settings` block supports:

* `asn` - The Autonomous System Number (ASN) to use as part of the BGP.

* `peering_address` - The BGP peer IP address of the virtual network
    gateway. This address is needed to configure the created gateway as a BGP Peer
    on the on-premises VPN devices.

* `peer_weight` - The weight added to routes which have been learned
    through BGP peering.

The `root_certificate` block supports:

* `name` - The user-defined name of the root certificate.

* `public_cert_data` - The public certificate of the root certificate
    authority. The certificate must be provided in Base-64 encoded X.509 format
    (PEM).

The `root_revoked_certificate` block supports:

* `name` - The user-defined name of the revoked certificate.

* `public_cert_data` - The SHA1 thumbprint of the certificate to be revoked.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Network Gateway.
