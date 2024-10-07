---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_attached_data_network"
description: |-
  Gets information about an existing Mobile Network Attached Data Network.
---

# azurerm_mobile_network_attached_data_network

Use this data source to access information about an existing Mobile Network Attached Data Network.

## Example Usage

```hcl
data "azurerm_mobile_network_packet_core_control_plane" "example" {
  name                = "example-mnpccp"
  resource_group_name = "example-rg"
}

data "azurerm_mobile_network_attached_data_network" "example" {
  mobile_network_data_network_name         = data.azurerm_mobile_network_packet_core_control_plane.example.name
  mobile_network_packet_core_data_plane_id = data.azurerm_mobile_network_packet_core_control_plane.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `mobile_network_data_network_name` - The Name of the `azurerm_mobile_network_data_network` this resource belongs to.

* `mobile_network_packet_core_data_plane_id` - The ID of the `azurerm_mobile_network_packet_core_data_plane` which the Mobile Network Attached Data Network belongs to.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Attached Data Network.

* `location` - The Azure Region where the Mobile Network Attached Data Network should exist. 

* `dns_addresses` - The DNS servers to signal to UEs to use for this attached data network.

* `network_address_port_translation` - A `network_address_port_translation` block as defined below.

* `user_plane_access_name` - The logical name for thie user data plane interface.

* `user_plane_access_ipv4_address` - The IPv4 address for the user data plane interface.

* `user_plane_access_ipv4_subnet` - The IPv4 subnet for the user data plane interface.

* `user_plane_access_ipv4_gateway` - The default IPv4 gateway for the user data plane interface.

* `user_equipment_address_pool_prefix` - The user equipment (UE) address pool prefixes for the attached data network from which the packet core instance will dynamically assign IP addresses to UEs. The packet core instance assigns an IP address to a UE when the UE sets up a PDU session.

* `user_equipment_static_address_pool_prefix` - The user equipment (UE) address pool prefixes for the attached data network from which the packet core instance will assign static IP addresses to UEs. The packet core instance assigns an IP address to a UE when the UE sets up a PDU session. The static IP address for a specific UE is set in StaticIPConfiguration on the corresponding SIM resource.

* `tags` - A mapping of tags which should be assigned to the Mobile Network Attached Data Network.

---

A `network_address_port_translation` block supports the following:

* `pinhole_limits` - Maximum number of UDP and TCP pinholes that can be open simultaneously on the core interface. For 5G networks, this is the N6 interface. For 4G networks, this is the SGi interface.

* `icmp_pinhole_timeouts_in_seconds` - Pinhole timeout for ICMP pinholes in seconds. 

* `tcp_pinhole_timeouts_in_seconds` - Pinhole timeout for TCP pinholes in seconds.

* `udp_pinhole_timeouts_in_seconds` - Pinhole timeout for UDP pinholes in seconds.

* `port_range` - A `port_range` block as defined below.

* `tcp_port_reuse_minimum_hold_time_in_seconds` - Minimum time in seconds that will pass before a TCP port that was used by a closed pinhole can be reused.

* `udp_port_reuse_minimum_hold_time_in_seconds` - Minimum time in seconds that will pass before a UDP port that was used by a closed pinhole can be reused.

---

A `port_range` block supports the following:

* `maximum` - The maximum port number.

* `minimum` - The minimum port number.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Attached Data Network.
