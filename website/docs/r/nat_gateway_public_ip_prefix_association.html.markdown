---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nat_gateway_public_ip_prefix_association"
description: |-
  Manages the association between a NAT Gateway and a Public IP Prefix.

---

# azurerm_nat_gateway_public_ip_prefix_association

Manages the association between a NAT Gateway and a Public IP Prefix.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_public_ip_prefix" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  prefix_length       = 30
  zones               = ["1"]
}

resource "azurerm_nat_gateway" "example" {
  name                = "example-NatGateway"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard"
}

resource "azurerm_nat_gateway_public_ip_prefix_association" "example" {
  nat_gateway_id      = azurerm_nat_gateway.example.id
  public_ip_prefix_id = azurerm_public_ip_prefix.example.id
}
```

## Argument Reference

The following arguments are supported:

* `nat_gateway_id` - (Required) The ID of the NAT Gateway. Changing this forces a new resource to be created.

* `public_ip_prefix_id` - (Required) The ID of the Public IP Prefix which this NAT Gateway which should be connected to. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The (Terraform specific) ID of the Association between the NAT Gateway and the Public IP Prefix.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the association between the NAT Gateway and the Public IP Prefix.
* `read` - (Defaults to 5 minutes) Used when retrieving the association between the NAT Gateway and the Public IP Prefix.
* `delete` - (Defaults to 30 minutes) Used when deleting the association between the NAT Gateway and the Public IP Prefix.

## Import

Associations between NAT Gateway and Public IP Prefixes can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_nat_gateway_public_ip_prefix_association.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/natGateways/gateway1|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/publicIPPrefixes/myPublicIpPrefix1"
```

-> **Note:** This is a Terraform Specific ID in the format `{natGatewayID}|{publicIPPrefixID}`

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network`: 2024-05-01
