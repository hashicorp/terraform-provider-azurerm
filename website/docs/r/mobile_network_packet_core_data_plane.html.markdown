---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_packet_core_data_plane"
description: |-
  Manages a Mobile Network Packet Core Data Plane.
---

# azurerm_mobile_network_packet_core_data_plane

Manages a Mobile Network Packet Core Data Plane.

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
  name                              = "example-mnpccp"
  resource_group_name               = azurerm_resource_group.example.name
  location                          = azurerm_resource_group.example.location
  sku                               = "G0"
  mobile_network_id                 = azurerm_mobile_network.example.id
  control_plane_access_name         = "default-interface"
  control_plane_access_ipv4_address = "192.168.1.199"
  control_plane_access_ipv4_gateway = "192.168.1.1"
  control_plane_access_ipv4_subnet  = "192.168.1.0/25"

  platform {
    type           = "AKS-HCI"
    edge_device_id = azurerm_databox_edge_device.example.id
  }
}

resource "azurerm_mobile_network_packet_core_data_plane" "example" {
  name                                        = "example-mnpcdp"
  mobile_network_packet_core_control_plane_id = azurerm_mobile_network_packet_core_control_plane.example.id
  location                                    = azurerm_resource_group.example.location
  user_plane_access_name                      = "default-interface"
  user_plane_access_ipv4_address              = "192.168.1.199"
  user_plane_access_ipv4_gateway              = "192.168.1.1"
  user_plane_access_ipv4_subnet               = "192.168.1.0/25"

  tags = {
    key = "value"
  }

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Mobile Network Packet Core Data Plane. Changing this forces a new Mobile Network Packet Core Data Plane to be created.

* `mobile_network_packet_core_control_plane_id` - (Required) Specifies the ID of the Mobile Network Packet Core Data Plane. Changing this forces a new Mobile Network Packet Core Data Plane to be created.

* `location` - (Required) Specifies the Azure Region where the Mobile Network Packet Core Data Plane should exist. Changing this forces a new Mobile Network Packet Core Data Plane to be created.

* `user_plane_access_name` - (Optional) Specifies the logical name for thie user plane interface. This should match one of the interfaces configured on your Azure Stack Edge device.

* `user_plane_access_ipv4_address` - (Optional) The IPv4 address for the user plane interface. This should match one of the interfaces configured on your Azure Stack Edge device.

* `user_plane_access_ipv4_subnet` - (Optional) The IPv4 subnet for the user plane interface. This should match one of the interfaces configured on your Azure Stack Edge device.

* `user_plane_access_ipv4_gateway` - (Optional) The default IPv4 gateway for the user plane interface. This should match one of the interfaces configured on your Azure Stack Edge device.

* `tags` - (Optional) A mapping of tags which should be assigned to the Mobile Network Packet Core Data Plane.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Packet Core Data Plane.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Mobile Network Packet Core Data Plane.
* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Packet Core Data Plane.
* `update` - (Defaults to 3 hours) Used when updating the Mobile Network Packet Core Data Plane.
* `delete` - (Defaults to 3 hours) Used when deleting the Mobile Network Packet Core Data Plane.

## Import

Mobile Network Packet Core Data Plane can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mobile_network_packet_core_data_plane.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.MobileNetwork/packetCoreControlPlanes/packetCoreControlPlane1/packetCoreDataPlanes/packetCoreDataPlane1
```
