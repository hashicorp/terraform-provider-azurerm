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

* `subnet_ids` - An array of IDs of the Subnets using this NAT Gateway resource.

## Import

NAT Gateway can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_nat_gateway.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/natGateways/ng1
```
