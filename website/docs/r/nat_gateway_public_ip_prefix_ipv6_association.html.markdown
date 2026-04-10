---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nat_gateway_public_ip_prefix_ipv6_association"
description: |-
  Manages the association between a NAT Gateway and an IPv6 Public IP Prefix.
---

# azurerm_nat_gateway_public_ip_prefix_ipv6_association

Manages the association between a NAT Gateway and an IPv6 Public IP Prefix.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_public_ip_prefix" "example" {
  name                = "example-public-ip-prefix"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  ip_version          = "IPv6"
  prefix_length       = 127
  sku                 = "StandardV2"
}

resource "azurerm_nat_gateway" "example" {
  name                = "example-nat-gateway"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "StandardV2"
}

resource "azurerm_nat_gateway_public_ip_prefix_ipv6_association" "example" {
  nat_gateway_id      = azurerm_nat_gateway.example.id
  public_ip_prefix_id = azurerm_public_ip_prefix.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `nat_gateway_id` - (Required) The ID of the NAT Gateway. Changing this forces a new resource to be created.

~> **Note:** `nat_gateway_id` must reference a NAT Gateway with SKU `StandardV2`.

* `public_ip_prefix_id` - (Required) The ID of the IPv6 Public IP Prefix which this NAT Gateway should be connected to. Changing this forces a new resource to be created.

~> **Note:** `public_ip_prefix_id` must reference an `IPv6` Public IP Prefix with SKU `StandardV2`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the NAT Gateway and IPv6 Public IP Prefix association.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the NAT Gateway and IPv6 Public IP Prefix association.
* `read` - (Defaults to 5 minutes) Used when retrieving the NAT Gateway and IPv6 Public IP Prefix association.
* `delete` - (Defaults to 30 minutes) Used when deleting the NAT Gateway and IPv6 Public IP Prefix association.

## Import

A NAT Gateway and IPv6 Public IP Prefix association can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_nat_gateway_public_ip_prefix_ipv6_association.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/natGateways/natGateway1|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/publicIPPrefixes/publicIPPrefix1"
```

-> **Note:** This is a Terraform-specific ID in the format `{natGatewayID}|{publicIPPrefixID}`.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network` - 2025-01-01