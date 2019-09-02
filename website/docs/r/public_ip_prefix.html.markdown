---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_public_ip_prefix"
sidebar_current: "docs-azurerm-resource-network-public-ip-prefix"
description: |-
  Manages a Public IP Prefix.
---

# azurerm_public_ip_prefix

Manages a Public IP Prefix.

-> **NOTE** Public IP Prefix are currently in Public Preview. You can find more information about [Public IP Preifx Preview here](https://docs.microsoft.com/en-us/azure/virtual-network/public-ip-address-prefix).

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "resourceGroup1"
  location = "West Europe"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acceptanceTestPublicIpPrefix1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  prefix_length = 31

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Public IP resource . Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the public IP.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Optional) The SKU of the Public IP Prefix. Accepted values are `Standard`. Defaults to `Standard`. Changing this forces a new resource to be created.

-> **Note** Public IP Prefix can only be created with Standard SKUs at this time.

* `prefix_length` - (Optional) Specifies the number of bits of the prefix. The value can be set between 24 (256 addresses) and 31 (2 addresses). Changing this forces a new resource to be created.

-> **Please Note:**: There may be Public IP address limits on the subscription . [More information available here](https://docs.microsoft.com/en-us/azure/azure-subscription-service-limits?toc=%2fazure%2fvirtual-network%2ftoc.json#publicip-address)

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `zones` - (Optional) A collection containing the availability zone to allocate the Public IP in.

-> **Please Note**: Availability Zones are [only supported in several regions at this time](https://docs.microsoft.com/en-us/azure/availability-zones/az-overview).

## Attributes Reference

The following attributes are exported:

* `id` - The Public IP Prefix ID.
* `ip_prefix` - The IP address prefix value that was allocated.

## Import

Public IP Prefixes can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_public_ip_prefix.myPublicIpPrefix /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/publicIPFixes/myPublicIpPrefix1
```
