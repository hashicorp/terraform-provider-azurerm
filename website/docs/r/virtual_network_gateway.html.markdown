---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_network_gateway"
description: |-
  Manages a virtual network gateway to establish secure, cross-premises connectivity.
---

# azurerm_virtual_network_gateway

Manages a Virtual Network Gateway to establish secure, cross-premises connectivity.

-> **Note:** Please be aware that provisioning a Virtual Network Gateway takes a long time (between 30 minutes and 1 hour)

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "test"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "test"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "example" {
  name                = "test"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  allocation_method = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "example" {
  name                = "test"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  type     = "Vpn"
  vpn_type = "RouteBased"

  active_active = false
  enable_bgp    = false
  sku           = "Basic"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.example.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.example.id
  }

  vpn_client_configuration {
    address_space = ["10.2.0.0/24"]

    root_certificate {
      name = "DigiCert-Federated-ID-Root-CA"

      public_cert_data = <<EOF
MIIDuzCCAqOgAwIBAgIQCHTZWCM+IlfFIRXIvyKSrjANBgkqhkiG9w0BAQsFADBn
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSYwJAYDVQQDEx1EaWdpQ2VydCBGZWRlcmF0ZWQgSUQg
Um9vdCBDQTAeFw0xMzAxMTUxMjAwMDBaFw0zMzAxMTUxMjAwMDBaMGcxCzAJBgNV
BAYTAlVTMRUwEwYDVQQKEwxEaWdpQ2VydCBJbmMxGTAXBgNVBAsTEHd3dy5kaWdp
Y2VydC5jb20xJjAkBgNVBAMTHURpZ2lDZXJ0IEZlZGVyYXRlZCBJRCBSb290IENB
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvAEB4pcCqnNNOWE6Ur5j
QPUH+1y1F9KdHTRSza6k5iDlXq1kGS1qAkuKtw9JsiNRrjltmFnzMZRBbX8Tlfl8
zAhBmb6dDduDGED01kBsTkgywYPxXVTKec0WxYEEF0oMn4wSYNl0lt2eJAKHXjNf
GTwiibdP8CUR2ghSM2sUTI8Nt1Omfc4SMHhGhYD64uJMbX98THQ/4LMGuYegou+d
GTiahfHtjn7AboSEknwAMJHCh5RlYZZ6B1O4QbKJ+34Q0eKgnI3X6Vc9u0zf6DH8
Dk+4zQDYRRTqTnVO3VT8jzqDlCRuNtq6YvryOWN74/dq8LQhUnXHvFyrsdMaE1X2
DwIDAQABo2MwYTAPBgNVHRMBAf8EBTADAQH/MA4GA1UdDwEB/wQEAwIBhjAdBgNV
HQ4EFgQUGRdkFnbGt1EWjKwbUne+5OaZvRYwHwYDVR0jBBgwFoAUGRdkFnbGt1EW
jKwbUne+5OaZvRYwDQYJKoZIhvcNAQELBQADggEBAHcqsHkrjpESqfuVTRiptJfP
9JbdtWqRTmOf6uJi2c8YVqI6XlKXsD8C1dUUaaHKLUJzvKiazibVuBwMIT84AyqR
QELn3e0BtgEymEygMU569b01ZPxoFSnNXc7qDZBDef8WfqAV/sxkTi8L9BkmFYfL
uGLOhRJOFprPdoDIUBB+tmCl3oDcBy3vnUeOEioz8zAkprcb3GHwHAK+vHmmfgcn
WsfMLH4JCLa/tRYL+Rw/N3ybCkDp00s0WUZ+AoDywSl0Q/ZEnNY0MsFiw6LyIdbq
M/s/1JRtO3bDSzD9TazRVzn2oBqzSa8VgIo5C1nOnoAKJTlsClJKvIhnRlaLQqk=
EOF

    }

    revoked_certificate {
      name       = "Verizon-Global-Root-CA"
      thumbprint = "912198EEF23DCAC40939312FEE97DD560BAE49B1"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `ip_configuration` - (Required) One or more (up to 3) `ip_configuration` blocks documented below. Changing this forces a new resource to be created. 
  An active-standby gateway requires exactly one `ip_configuration` block,
  an active-active gateway requires exactly two `ip_configuration` blocks whereas
  an active-active zone redundant gateway with P2S configuration requires exactly three `ip_configuration` blocks.

* `location` - (Required) The location/region where the Virtual Network Gateway is located. Changing this forces a new resource to be created.

* `name` - (Required) The name of the Virtual Network Gateway. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Virtual Network Gateway. Changing this forces a new resource to be created.

* `sku` - (Required) Configuration of the size and capacity of the virtual network gateway. Valid options are `Basic`, `Standard`, `HighPerformance`, `UltraPerformance`, `ErGw1AZ`, `ErGw2AZ`, `ErGw3AZ`, `VpnGw1`, `VpnGw2`, `VpnGw3`, `VpnGw4`,`VpnGw5`, `VpnGw1AZ`, `VpnGw2AZ`, `VpnGw3AZ`,`VpnGw4AZ` and `VpnGw5AZ` and depend on the `type`, `vpn_type` and `generation` arguments. A `PolicyBased` gateway only supports the `Basic` SKU. Further, the `UltraPerformance` SKU is only supported by an `ExpressRoute` gateway.

~> **Note:** To build a UltraPerformance ExpressRoute Virtual Network gateway, the associated Public IP needs to be SKU "Basic" not "Standard"

~> **Note:** Not all SKUs (e.g. `ErGw1AZ`) are available in all regions. If you see `StatusCode=400 -- Original Error: Code="InvalidGatewaySkuSpecifiedForGatewayDeploymentType"` please try another region.

* `type` - (Required) The type of the Virtual Network Gateway. Valid options are `Vpn` or `ExpressRoute`. Changing the type forces a new resource to be created.

---

* `active_active` - (Optional) If `true`, an active-active Virtual Network Gateway will be created. An active-active gateway requires a `HighPerformance` or an `UltraPerformance` SKU. If `false`, an active-standby gateway will be created. Defaults to `false`.

* `default_local_network_gateway_id` - (Optional) The ID of the local network gateway through which outbound Internet traffic from the virtual network in which the gateway is created will be routed (*forced tunnelling*). Refer to the [Azure documentation on forced tunnelling](https://docs.microsoft.com/azure/vpn-gateway/vpn-gateway-forced-tunneling-rm). If not specified, forced tunnelling is disabled.

* `edge_zone` - (Optional) Specifies the Edge Zone within the Azure Region where this Virtual Network Gateway should exist. Changing this forces a new Virtual Network Gateway to be created.

* `enable_bgp` - (Optional) If `true`, BGP (Border Gateway Protocol) will be enabled for this Virtual Network Gateway. Defaults to `false`.

* `bgp_settings` - (Optional) A `bgp_settings` block which is documented below. In this block the BGP specific settings can be defined.

* `custom_route` - (Optional) A `custom_route` block as defined below. Specifies a custom routes address space for a virtual network gateway and a VpnClient.

* `generation` - (Optional) The Generation of the Virtual Network gateway. Possible values include `Generation1`, `Generation2` or `None`. Changing this forces a new resource to be created.

-> **Note:** The available values depend on the `type` and `sku` arguments - where `Generation2` is only value for a `sku` larger than `VpnGw2` or `VpnGw2AZ`.

* `private_ip_address_enabled` - (Optional) Should private IP be enabled on this gateway for connections? Changing this forces a new resource to be created.

* `bgp_route_translation_for_nat_enabled` - (Optional) Is BGP Route Translation for NAT enabled? Defaults to `false`.

* `dns_forwarding_enabled` - (Optional) Is DNS forwarding enabled?

* `ip_sec_replay_protection_enabled` - (Optional) Is IP Sec Replay Protection enabled? Defaults to `true`.

* `policy_group` - (Optional) One or more `policy_group` blocks as defined below.

* `remote_vnet_traffic_enabled` - (Optional) Is remote vnet traffic that is used to configure this gateway to accept traffic from other Azure Virtual Networks enabled? Defaults to `false`.

* `virtual_wan_traffic_enabled` - (Optional) Is remote vnet traffic that is used to configure this gateway to accept traffic from remote Virtual WAN networks enabled? Defaults to `false`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `vpn_client_configuration` - (Optional) A `vpn_client_configuration` block which is documented below. In this block the Virtual Network Gateway can be configured to accept IPSec point-to-site connections.

* `vpn_type` - (Optional) The routing type of the Virtual Network Gateway. Valid options are `RouteBased` or `PolicyBased`. Defaults to `RouteBased`. Changing this forces a new resource to be created.

---

The `ip_configuration` block supports:

* `name` - (Optional) A user-defined name of the IP configuration. Defaults to `vnetGatewayConfig`.

* `private_ip_address_allocation` - (Optional) Defines how the private IP address of the gateways virtual interface is assigned. The only valid value is `Dynamic` for Virtual Network Gateway (`Static` is not supported by the service yet). Defaults to `Dynamic`.

* `subnet_id` - (Required) The ID of the gateway subnet of a virtual network in which the virtual network gateway will be created. It is mandatory that the associated subnet is named `GatewaySubnet`. Therefore, each virtual network can contain at most a single Virtual Network Gateway.

* `public_ip_address_id` - (Required) The ID of the public IP address to associate with the Virtual Network Gateway.

---

The `policy_group` block supports:

* `name` - (Required) The name of the Virtual Network Gateway Policy Group.

* `policy_member` - (Required) One or more `policy_member` blocks as defined below.

* `is_default` - (Optional) Is this a Default Virtual Network Gateway Policy Group? Defaults to `false`.

* `priority` - (Optional) The priority for the Virtual Network Gateway Policy Group. Defaults to `0`.

---

The `policy_member` block supports:

* `name` - (Required) The name of the Virtual Network Gateway Policy Group Member.

* `type` - (Required) The VPN Policy Member attribute type. Possible values are `AADGroupId`, `CertificateGroupId` and `RadiusAzureGroupId`.

* `value` - (Required) The value of attribute that is used for this Virtual Network Gateway Policy Group Member.

---

The `vpn_client_configuration` block supports:

* `address_space` - (Required) The address space out of which IP addresses for vpn clients will be taken. You can provide more than one address space, e.g. in CIDR notation.

* `aad_tenant` - (Optional) AzureAD Tenant URL

* `aad_audience` - (Optional) The client id of the Azure VPN application.
    See [Create an Active Directory (AD) tenant for P2S OpenVPN protocol connections](https://docs.microsoft.com/en-gb/azure/vpn-gateway/openvpn-azure-ad-tenant-multi-app) for values

* `aad_issuer` - (Optional) The STS url for your tenant

* `ipsec_policy` - (Optional) An `ipsec_policy` block as defined below.

* `root_certificate` - (Optional) One or more `root_certificate` blocks which are defined below. These root certificates are used to sign the client certificate used by the VPN clients to connect to the gateway.

* `revoked_certificate` - (Optional) One or more `revoked_certificate` blocks which are defined below.

* `radius_server` - (Optional) One or more `radius_server` blocks as defined below.

* `radius_server_address` - (Optional) The address of the Radius server.

* `radius_server_secret` - (Optional) The secret used by the Radius server.

* `vpn_client_protocols` - (Optional) List of the protocols supported by the vpn client.
    The supported values are `SSTP`, `IkeV2` and `OpenVPN`.
    Values `SSTP` and `IkeV2` are incompatible with the use of
    `aad_tenant`, `aad_audience` and `aad_issuer`.

* `vpn_auth_types` - (Optional) List of the vpn authentication types for the virtual network gateway.
    The supported values are `AAD`, `Radius` and `Certificate`.

-> **Note:** `vpn_auth_types` must be set when using multiple vpn authentication types.

* `virtual_network_gateway_client_connection` - (Optional) One or more `virtual_network_gateway_client_connection` blocks as defined below.

---

The `bgp_settings` block supports:

* `asn` - (Optional) The Autonomous System Number (ASN) to use as part of the BGP.

* `peering_addresses` - (Optional) A list of `peering_addresses` blocks as defined below. Only one `peering_addresses` block can be specified except when `active_active` of this Virtual Network Gateway is `true`.

* `peer_weight` - (Optional) The weight added to routes which have been learned through BGP peering. Valid values can be between `0` and `100`.

---

A `custom_route` block supports the following:

* `address_prefixes` - (Optional) A list of address blocks reserved for this virtual network in CIDR notation.

---

A `peering_addresses` block supports the following:

* `ip_configuration_name` - (Optional) The name of the IP configuration of this Virtual Network Gateway. In case there are multiple `ip_configuration` blocks defined, this property is **required** to specify.

* `apipa_addresses` - (Optional) A list of Azure custom APIPA addresses assigned to the BGP peer of the Virtual Network Gateway.

~> **Note:** The valid range for the reserved APIPA address in Azure Public is from `169.254.21.0` to `169.254.22.255`.

---

The `root_certificate` block supports:

* `name` - (Required) A user-defined name of the root certificate.

* `public_cert_data` - (Required) The public certificate of the root certificate authority. The certificate must be provided in Base-64 encoded X.509 format (PEM). In particular, this argument *must not* include the `-----BEGIN CERTIFICATE-----` or `-----END CERTIFICATE-----` markers, nor any newlines.

---

The `revoked_certificate` block supports:

* `name` - (Required) Specifies the name of the certificate resource.

* `thumbprint` - (Required) Specifies the public data of the certificate.

---

The `ipsec_policy` block supports:

* `dh_group` - (Required) The DH Group, used in IKE Phase 1. Possible values are `DHGroup1`, `DHGroup2`, `DHGroup14`, `DHGroup24`, `DHGroup2048`, `ECP256`, `ECP384` and `None`.

* `ike_encryption` - (Required) The IKE encryption algorithm, used for IKE Phase 2. Possible values are `AES128`, `AES192`, `AES256`, `DES`, `DES3`, `GCMAES128` and `GCMAES256`.

* `ike_integrity` - (Required) The IKE encryption integrity algorithm, used for IKE Phase 2. Possible values are `GCMAES128`, `GCMAES256`, `MD5`, `SHA1`, `SHA256` and `SHA384`.

* `ipsec_encryption` - (Required) The IPSec encryption algorithm, used for IKE phase 1. Possible values are `AES128`, `AES192`, `AES256`, `DES`, `DES3`, `GCMAES128`, `GCMAES192`, `GCMAES256` and `None`.

* `ipsec_integrity` - (Required) The IPSec integrity algorithm, used for IKE phase 1. Possible values are `GCMAES128`, `GCMAES192`, `GCMAES256`, `MD5`, `SHA1` and `SHA256`.

* `pfs_group` - (Required) The Pfs Group, used in IKE Phase 2. Possible values are `ECP256`, `ECP384`, `PFS1`, `PFS2`, `PFS14`, `PFS24`, `PFS2048`, `PFSMM` and `None`.

* `sa_lifetime_in_seconds` - (Required) The IPSec Security Association lifetime in seconds for a Site-to-Site VPN tunnel. Possible values are between `300` and `172799`.

* `sa_data_size_in_kilobytes` - (Required) The IPSec Security Association payload size in KB for a Site-to-Site VPN tunnel. Possible values are between `1024` and `2147483647`.

---

The `radius_server` block supports:

* `address` - (Required) The address of the Radius Server.

* `secret` - (Required) The secret that is used to communicate with the Radius Server.

* `score` - (Required) The score of the Radius Server determines the priority of the server. Possible values are between `1` and `30`.

---

The `virtual_network_gateway_client_connection` block supports:

* `name` - (Required) The name of the Virtual Network Gateway Client Connection.

* `policy_group_names` - (Required) A list of names of Virtual Network Gateway Policy Groups.

* `address_prefixes` - (Required) A list of address prefixes for P2S VPN Client.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Network Gateway.

* `bgp_settings` - A block of `bgp_settings`.

---

The `bgp_settings` block supports:

* `peering_addresses` - A list of `peering_addresses` as defined below.

---

The `peering_addresses` block supports:

* `default_addresses` - A list of peering address assigned to the BGP peer of the Virtual Network Gateway.

* `tunnel_ip_addresses` - A list of tunnel IP addresses assigned to the BGP peer of the Virtual Network Gateway.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Virtual Network Gateway.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Network Gateway.
* `update` - (Defaults to 1 hour) Used when updating the Virtual Network Gateway.
* `delete` - (Defaults to 1 hour) Used when deleting the Virtual Network Gateway.

## Import

Virtual Network Gateways can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_network_gateway.exampleGateway /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.Network/virtualNetworkGateways/myGateway1
```
