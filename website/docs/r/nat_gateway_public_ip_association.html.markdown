---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nat_gateway_public_ip_association"
description: |-
  Manages the association between a Nat Gateway and a Public IP.

---

# azurerm_nat_gateway_public_ip_association

Manages the association between a Nat Gateway and a Public IP.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_public_ip" "example" {
  name                = "example-PIP"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_nat_gateway" "example" {
  name                = "example-NatGateway"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard"
}

resource "azurerm_nat_gateway_public_ip_association" "example" {
  nat_gateway_id       = azurerm_nat_gateway.example.id
  public_ip_address_id = azurerm_public_ip.example.id
}
```

## Argument Reference

The following arguments are supported:

* `nat_gateway_id` - (Required) The ID of the Nat Gateway. Changing this forces a new resource to be created.

* `public_ip_address_id` - (Required) The ID of the Public IP which this Nat Gateway which should be connected to. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The (Terraform specific) ID of the Association between the Nat Gateway and the Public IP.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the association between the Nat Gateway and the Public IP.
* `read` - (Defaults to 5 minutes) Used when retrieving the association between the Nat Gateway and the Public IP.
* `delete` - (Defaults to 30 minutes) Used when deleting the association between the Nat Gateway and the Public IP.

## Import

Associations between Nat Gateway and Public IP Addresses can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_nat_gateway_public_ip_association.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/natGateways/gateway1|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/publicIPAddresses/myPublicIpAddress1"
```

-> **Note:** This is a Terraform Specific ID in the format `{natGatewayID}|{publicIPAddressID}`
