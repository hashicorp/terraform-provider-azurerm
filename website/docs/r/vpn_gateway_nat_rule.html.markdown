---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_vpn_gateway_nat_rule"
description: |-
  Manages a VPN Gateway Nat Rule.
---

# azurerm_vpn_gateway_nat_rule

Manages a VPN Gateway Nat Rule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_wan" "example" {
  name                = "example-vwan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_hub" "example" {
  name                = "example-vhub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_prefix      = "10.0.1.0/24"
  virtual_wan_id      = azurerm_virtual_wan.example.id
}

resource "azurerm_vpn_gateway" "example" {
  name                = "example-vpngateway"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  virtual_hub_id      = azurerm_virtual_hub.example.id
}

resource "azurerm_vpn_gateway_nat_rule" "example" {
  name                            = "example-vpngatewaynatrule"
  resource_group_name             = azurerm_resource_group.example.name
  vpn_gateway_id                  = azurerm_vpn_gateway.example.id
  external_address_space_mappings = ["192.168.21.0/26"]
  internal_address_space_mappings = ["10.4.0.0/26"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this VPN Gateway Nat Rule. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The Name of the Resource Group in which this VPN Gateway Nat Rule should be created. Changing this forces a new resource to be created.

* `vpn_gateway_id` - (Required) The ID of the VPN Gateway that this VPN Gateway Nat Rule belongs to. Changing this forces a new resource to be created.

* `external_address_space_mappings` - (Required) A list of CIDR Ranges which are used for external mapping of the VPN Gateway Nat Rule.

* `internal_address_space_mappings` - (Required) A list of CIDR Ranges which are used for internal mapping of the VPN Gateway Nat Rule.

* `ip_configuration_id` - (Optional) The ID of the IP Configuration this VPN Gateway Nat Rule applies to. Possible values are `Instance0` and `Instance1`.

* `mode` - (Optional) The source Nat direction of the VPN Nat. Possible values are `EgressSnat` and `IngressSnat`. Defaults to `EgressSnat`. Changing this forces a new resource to be created.

* `type` - (Optional) The type of the VPN Gateway Nat Rule. Possible values are `Dynamic` and `Static`. Defaults to `Static`. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the VPN Gateway Nat Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the VPN Gateway Nat Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the VPN Gateway Nat Rule.
* `update` - (Defaults to 30 minutes) Used when updating the VPN Gateway Nat Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the VPN Gateway Nat Rule.

## Import

VPN Gateway Nat Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_vpn_gateway_nat_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Network/vpnGateways/vpnGateway1/natRules/natRule1
```
