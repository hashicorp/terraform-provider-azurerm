---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nat_gateway"
sidebar_current: "docs-azurerm-resource-nat-gateway"
description: |-
  Manages a Azure NAT Gateway.
---
# azurerm_nat_gateway

Manages a Azure NAT Gateway.

-> **NOTE:** The Azure NAT Gateway service is currently in private preview. Your subscription must be on the NAT Gateway private preview whitelist for this resource to be provisioned correctly. If you attempt to provision this resource and receive an `InvalidResourceType` error may mean that your subscription is not part of the NAT Gateway private preview or you are using a region which does not yet support the NAT Gateway private preview service. The NAT Gateway private preview service is currently available in a limited set of regions. Private preview resources may have multiple breaking changes over their lifecycle until they GA. You can opt into the Private Preview by contacting your Microsoft Representative.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "nat-gateway-example-rg"
  location = "eastus2"
}

resource "azurerm_public_ip" "example" {
  name                = "nat-gateway-publicIP"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_public_ip_prefix" "example" {
  name                = "nat-gateway-publicIPPrefix"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  prefix_length       = 30
  zones               = ["1"]
}

resource "azurerm_nat_gateway" "example" {
  name                    = "nat-Gateway"
  location                = "${azurerm_resource_group.example.location}"
  resource_group_name     = "${azurerm_resource_group.example.name}"
  public_ip_address_ids   = ["${azurerm_public_ip.example.id}"]
  public_ip_prefix_ids    = ["${azurerm_public_ip_prefix.example.id}"]
  sku_name                = "Standard"
  idle_timeout_in_minutes = 10
  zones                   = ["1"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the NAT Gateway. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which the NAT Gateway should exist. Changing this forces a new resource to be created.

* `location` - (Optional) Specifies the supported Azure location where the NAT Gateway should exist. Changing this forces a new resource to be created.

* `idle_timeout_in_minutes` - (Optional) The idle timeout which should be used in minutes. Defaults to `4`.

* `public_ip_address_ids` - (Optional) A list of Public IP Address ID's which should be associated with the NAT Gateway resource.

* `public_ip_prefix_ids` - (Optional) A list of Public IP Prefix ID's which should be associated with the NAT Gateway resource.

* `sku_name` - (Optional) The SKU which should be used. At this time the only supported value is `Standard`. Defaults to `Standard`.

* `tags` - (Optional) A mapping of tags to assign to the resource. Changing this forces a new resource to be created.

* `zones` - (Optional) A list of availability zones where the NAT Gateway should be provisioned. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the NAT Gateway.

* `resource_guid` - The resource GUID property of the NAT Gateway.

## Import

NAT Gateway can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_nat_gateway.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/natGateways/gateway1
```
