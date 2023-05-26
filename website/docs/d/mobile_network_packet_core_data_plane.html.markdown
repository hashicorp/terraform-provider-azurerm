---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_packet_core_data_plane"
description: |-
  Get information a Mobile Network Packet Core Data Plane.
---

# azurerm_mobile_network_packet_core_data_plane

Get information a Mobile Network Packet Core Data Plane.

## Example Usage

```hcl
data "azurerm_mobile_network_packet_core_control_plane" "example" {
  name                = "example-mnpccp"
  resource_group_name = "example-rg"
}

data "azurerm_mobile_network_packet_core_data_plane" "example" {
  name                                        = "example-mnpcdp"
  mobile_network_packet_core_control_plane_id = data.azurerm_mobile_network_packet_core_control_plane.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Mobile Network Packet Core Data Plane. 

* `mobile_network_packet_core_control_plane_id` - (Required) The ID of the Mobile Network Packet Core Data Plane.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Packet Core Data Plane.

* `location` - The Azure Region where the Mobile Network Packet Core Data Plane should exist.

* `user_plane_access_name` - The logical name for thie user plane interface.

* `user_plane_access_ipv4_address` - The IPv4 address for the user plane interface.

* `user_plane_access_ipv4_subnet` - The IPv4 subnet for the user plane interface.

* `user_plane_access_ipv4_gateway` - The default IPv4 gateway for the user plane interface.

* `tags` - A mapping of tags which should be assigned to the Mobile Network Packet Core Data Plane.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Packet Core Data Plane.

