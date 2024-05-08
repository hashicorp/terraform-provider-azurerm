---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_network_gateway_nat_rule"
description: |-
  Manages a Virtual Network Gateway Nat Rule.
---

# azurerm_virtual_network_gateway_nat_rule

Manages a Virtual Network Gateway Nat Rule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "example" {
  name                = "example-pip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "example" {
  name                = "example-vnetgw"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.example.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.example.id
  }
}

data "azurerm_virtual_network_gateway" "example" {
  name                = azurerm_virtual_network_gateway.example.name
  resource_group_name = azurerm_virtual_network_gateway.example.resource_group_name
}

resource "azurerm_virtual_network_gateway_nat_rule" "example" {
  name                       = "example-vnetgwnatrule"
  resource_group_name        = azurerm_resource_group.example.name
  virtual_network_gateway_id = data.azurerm_virtual_network_gateway.example.id
  mode                       = "EgressSnat"
  type                       = "Dynamic"
  ip_configuration_id        = data.azurerm_virtual_network_gateway.example.ip_configuration[0].id

  external_mapping {
    address_space = "10.2.0.0/26"
    port_range    = "200"
  }

  internal_mapping {
    address_space = "10.4.0.0/26"
    port_range    = "400"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Virtual Network Gateway Nat Rule. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The Name of the Resource Group in which this Virtual Network Gateway Nat Rule should be created. Changing this forces a new resource to be created.

* `virtual_network_gateway_id` - (Required) The ID of the Virtual Network Gateway that this Virtual Network Gateway Nat Rule belongs to. Changing this forces a new resource to be created.

* `external_mapping` - (Required) One or more `external_mapping` blocks as documented below.

* `internal_mapping` - (Required) One or more `internal_mapping` blocks as documented below.

* `ip_configuration_id` - (Optional) The ID of the IP Configuration this Virtual Network Gateway Nat Rule applies to.

* `mode` - (Optional) The source Nat direction of the Virtual Network Gateway Nat. Possible values are `EgressSnat` and `IngressSnat`. Defaults to `EgressSnat`. Changing this forces a new resource to be created.

* `type` - (Optional) The type of the Virtual Network Gateway Nat Rule. Possible values are `Dynamic` and `Static`. Defaults to `Static`. Changing this forces a new resource to be created.

---

A `external_mapping` block exports the following:

* `address_space` - (Required) The string CIDR representing the address space for the Virtual Network Gateway Nat Rule external mapping.

* `port_range` - (Optional) The single port range for the Virtual Network Gateway Nat Rule external mapping.

---

A `internal_mapping` block exports the following:

* `address_space` - (Required) The string CIDR representing the address space for the Virtual Network Gateway Nat Rule internal mapping.

* `port_range` - (Optional) The single port range for the Virtual Network Gateway Nat Rule internal mapping.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Network Gateway Nat Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Network Gateway Nat Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Network Gateway Nat Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Network Gateway Nat Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Network Gateway Nat Rule.

## Import

Virtual Network Gateway Nat Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_network_gateway_nat_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworkGateways/gw1/natRules/rule1
```
