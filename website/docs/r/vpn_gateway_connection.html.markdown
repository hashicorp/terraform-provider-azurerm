---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_vpn_gateway_connection"
description: |-
  Manages a VPN Gateway Connection.
---

# azurerm_vpn_gateway_connection

Manages a VPN Gateway Connection.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
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
  address_prefix      = "10.0.0.0/24"
}

resource "azurerm_vpn_gateway" "example" {
  name                = "example-vpng"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  virtual_hub_id      = azurerm_virtual_hub.example.id
}

resource "azurerm_vpn_site" "example" {
  name                = "example-vpn-site"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  virtual_wan_id      = azurerm_virtual_wan.example.id
  link {
    name       = "link1"
    ip_address = "10.1.0.0"
  }
  link {
    name       = "link2"
    ip_address = "10.2.0.0"
  }
}

resource "azurerm_vpn_gateway_connection" "example" {
  name               = "example"
  vpn_gateway_id     = azurerm_vpn_gateway.example.id
  remote_vpn_site_id = azurerm_vpn_site.example.id

  vpn_link {
    name             = "link1"
    vpn_site_link_id = azurerm_vpn_site.example.link[0].id
  }

  vpn_link {
    name             = "link2"
    vpn_site_link_id = azurerm_vpn_site.example.link[1].id
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this VPN Gateway Connection. Changing this forces a new VPN Gateway Connection to be created.

* `remote_vpn_site_id` - (Required) The ID of the remote VPN Site, which will connect to the VPN Gateway. Changing this forces a new VPN Gateway Connection to be created.

* `vpn_gateway_id` - (Required) The ID of the VPN Gateway that this VPN Gateway Connection belongs to. Changing this forces a new VPN Gateway Connection to be created.

* `vpn_link` - (Required) One or more `vpn_link` blocks as defined below.

* `internet_security_enabled` - (Optional) Whether Internet Security is enabled for this VPN Connection. Defaults to `false`.

* `routing` - (Optional) A `routing` block as defined below. If this is not specified, there will be a default route table created implicitly.

* `traffic_selector_policy` - (Optional) One or more `traffic_selector_policy` blocks as defined below.

---

A `ipsec_policy` block supports the following:

* `dh_group` - (Required) The DH Group used in IKE Phase 1 for initial SA. Possible values are `None`, `DHGroup1`, `DHGroup2`, `DHGroup14`, `DHGroup24`, `DHGroup2048`, `ECP256`, `ECP384`.

* `ike_encryption_algorithm` - (Required) The IKE encryption algorithm (IKE phase 2). Possible values are `DES`, `DES3`, `AES128`, `AES192`, `AES256`, `GCMAES128`, `GCMAES256`.

* `ike_integrity_algorithm` - (Required) The IKE integrity algorithm (IKE phase 2). Possible values are `MD5`, `SHA1`, `SHA256`, `SHA384`, `GCMAES128`, `GCMAES256`.

* `encryption_algorithm` - (Required) The IPSec encryption algorithm (IKE phase 1). Possible values are `AES128`, `AES192`, `AES256`, `DES`, `DES3`, `GCMAES128`, `GCMAES192`, `GCMAES256`, `None`.

* `integrity_algorithm` - (Required) The IPSec integrity algorithm (IKE phase 1). Possible values are `MD5`, `SHA1`, `SHA256`, `GCMAES128`, `GCMAES192`, `GCMAES256`.

* `pfs_group` - (Required) The Pfs Group used in IKE Phase 2 for the new child SA. Possible values are `None`, `PFS1`, `PFS2`, `PFS14`, `PFS24`, `PFS2048`, `PFSMM`, `ECP256`, `ECP384`.

* `sa_data_size_kb` - (Required) The IPSec Security Association (also called Quick Mode or Phase 2 SA) payload size in KB for the site to site VPN tunnel.

* `sa_lifetime_sec` - (Required) The IPSec Security Association (also called Quick Mode or Phase 2 SA) lifetime in seconds for the site to site VPN tunnel.

---

A `vpn_link` block supports the following:

* `name` - (Required) The name which should be used for this VPN Link Connection.

* `egress_nat_rule_ids` - (Optional) A list of the egress NAT Rule Ids.

* `ingress_nat_rule_ids` - (Optional) A list of the ingress NAT Rule Ids.

* `vpn_site_link_id` - (Required) The ID of the connected VPN Site Link. Changing this forces a new VPN Gateway Connection to be created.

* `bandwidth_mbps` - (Optional) The expected connection bandwidth in MBPS. Defaults to `10`.

* `bgp_enabled` - (Optional) Should the BGP be enabled? Defaults to `false`. Changing this forces a new VPN Gateway Connection to be created.

* `connection_mode` - (Optional) The connection mode of this VPN Link. Possible values are `Default`, `InitiatorOnly` and `ResponderOnly`. Defaults to `Default`.

* `ipsec_policy` - (Optional) One or more `ipsec_policy` blocks as defined above.

* `protocol` - (Optional) The protocol used for this VPN Link Connection. Possible values are `IKEv1` and `IKEv2`. Defaults to `IKEv2`.

* `ratelimit_enabled` - (Optional) Should the rate limit be enabled? Defaults to `false`.

* `route_weight` - (Optional) Routing weight for this VPN Link Connection. Defaults to `0`.

* `shared_key` - (Optional) SharedKey for this VPN Link Connection.

* `local_azure_ip_address_enabled` - (Optional) Whether to use local Azure IP to initiate connection? Defaults to `false`.

* `policy_based_traffic_selector_enabled` - (Optional) Whether to enable policy-based traffic selectors? Defaults to `false`.

* `custom_bgp_address` - (Optional) One or more `custom_bgp_address` blocks as defined below.

---

A `routing` block supports the following:

* `associated_route_table` - (Required) The ID of the Route Table associated with this VPN Connection.

* `propagated_route_table` - (Optional) A `propagated_route_table` block as defined below.

* `inbound_route_map_id` - (Optional) The resource ID of the Route Map associated with this Routing Configuration for inbound learned routes.

* `outbound_route_map_id` - (Optional) The resource ID of the Route Map associated with this Routing Configuration for outbound advertised routes.

---

A `traffic_selector_policy` block supports the following:

* `local_address_ranges` - (Required) A list of local address spaces in CIDR format for this VPN Gateway Connection.

* `remote_address_ranges` - (Required) A list of remote address spaces in CIDR format for this VPN Gateway Connection.

---

A `propagated_route_table` block supports the following:

* `route_table_ids` - (Required) A list of Route Table IDs to associated with this VPN Gateway Connection.

* `labels` - (Optional) A list of labels to assign to this route table.

---

A `custom_bgp_address` block supports the following:

* `ip_address` - (Required) The custom bgp ip address which belongs to the IP Configuration.

* `ip_configuration_id` - (Required) The ID of the IP Configuration which belongs to the VPN Gateway.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the VPN Gateway Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the VPN Gateway Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the VPN Gateway Connection.
* `update` - (Defaults to 30 minutes) Used when updating the VPN Gateway Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the VPN Gateway Connection.

## Import

VPN Gateway Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_vpn_gateway_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/vpnGateways/gateway1/vpnConnections/conn1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network`: 2024-05-01
