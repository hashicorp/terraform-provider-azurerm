---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_point_to_site_vpn_gateway"
sidebar_current: "docs-azurerm-resource-network-point-to-site-vpn-gateway"
description: |-
  Manages a Point-to-Site VPN Gateway.

---

# azurerm_point_to_site_vpn_gateway

Manages a Point-to-Site VPN Gateway.

## Example Usage

```hcl
resource "azurerm_point_to_site_vpn_gateway" "example" {
  name                        = "example-vpn-gateway"
  location                    = azurerm_resource_group.example.location
  resource_group_name         = azurerm_resource_group.example.resource_group_name
  virtual_hub_id              = azurerm_virtual_hub.example.id
  vpn_server_configuration_id = azurerm_vpn_server_configuration.example.id
  scale_unit                  = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Point-to-Site VPN Gateway. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Point-to-Site VPN Gateway. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `connection_configuration` - (Required) A `connection_configuration` block as defined below.

* `scale_unit` - (Required) The Scale Unit for this Point-to-Site VPN Gateway.

* `virtual_hub_id` - (Required) The ID of the Virtual Hub where this Point-to-Site VPN Gateway should exist. Changing this forces a new resource to be created.

* `vpn_server_configuration_id` - (Required) The ID of the VPN Server Configuration which this Point-to-Site VPN Gateway should use. Changing this forces a new resource to be created. 

* `tags` - (Optional) A mapping of tags to assign to the Point-to-Site VPN Gateway.

---

A `connection_configuration` block supports the following:

* `name` - (Required) The Name which should be used for this Connection Configuration.

* `vpn_client_address_pool` - (Required) A `vpn_client_address_pool` block as defined below.

---

A `vpn_client_address_pool` block supports the following:

* `address_prefixes` - (Required) A list of CIDR Ranges which should be used as Address Prefixes.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Point-to-Site VPN Gateway.

## Import

Point-to-Site VPN Gateway's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_point_to_site_vpn_gateway.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/p2svpnGateways/gateway1
```

