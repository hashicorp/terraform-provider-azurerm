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

* `name` - (Required) Specifies the name which should be used for this Mobile Network Packet Core Data Plane. Changing this forces a new Mobile Network Packet Core Data Plane to be created.

* `mobile_network_packet_core_control_plane_id` - (Required) Specifies the ID of the Mobile Network Packet Core Data Plane. Changing this forces a new Mobile Network Packet Core Data Plane to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Packet Core Data Plane.

* `location` - The Azure Region where the Mobile Network Packet Core Data Plane should exist.

* `user_plane_access_interface` - A `user_plane_access_interface` block as defined below.

* `tags` - A mapping of tags which should be assigned to the Mobile Network Packet Core Data Plane.

---

A `user_plane_access_interface` block supports the following:

* `name` - The logical name for this interface. 

* `ipv4_address` - The IPv4 address.

* `ipv4_subnet` - The IPv4 subnet.

* `ipv4_gateway` - The default IPv4 gateway (router).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Packet Core Data Plane.

