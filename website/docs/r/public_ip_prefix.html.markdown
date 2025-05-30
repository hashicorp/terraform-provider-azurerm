---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_public_ip_prefix"
description: |-
  Manages a Public IP Prefix.
---

# azurerm_public_ip_prefix

Manages a Public IP Prefix.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_public_ip_prefix" "example" {
  name                = "acceptanceTestPublicIpPrefix1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  prefix_length = 31

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Public IP Prefix resource . Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Public IP Prefix. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Optional) The SKU of the Public IP Prefix. Accepted values are `Standard`. Defaults to `Standard`. Changing this forces a new resource to be created.

-> **Note:** Public IP Prefix can only be created with Standard SKUs at this time.

* `sku_tier` - (Optional) The SKU Tier that should be used for the Public IP. Possible values are `Regional` and `Global`. Defaults to `Regional`. Changing this forces a new resource to be created.

* `ip_version` - (Optional) The IP Version to use, `IPv6` or `IPv4`. Changing this forces a new resource to be created. Default is `IPv4`.

* `prefix_length` - (Optional) Specifies the number of bits of the prefix. The value can be set between 0 (4,294,967,296 addresses) and 31 (2 addresses). Defaults to `28`(16 addresses). Changing this forces a new resource to be created.

-> **Note:** There may be Public IP address limits on the subscription . [More information available here](https://docs.microsoft.com/azure/azure-subscription-service-limits?toc=%2fazure%2fvirtual-network%2ftoc.json#publicip-address)

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `zones` - (Optional) Specifies a list of Availability Zones in which this Public IP Prefix should be located. Changing this forces a new Public IP Prefix to be created.

-> **Note:** Availability Zones are [only supported in several regions at this time](https://docs.microsoft.com/azure/availability-zones/az-overview).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Public IP Prefix ID.
* `ip_prefix` - The IP address prefix value that was allocated.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Public IP Prefix.
* `read` - (Defaults to 5 minutes) Used when retrieving the Public IP Prefix.
* `update` - (Defaults to 30 minutes) Used when updating the Public IP Prefix.
* `delete` - (Defaults to 30 minutes) Used when deleting the Public IP Prefix.

## Import

Public IP Prefixes can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_public_ip_prefix.myPublicIpPrefix /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/publicIPPrefixes/myPublicIpPrefix1
```
