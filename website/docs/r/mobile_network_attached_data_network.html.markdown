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
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_databox_edge_device" "example" {
  name                = "example-device"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku_name = "EdgeP_Base-Standard"
}

resource "azurerm_mobile_network" "example" {
  name                = "example-mn"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  mobile_country_code = "001"
  mobile_network_code = "01"
}

resource "azurerm_mobile_network_packet_core_control_plane" "example" {
  name                = "example-mnpccp"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
  sku                 = "G0"
  mobile_network_id   = azurerm_mobile_network.example.id

  control_plane_access_interface {
    name         = "default-interface"
    ipv4_address = "192.168.1.199"
    ipv4_gateway = "192.168.1.1"
    ipv4_subnet  = "192.168.1.0/25"
  }

  platform {
    type = "AKS-HCI"
    edge_device_id = azurerm_databox_edge_device.example.id
  }
}

resource "azurerm_mobile_network_packet_core_data_plane" "example" {
  name                                        = "example-mnpcdp"
  mobile_network_packet_core_control_plane_id = azurerm_mobile_network_packet_core_control_plane.example.id
  location                                    = azurerm_resource_group.example.location

  user_plane_access_interface {
    name         = "default-interface"
    ipv4_address = "192.168.1.199"
    ipv4_gateway = "192.168.1.1"
    ipv4_subnet  = "192.168.1.0/25"
  }
}

resource "azurerm_mobile_network_data_network" "example" {
  name              = "example-data-network"
  mobile_network_id = azurerm_mobile_network.example.id
  location          = azurerm_resource_group.example.location
}

resource "azurerm_mobile_network_attached_data_network" "example" {
  name                                        = "example-data-network"
  mobile_network_packet_core_data_plane_id    = azurerm_mobile_network_packet_core_data_plane.example.id
  location                                    = azurerm_resource_group.example.location
  dns_addresses                               = ["1.1.1.1"]
  user_equipment_address_pool_prefixes        = ["2.4.0.0/16"]
  user_equipment_static_address_pool_prefixes = ["2.4.0.0/16"]

  network_address_port_translation_configuration {
    enabled                = true
    pinhole_maximum_number = 65536
    pinhole_timeouts_in_seconds {
      icmp = 30
      tcp  = 100
      udp  = 39
    }
    port_range {
      max_port = 49999
      min_port = 1024
    }
    port_reuse_minimum_hold_time_in_seconds {
      tcp = 120
      udp = 60
    }
  }

  user_plane_data_interface {
    name         = "test"
    ipv4_address = "10.204.141.4"
    ipv4_gateway = "10.204.141.1"
    ipv4_subnet  = "10.204.141.0/24"
  }

  tags = {
    key = "value"
  }

  depends_on = [azurerm_mobile_network_data_network.example]

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Mobile Network Attached Data Network. Must be as same as `azurerm_mobile_network_data_network`, Changing this forces a new Mobile Network Attached Data Network to be created.

* `mobile_network_packet_core_data_plane_id` - (Required) Specifies the ID of the Mobile Network Attached Data Network. Changing this forces a new Mobile Network Attached Data Network to be created.

* `location` - (Required) Specifies the Azure Region where the Mobile Network Attached Data Network should exist. Changing this forces a new Mobile Network Attached Data Network to be created.

* `dns_addresses` - (Required) Specifies the DNS servers to signal to UEs to use for this attached data network.

* `network_address_port_translation_configuration` - (Optional) A `network_address_port_translation_configuration` block as defined below.

* `user_plane_data_interface` - (Optinal) A `user_plane_data_interface` block as defined below.

* `user_equipment_address_pool_prefix` - (Optional) Specifies the user equipment (UE) address pool prefixes for the attached data network from which the packet core instance will dynamically assign IP addresses to UEs. The packet core instance assigns an IP address to a UE when the UE sets up a PDU session. At least one of `user_equipment_address_pool_prefix` and `user_equipment_static_address_pool_prefix`. If you define both, they must be of the same size.

* `user_equipment_static_address_pool_prefix` - (Optional) Specifies the user equipment (UE) address pool prefixes for the attached data network from which the packet core instance will assign static IP addresses to UEs. The packet core instance assigns an IP address to a UE when the UE sets up a PDU session. The static IP address for a specific UE is set in StaticIPConfiguration on the corresponding SIM resource. At least one of `user_equipment_address_pool_prefix` and `user_equipment_static_address_pool_prefix`. If you define both, they must be of the same size.

* `tags` - (Optional) A mapping of tags which should be assigned to the Mobile Network Attached Data Network.

---

A `user_plane_data_interface` block supports the following:

* `name` - (Optional) Specifies the logical name for this interface. This should match one of the interfaces configured on your Azure Stack Edge device.

* `ipv4_address` - (Optional) The IPv4 address.

* `ipv4_subnet` - (Optional) The IPv4 subnet.

* `ipv4_gateway` - (Optional) The default IPv4 gateway (router).

---

A `network_address_port_translation_configuration` block supports the following:

* `enabled` - (Required) Whether NAPT is enabled for connections to this attached data network.

* `pinhole_limits` - (Optional) Maximum number of UDP and TCP pinholes that can be open simultaneously on the core interface. For 5G networks, this is the N6 interface. For 4G networks, this is the SGi interface. Must be between 1 and 65536.

* `pinhole_timeouts_in_seconds` - (Optional) A `pinhole_timeouts_in_seconds` block as defined below.

* `port_range` - (Optional) A `port_range` block as defined below.

* `port_reuse_minimum_hold_time_in_seconds` - (Optional) A `port_reuse_minimum_hold_time_in_seconds` block as defined below.

---

A `pinhole_timeouts_in_seconds` block supports the following:

* `icmp` - (Optional) Pinhole timeout for ICMP pinholes in seconds. Default for ICMP Echo is 60 seconds, as per RFC 5508 section 3.2.

* `tcp` - (Optional) Pinhole timeout for TCP pinholes in seconds. Default for TCP is 2 hours 4 minutes, as per RFC 5382 section 5.

* `udp` - (Optional) Pinhole timeout for UDP pinholes in seconds. Default for UDP is 5 minutes, as per RFC 4787 section 4.3.

---

A `port_range` block supports the following:

* `max_port` - (Optional) Specifies the maximum port number.

* `min_port` - (Optional) Specifies the minimum port number.

---

A `port_reuse_minimum_hold_time_in_seconds` block supports the following:

* `tcp` - (Optional) Minimum time in seconds that will pass before a TCP port that was used by a closed pinhole can be reused. Default for TCP is 2 minutes.

* `udp` - (Optional) Minimum time in seconds that will pass before a UDP port that was used by a closed pinhole can be reused. Default for UDP is 1 minute.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Attached Data Network.



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 180 minutes) Used when creating the Mobile Network Attached Data Network.
* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Attached Data Network.
* `update` - (Defaults to 180 minutes) Used when updating the Mobile Network Attached Data Network.
* `delete` - (Defaults to 180 minutes) Used when deleting the Mobile Network Attached Data Network.

## Import

Mobile Network Attached Data Network can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mobile_network_attached_data_network.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.MobileNetwork/packetCoreControlPlanes/packetCoreControlPlane1/packetCoreDataPlanes/packetCoreDataPlane1/attachedDataNetworks/attachedDataNetwork1
```
