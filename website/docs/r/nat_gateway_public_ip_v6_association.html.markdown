---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nat_gateway_public_ip_v6_association"
description: |-
  Manages the association between a NAT Gateway and an IPv6 Public IP Address.
---

# azurerm_nat_gateway_public_ip_v6_association

Manages the association between a NAT Gateway and an IPv6 Public IP Address.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_public_ip" "example" {
  name                = "example-PIPv6"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "StandardV2"
  ip_version          = "IPv6"
}

resource "azurerm_nat_gateway" "example" {
  name                = "example-NatGateway"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "StandardV2"
  zones               = [1, 2, 3]
}

resource "azurerm_nat_gateway_public_ip_v6_association" "example" {
  nat_gateway_id       = azurerm_nat_gateway.example.id
  public_ip_address_id = azurerm_public_ip.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `nat_gateway_id` - (Required) The ID of the NAT Gateway. Changing this forces a new resource to be created.

* `public_ip_address_id` - (Required) The ID of the IPv6 Public IP which this NAT Gateway should be connected to. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The (Terraform specific) ID of the Association between the NAT Gateway and the IPv6 Public IP.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the association between the NAT Gateway and the IPv6 Public IP.
* `read` - (Defaults to 5 minutes) Used when retrieving the association between the NAT Gateway and the IPv6 Public IP.
* `delete` - (Defaults to 30 minutes) Used when deleting the association between the NAT Gateway and the IPv6 Public IP.

## Import

Associations between NAT Gateway and IPv6 Public IP Addresses can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_nat_gateway_public_ip_v6_association.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/natGateways/gateway1|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/publicIPAddresses/myPublicIpAddressV6"
```

-> **Note:** This is a Terraform Specific ID in the format `{natGatewayID}|{publicIPAddressID}`

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network` - 2025-01-01
