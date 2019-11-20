---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub"
sidebar_current: "docs-azurerm-resource-virtual-hub"
description: |-
  Manages a Virtual Hub.
---

# azurerm_virtual_hub

Manages a Virtual Hub.


## Virtual Hub Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_virtual_wan" "example" {
  name                = "example-virtualwan"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
}
resource "azurerm_virtual_hub" "example" {
  name                = "example-virtualhub"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  address_prefix      = "10.0.1.0/24"
  virtual_wan_id      = "${azurerm_virtual_wan.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Virtual Hub. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the Virtual Hub should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `address_prefix` - (Required) Address-prefix for this Virtual Hub.

* `virtual_wan_id` - (Required) The resource id of virtual wan.

* `s2s_vpn_gateway_id` - (Optional) The resource id of s2s vpn gateway.

* `p2s_vpn_gateway_id` - (Optional) The resource id of p2s vpn gateway.

* `express_route_gateway_id` - (Optional) The resource id of express route gateway.

* `virtual_network_connection` - (Optional) One or more `virtual_network_connection` block defined below.

* `route` - (Optional) One `route` block defined below.

* `tags` - (Optional) Resource tags. Changing this forces a new resource to be created.

---

The `route` block supports the following:

* `address_prefixes` - (Required) List of all addressPrefixes.

* `next_hop_ip_address` - (Required) NextHop ip address.

---

The `virtual_network_connection` block supports the following:

* `name` - (Required) The name of the resource that is unique within a resource group. This name can be used to access the resource.

* `remote_virtual_network_id` - (Required) The resource id of remote virtual network.

* `allow_hub_to_remote_vnet_transit` - (Optional) VirtualHub to RemoteVnet transit to enabled or not.

* `allow_remote_vnet_to_use_hub_vnet_gateways` - (Optional) Allow RemoteVnet to use Virtual Hub's gateways.

* `enable_internet_security` - (Optional) Enable internet security.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Hub.

## Import

Virtual Hub can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_virtual_hub.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/vhub1
```