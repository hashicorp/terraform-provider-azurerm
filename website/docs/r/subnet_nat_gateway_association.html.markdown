---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subnet_nat_gateway_association"
sidebar_current: "docs-azurerm-resource-network-subnet-nat-gateway-association"
description: |-
  Associates a [NAT Gateway](nat_gateway.html) with a [Subnet](subnet.html) within a [Virtual Network](virtual_network.html).
---

# azurerm_subnet_route_table_association

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
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  virtual_network_name = "${azurerm_virtual_network.example.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_nat_gateway" "example" {
  name                = "example-natgateway"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_subnet_nat_gateway_association" "example" {
  subnet_id      = "${azurerm_subnet.example.id}"
  nat_gateway_id = "${azurerm_nat_gateway.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `nat_gateway_id` - (Required) The Azure resource ID of the NAT Gateway which should be associated with the Subnet. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The Azure resource ID of the Subnet. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The Azure resource ID of the Subnet.

## Import

Subnet NAT Gateway Associations can be imported using the `resource id` of the Subnet, e.g.

```shell
terraform import azurerm_subnet_nat_gateway_association.association1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1/subnets/mysubnet1
```
