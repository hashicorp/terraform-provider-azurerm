---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nat_gateway"
description: |-
  Manages a Azure NAT Gateway.
---
# azurerm_nat_gateway

Manages an Azure NAT Gateway.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "nat-gateway-example-rg"
  location = "West Europe"
}

resource "azurerm_nat_gateway" "example" {
  name                    = "nat-gateway"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  sku_name                = "Standard"
  idle_timeout_in_minutes = 10
  zones                   = ["1"]
}
```

For more complete examples, please see the [azurerm_nat_gateway_public_ip_association](nat_gateway_public_ip_association.html) and [azurerm_nat_gateway_public_ip_prefix_association](nat_gateway_public_ip_prefix_association.html) resources.

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the NAT Gateway. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which the NAT Gateway should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the NAT Gateway should exist. Changing this forces a new resource to be created.

* `idle_timeout_in_minutes` - (Optional) The idle timeout which should be used in minutes. Defaults to `4`.

* `sku_name` - (Optional) The SKU which should be used. At this time the only supported value is `Standard`. Defaults to `Standard`.

* `tags` - (Optional) A mapping of tags to assign to the resource. 

* `zones` - (Optional) A list of Availability Zones in which this NAT Gateway should be located. Changing this forces a new NAT Gateway to be created.

-> **Note:** Only one Availability Zone can be defined. For more information, please check out the [Azure documentation](https://learn.microsoft.com/en-us/azure/nat-gateway/nat-overview#availability-zones)

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the NAT Gateway.

* `resource_guid` - The resource GUID property of the NAT Gateway.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the NAT Gateway.
* `read` - (Defaults to 5 minutes) Used when retrieving the NAT Gateway.
* `update` - (Defaults to 1 hour) Used when updating the NAT Gateway.
* `delete` - (Defaults to 1 hour) Used when deleting the NAT Gateway.

## Import

NAT Gateway can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_nat_gateway.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/natGateways/gateway1
```
