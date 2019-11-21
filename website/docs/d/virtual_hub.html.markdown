---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub"
sidebar_current: "docs-azurerm-datasource-virtual-hub"
description: |-
  Gets information about an existing Virtual Hub
---

# Data Source: azurerm_virtual_hub

Uses this data source to access information about an existing Virtual Hub.


## Virtual Hub Usage

```hcl
data "azurerm_virtual_hub" "example" {
  resource_group = "acctestRG"
  name           = "acctestvirtualhub"
}
output "virtual_hub_id" {
  value = "${data.azurerm_virtual_hub.example.id}"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Virtual Hub.

* `resource_group_name` - (Required) The Name of the Resource Group where the Virtual Hub exists.


## Attributes Reference

The following attributes are exported:

* `location` - The Azure Region where the Virtual Hub exists.

* `address_prefix` - Address-prefix for this Virtual Hub.

* `virtual_wan_id` - The resource id of virtual wan.

* `s2s_vpn_gateway_id` - The resource id of s2s vpn gateway.

* `p2s_vpn_gateway_id` - The resource id of p2s vpn gateway.

* `express_route_gateway_id` - The resource id of express route gateway.

* `route` - One `route` block defined below.

* `virtual_network_connection` - One or more `virtual_network_connection` block defined below.

* `tags` - Resource tags.

---

The `route` block contains the following:

* `address_prefixes` - List of all addressPrefixes.

* `next_hop_ip_address` - NextHop ip address.

---

The `virtual_network_connection` block contains the following:

* `name` - The name of the resource that is unique within a resource group. This name can be used to access the resource.

* `remote_virtual_network_id` - The resource id of remote virtual network.

* `allow_hub_to_remote_vnet_transit` - VirtualHub to RemoteVnet transit to enabled or not.

* `allow_remote_vnet_to_use_hub_vnet_gateways` - Allow RemoteVnet to use Virtual Hub's gateways.

* `enable_internet_security` - Enable internet security.
