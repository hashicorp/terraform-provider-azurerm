---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_public_ip"
description: |-
  Manages a Public IP Address.
---

# azurerm_public_ip

Manages a Public IP Address.

~> **Note:** If this resource is to be associated with a resource that requires disassociation before destruction (such as `azurerm_network_interface`) it is recommended to set the `lifecycle` argument `create_before_destroy = true`. Otherwise, it can fail to disassociate on destruction.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_public_ip" "example" {
  name                = "acceptanceTestPublicIp1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  allocation_method   = "Static"

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Public IP. Changing this forces a new Public IP to be created.

* `resource_group_name` - (Required) The name of the Resource Group where this Public IP should exist. Changing this forces a new Public IP to be created.

* `location` - (Required) Specifies the supported Azure location where the Public IP should exist. Changing this forces a new resource to be created.

* `allocation_method` - (Required) Defines the allocation method for this IP address. Possible values are `Static` or `Dynamic`.

~> **Note:** `Dynamic` Public IP Addresses aren't allocated until they're assigned to a resource (such as a Virtual Machine or a Load Balancer) by design within Azure. See `ip_address` argument.

---

* `zones` - (Optional) A collection containing the availability zone to allocate the Public IP in. Changing this forces a new resource to be created.

-> **Note:** Availability Zones are only supported with a [Standard SKU](https://docs.microsoft.com/azure/virtual-network/virtual-network-ip-addresses-overview-arm#standard) and [in select regions](https://docs.microsoft.com/azure/availability-zones/az-overview) at this time. Standard SKU Public IP Addresses that do not specify a zone are **not** zone-redundant by default.

* `ddos_protection_mode` - (Optional) The DDoS protection mode of the public IP. Possible values are `Disabled`, `Enabled`, and `VirtualNetworkInherited`. Defaults to `VirtualNetworkInherited`.
 
* `ddos_protection_plan_id` - (Optional) The ID of DDoS protection plan associated with the public IP. 

-> **Note:** `ddos_protection_plan_id` can only be set when `ddos_protection_mode` is `Enabled`.

* `domain_name_label` - (Optional) Label for the Domain Name. Will be used to make up the FQDN. If a domain name label is specified, an A DNS record is created for the public IP in the Microsoft Azure DNS system.

* `domain_name_label_scope` - (Optional) Scope for the domain name label. If a domain name label scope is specified, an A DNS record is created for the public IP in the Microsoft Azure DNS system with a hashed value includes in FQDN. Possible values are `NoReuse`, `ResourceGroupReuse`, `SubscriptionReuse` and `TenantReuse`. Changing this forces a new Public IP to be created.

* `edge_zone` - (Optional) Specifies the Edge Zone within the Azure Region where this Public IP should exist. Changing this forces a new Public IP to be created.

* `idle_timeout_in_minutes` - (Optional) Specifies the timeout for the TCP idle connection. The value can be set between 4 and 30 minutes.

* `ip_tags` - (Optional) A mapping of IP tags to assign to the public IP. Changing this forces a new resource to be created.

-> **Note:** IP Tag `RoutingPreference` requires multiple `zones` and `Standard` SKU to be set.

* `ip_version` - (Optional) The IP Version to use, IPv6 or IPv4. Changing this forces a new resource to be created. Defaults to `IPv4`.

-> **Note:** Only `static` IP address allocation is supported for IPv6.

* `public_ip_prefix_id` - (Optional) If specified then public IP address allocated will be provided from the public IP prefix resource. Changing this forces a new resource to be created.

* `reverse_fqdn` - (Optional) A fully qualified domain name that resolves to this public IP address. If the reverseFqdn is specified, then a PTR DNS record is created pointing from the IP address in the in-addr.arpa domain to the reverse FQDN.

* `sku` - (Optional) The SKU of the Public IP. Accepted values are `Basic` and `Standard`. Defaults to `Standard`. Changing this forces a new resource to be created.

-> **Note:** Public IP Standard SKUs require `allocation_method` to be set to `Static`.

* `sku_tier` - (Optional) The SKU Tier that should be used for the Public IP. Possible values are `Regional` and `Global`. Defaults to `Regional`. Changing this forces a new resource to be created.

-> **Note:** When `sku_tier` is set to `Global`, `sku` must be set to `Standard`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of this Public IP.

* `ip_address` - The IP address value that was allocated.

~> **Note:** `Dynamic` Public IP Addresses aren't allocated until they're attached to a device (e.g. a Virtual Machine/Load Balancer). Instead you can obtain the IP Address once the Public IP has been assigned via the [`azurerm_public_ip` Data Source](../d/public_ip.html).

* `fqdn` - Fully qualified domain name of the A DNS record associated with the public IP. `domain_name_label` must be specified to get the `fqdn`. This is the concatenation of the `domain_name_label` and the regionalized DNS zone

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Public IP.
* `read` - (Defaults to 5 minutes) Used when retrieving the Public IP.
* `update` - (Defaults to 30 minutes) Used when updating the Public IP.
* `delete` - (Defaults to 30 minutes) Used when deleting the Public IP.

## Import

Public IPs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_public_ip.myPublicIp /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/publicIPAddresses/myPublicIpAddress1
```
