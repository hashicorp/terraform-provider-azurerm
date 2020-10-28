---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subnet_nat_gateway_association"
description: |-
  Associates a [NAT Gateway](nat_gateway.html) with a [Subnet](subnet.html) within a [Virtual Network](virtual_network.html).
---

# azurerm_subnet_nat_gateway_association

Associates a [NAT Gateway](nat_gateway.html) with a [Subnet](subnet.html) within a [Virtual Network](virtual_network.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-nat-gateway-rg"
  location = "East US 2"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name               = "example-subnet"
  virtual_network_id = azurerm_virtual_network.example.id
  address_prefixes   = ["10.0.2.0/24"]
}

resource "azurerm_nat_gateway" "example" {
  name                = "example-natgateway"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet_nat_gateway_association" "example" {
  subnet_id      = azurerm_subnet.example.id
  nat_gateway_id = azurerm_nat_gateway.example.id
}
```

## Argument Reference

The following arguments are supported:

* `nat_gateway_id` - (Required) The ID of the NAT Gateway which should be associated with the Subnet. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the Subnet. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Subnet.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Subnet NAT Gateway Association.
* `update` - (Defaults to 30 minutes) Used when updating the Subnet NAT Gateway Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Subnet NAT Gateway Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the Subnet NAT Gateway Association.

## Import

Subnet NAT Gateway Associations can be imported using the `resource id` of the Subnet, e.g.

```shell
terraform import azurerm_subnet_nat_gateway_association.association1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1/subnets/mysubnet1
```
