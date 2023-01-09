---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_attached_data_network"
description: |-
  Manages a Mobile Network Attached Data Network.
---

# azurerm_mobile_network_attached_data_network

Manages a Mobile Network Attached Data Network.

## Example Usage

```hcl
data "azurerm_mobile_network_packet_core_control_plane" "example" {
  name                = "example-mnpccp"
  resource_group_name = "example-rg"
}

resource "azurerm_mobile_network_attached_data_network" "example" {
  name                                     = "example-data-network"
  mobile_network_packet_core_data_plane_id = data.azurerm_mobile_network_packet_core_control_plane.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Mobile Network Attached Data Network. Must be as same as `azurerm_mobile_network_data_network`, Changing this forces a new Mobile Network Attached Data Network to be created.

* `mobile_network_packet_core_data_plane_id` - (Required) Specifies the ID of the Mobile Network Attached Data Network. Changing this forces a new Mobile Network Attached Data Network to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Attached Data Network.

* `location` - Specifies the Azure Region where the Mobile Network Attached Data Network should exist. 

* `dns_addresses` - Specifies the DNS servers to signal to UEs to use for this attached data network.

* `network_address_port_translation_configuration` - A `network_address_port_translation_configuration` block as defined below.

* `user_plane_data_interface` - A `user_plane_data_interface` block as defined below.

* `user_equipment_address_pool_prefix` - Specifies the user equipment (UE) address pool prefixes for the attached data network from which the packet core instance will dynamically assign IP addresses to UEs. The packet core instance assigns an IP address to a UE when the UE sets up a PDU session.

* `user_equipment_static_address_pool_prefix` - Specifies the user equipment (UE) address pool prefixes for the attached data network from which the packet core instance will assign static IP addresses to UEs. The packet core instance assigns an IP address to a UE when the UE sets up a PDU session. The static IP address for a specific UE is set in StaticIPConfiguration on the corresponding SIM resource.

* `tags` - A mapping of tags which should be assigned to the Mobile Network Attached Data Network.

---

A `user_plane_data_interface` block supports the following:

* `name` - The logical name for this interface.

* `ipv4_address` - The IPv4 address.

* `ipv4_subnet` - The IPv4 subnet.

* `ipv4_gateway` - The default IPv4 gateway (router).

---

A `network_address_port_translation_configuration` block supports the following:

* `enabled` - Whether NAPT is enabled for connections to this attached data network.

* `pinhole_limits` - Maximum number of UDP and TCP pinholes that can be open simultaneously on the core interface. For 5G networks, this is the N6 interface. For 4G networks, this is the SGi interface.

* `pinhole_timeouts_in_seconds` - A `pinhole_timeouts_in_seconds` block as defined below.

* `port_range` - A `port_range` block as defined below.

* `port_reuse_minimum_hold_time_in_seconds` - A `port_reuse_minimum_hold_time_in_seconds` block as defined below.

---

A `pinhole_timeouts_in_seconds` block supports the following:

* `icmp` - Pinhole timeout for ICMP pinholes in seconds. Default for ICMP Echo is 60 seconds, as per RFC 5508 section 3.2.

* `tcp` - Pinhole timeout for TCP pinholes in seconds. Default for TCP is 2 hours 4 minutes, as per RFC 5382 section 5.

* `udp` - Pinhole timeout for UDP pinholes in seconds. Default for UDP is 5 minutes, as per RFC 4787 section 4.3.

---

A `port_range` block supports the following:

* `max_port` - Specifies the maximum port number.

* `min_port` - Specifies the minimum port number.

---

A `port_reuse_minimum_hold_time_in_seconds` block supports the following:

* `tcp` - Minimum time in seconds that will pass before a TCP port that was used by a closed pinhole can be reused.

* `udp` - Minimum time in seconds that will pass before a UDP port that was used by a closed pinhole can be reused.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 180 minutes) Used when creating the Mobile Network Attached Data Network.
* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Attached Data Network.
* `update` - (Defaults to 180 minutes) Used when updating the Mobile Network Attached Data Network.
* `delete` - (Defaults to 180 minutes) Used when deleting the Mobile Network Attached Data Network.

## Import

Mobile Network Attached Data Network can be imported using the `resource id`, e.g.
