---
subcategory: ""
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nat_gateway"
sidebar_current: "docs-azurerm-resource-nat-gateway"
description: |-
  Manages an Azure NAT Gateway instance.
---
# azurerm_nat_gateway

Manages an Azure NAT Gateway instance.

-> **NOTE:** The Azure NAT Gateway service is currently in private preview. Your subscription must be on the NAT Gateway private preview whitelist for this resource to be provisioned correctly. If you attempt to provision this resource and receive an `InvalidResourceType` error that means that your subscription is not part of the NAT Gateway private preview whitelist and you will not be able to use this resource. The NAT Gateway private preview service is currently only available in the `East US 2` and `West Central US` regions. You can opt into the Private Preview by contacting your Microsoft Representative.

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

* `name` - (Required) Specifies the name of the Nat Gateway. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the NAT Gateway exists.

* `location` - (Optional) Specifies the supported Azure location where the resource should exist. Changing this forces a new resource to be created.

* `idle_timeout_in_minutes` - (Optional) The idle timeout of the Nat Gateway. Defaults to `4`.

* `public_ip_address_ids` - (Optional) An array of the IDs of Public IP Addresses associated with the NAT Gateway resource.

* `public_ip_prefix_ids` - (Optional) An array of the IDs of Public IP Prefixes associated with the NAT Gateway resource.

* `sku_name` - (Optional) The nat gateway SKU, supported value, `Standard`. Defaults to `Standard`.

* `zones` - (Optional) A list of availability zones where the Nat Gateway should be provisioned. Supported values are `1`, `2`, and `3`. For more information on `zones` please refer to the [product documentation](https://docs.microsoft.com/en-us/azure/availability-zones/az-overview). Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `resource_guid` - The resource GUID property of the Nat Gateway.

* `subnet_ids` - A list subnet IDs that are using this NAT Gateway resource.

## Import

NAT Gateway can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_nat_gateway.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/natGateways/ng1
```
