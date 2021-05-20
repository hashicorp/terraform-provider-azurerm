---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_network_gateway_connection"
description: |-
  Manages a connection in an existing Virtual Network Gateway.
---

# azurerm_virtual_network_gateway_connection

Manages a connection in an existing Virtual Network Gateway.

## Example Usage

### Site-to-Site connection

The following example shows a connection between an Azure virtual network
and an on-premises VPN device and network.

```hcl
resource "azurerm_resource_group" "example" {
  name     = "test"
  location = "West US"
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

resource "azurerm_local_network_gateway" "onpremise" {
  name                = "onpremise"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  gateway_address     = "168.62.225.23"
  address_space       = ["10.1.1.0/24"]
}

resource "azurerm_public_ip" "example" {
  name                = "test"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Dynamic"
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
    public_ip_address_id          = azurerm_public_ip.example.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.example.id
  }
}

resource "azurerm_virtual_network_gateway_connection" "onpremise" {
  name                = "onpremise"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.example.id
  local_network_gateway_id   = azurerm_local_network_gateway.onpremise.id

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}
```

### VNet-to-VNet connection

The following example shows a connection between two Azure virtual network
in different locations/regions.

```hcl
resource "azurerm_resource_group" "us" {
  name     = "us"
  location = "East US"
}

resource "azurerm_virtual_network" "us" {
  name                = "us"
  location            = azurerm_resource_group.us.location
  resource_group_name = azurerm_resource_group.us.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "us_gateway" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.us.name
  virtual_network_name = azurerm_virtual_network.us.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "us" {
  name                = "us"
  location            = azurerm_resource_group.us.location
  resource_group_name = azurerm_resource_group.us.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "us" {
  name                = "us-gateway"
  location            = azurerm_resource_group.us.location
  resource_group_name = azurerm_resource_group.us.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.us.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.us_gateway.id
  }
}

resource "azurerm_resource_group" "europe" {
  name     = "europe"
  location = "West Europe"
}

resource "azurerm_virtual_network" "europe" {
  name                = "europe"
  location            = azurerm_resource_group.europe.location
  resource_group_name = azurerm_resource_group.europe.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "europe_gateway" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.europe.name
  virtual_network_name = azurerm_virtual_network.europe.name
  address_prefixes     = ["10.1.1.0/24"]
}

resource "azurerm_public_ip" "europe" {
  name                = "europe"
  location            = azurerm_resource_group.europe.location
  resource_group_name = azurerm_resource_group.europe.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "europe" {
  name                = "europe-gateway"
  location            = azurerm_resource_group.europe.location
  resource_group_name = azurerm_resource_group.europe.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.europe.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.europe_gateway.id
  }
}

resource "azurerm_virtual_network_gateway_connection" "us_to_europe" {
  name                = "us-to-europe"
  location            = azurerm_resource_group.us.location
  resource_group_name = azurerm_resource_group.us.name

  type                            = "Vnet2Vnet"
  virtual_network_gateway_id      = azurerm_virtual_network_gateway.us.id
  peer_virtual_network_gateway_id = azurerm_virtual_network_gateway.europe.id

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}

resource "azurerm_virtual_network_gateway_connection" "europe_to_us" {
  name                = "europe-to-us"
  location            = azurerm_resource_group.europe.location
  resource_group_name = azurerm_resource_group.europe.name

  type                            = "Vnet2Vnet"
  virtual_network_gateway_id      = azurerm_virtual_network_gateway.europe.id
  peer_virtual_network_gateway_id = azurerm_virtual_network_gateway.us.id

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the connection. Changing the name forces a
    new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the connection Changing the name forces a new resource to be created.

* `location` - (Required) The location/region where the connection is
    located. Changing this forces a new resource to be created.

* `type` - (Required) The type of connection. Valid options are `IPsec`
    (Site-to-Site), `ExpressRoute` (ExpressRoute), and `Vnet2Vnet` (VNet-to-VNet).
    Each connection type requires different mandatory arguments (refer to the
    examples above). Changing the connection type will force a new connection
    to be created.

* `virtual_network_gateway_id` - (Required) The ID of the Virtual Network Gateway
    in which the connection will be created. Changing the gateway forces a new
    resource to be created.

* `authorization_key` - (Optional) The authorization key associated with the
    Express Route Circuit. This field is required only if the type is an
    ExpressRoute connection.

* `dpd_timeout_seconds` - (Optional) The dead peer detection timeout of this connection in seconds. Changing this forces a new resource to be created.

* `express_route_circuit_id` - (Optional) The ID of the Express Route Circuit
    when creating an ExpressRoute connection (i.e. when `type` is `ExpressRoute`).
    The Express Route Circuit can be in the same or in a different subscription.

* `peer_virtual_network_gateway_id` - (Optional) The ID of the peer virtual
    network gateway when creating a VNet-to-VNet connection (i.e. when `type`
    is `Vnet2Vnet`). The peer Virtual Network Gateway can be in the same or
    in a different subscription.

* `local_azure_ip_address_enabled` - (Optional) Use private local Azure IP for the connection. Changing this forces a new resource to be created.

* `local_network_gateway_id` - (Optional) The ID of the local network gateway
    when creating Site-to-Site connection (i.e. when `type` is `IPsec`).

* `routing_weight` - (Optional) The routing weight. Defaults to `10`.

* `shared_key` - (Optional) The shared IPSec key. A key could be provided if a
    Site-to-Site, VNet-to-VNet or ExpressRoute connection is created.

* `connection_protocol` - (Optional) The IKE protocol version to use. Possible
    values are `IKEv1` and `IKEv2`. Defaults to `IKEv2`.
    Changing this value will force a resource to be created.
-> **Note**: Only valid for `IPSec` connections on virtual network gateways with SKU `VpnGw1`, `VpnGw2`, `VpnGw3`, `VpnGw1AZ`, `VpnGw2AZ` or `VpnGw3AZ`.

* `enable_bgp` - (Optional) If `true`, BGP (Border Gateway Protocol) is enabled
    for this connection. Defaults to `false`.

* `express_route_gateway_bypass` - (Optional) If `true`, data packets will bypass ExpressRoute Gateway for data forwarding This is only valid for ExpressRoute connections.

* `use_policy_based_traffic_selectors` - (Optional) If `true`, policy-based traffic
    selectors are enabled for this connection. Enabling policy-based traffic
    selectors requires an `ipsec_policy` block. Defaults to `false`.

* `ipsec_policy` (Optional) A `ipsec_policy` block which is documented below.
    Only a single policy can be defined for a connection. For details on
    custom policies refer to [the relevant section in the Azure documentation](https://docs.microsoft.com/en-us/azure/vpn-gateway/vpn-gateway-ipsecikepolicy-rm-powershell).

* `traffic_selector_policy` A `traffic_selector_policy` which allows to specify traffic selector policy proposal to be used in a virtual network gateway connection.
    Only one block can be defined for a connection.
    For details about traffic selectors refer to [the relevant section in the Azure documentation](https://docs.microsoft.com/en-us/azure/vpn-gateway/vpn-gateway-connect-multiple-policybased-rm-ps).

* `tags` - (Optional) A mapping of tags to assign to the resource.

The `ipsec_policy` block supports:

* `dh_group` - (Required) The DH group used in IKE phase 1 for initial SA. Valid
    options are `DHGroup1`, `DHGroup14`, `DHGroup2`, `DHGroup2048`, `DHGroup24`,
    `ECP256`, `ECP384`, or `None`.

* `ike_encryption` - (Required) The IKE encryption algorithm. Valid
    options are `AES128`, `AES192`, `AES256`, `DES`, `DES3`, `GCMAES128`, or `GCMAES256`.

* `ike_integrity` - (Required) The IKE integrity algorithm. Valid
    options are `GCMAES128`, `GCMAES256`, `MD5`, `SHA1`, `SHA256`, or `SHA384`.

* `ipsec_encryption` - (Required) The IPSec encryption algorithm. Valid
    options are `AES128`, `AES192`, `AES256`, `DES`, `DES3`, `GCMAES128`, `GCMAES192`, `GCMAES256`, or `None`.

* `ipsec_integrity` - (Required) The IPSec integrity algorithm. Valid
    options are `GCMAES128`, `GCMAES192`, `GCMAES256`, `MD5`, `SHA1`, or `SHA256`.

* `pfs_group` - (Required) The DH group used in IKE phase 2 for new child SA.
    Valid options are `ECP256`, `ECP384`, `PFS1`, `PFS14`, `PFS2`, `PFS2048`, `PFS24`, `PFSMM`,
    or `None`.

* `sa_datasize` - (Optional) The IPSec SA payload size in KB. Must be at least
    `1024` KB. Defaults to `102400000` KB.

* `sa_lifetime` - (Optional) The IPSec SA lifetime in seconds. Must be at least
    `300` seconds. Defaults to `27000` seconds.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Network Gateway Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Network Gateway Connection.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Network Gateway Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Network Gateway Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Network Gateway Connection.

## Import

Virtual Network Gateway Connections can be imported using their `resource id`, e.g.

```
terraform import azurerm_virtual_network_gateway_connection.exampleConnection /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.Network/connections/myConnection1
```
