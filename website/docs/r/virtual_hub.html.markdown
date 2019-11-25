---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub"
sidebar_current: "docs-azurerm-resource-virtual-hub"
description: |-
  Manages a Virtual Hub within a Virtual WAN.
---

# azurerm_virtual_hub

Manages a Virtual Hub within a Virtual WAN.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_wan" "example" {
  name                = "example-virtualwan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_hub" "example" {
  name                = "example-virtualhub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_prefix      = "10.0.1.0/24"
  virtual_wan_id      = azurerm_virtual_wan.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Virtual Hub. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Virtual Hub should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the Virtual Hub should exist. Changing this forces a new resource to be created.

* `address_prefix` - (Required) The Address Prefix which should be used for this Virtual Hub.

* `virtual_wan_id` - (Required) The ID of a Virtual WAN within which the Virtual Hub should be created.

---

* `express_route_gateway_id` - (Optional) The ID of an Express Route Gateway which should be used for Express Route connections.

~> **NOTE:** This functionality is in Preview and must be opted into via `az feature register --namespace Microsoft.Network --name AllowCortexExpressRouteGateway` and then `az provider register -n Microsoft.Network`.

* `p2s_vpn_gateway_id` - (Optional) The ID of a Point-to-Site VPN Gateway which should be used for Point-to-Site connections.

* `route` - (Optional) One or more `route` blocks as defined below.

* `s2s_vpn_gateway_id` - (Optional) The ID of a Site-to-Site VPN Gateway which should be used for Site-to-Site connections.

* `tags` - (Optional) A mapping of tags to assign to the Virtual Hub.

* `virtual_network_connection` - (Optional) One or more `virtual_network_connection` blocks as defined below.

---

The `route` block supports the following:

* `address_prefixes` - (Required) A list of Address Prefixes.

* `next_hop_ip_address` - (Required) The IP Address that Packets should be forwarded to as the Next Hop.

---

The `virtual_network_connection` block supports the following:

* `name` - (Required) The name of the resource that is unique within a resource group. This name can be used to access the resource.

* `remote_virtual_network_id` - (Required) The ID of a Virtual Network.

* `allow_hub_to_remote_vnet_transit` - (Optional) Should the Virtual Hub be able to transit via Remote Virtual Networks?

* `allow_remote_vnet_to_use_hub_vnet_gateways` - (Optional) Should the Remote Virtual Network be able to use the Hub's Virtual Network Gateways?

* `enable_internet_security` - (Optional) Should internet security be enabled?

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Hub.

## Import

Virtual Hub can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_virtual_hub.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/vhub1
```
