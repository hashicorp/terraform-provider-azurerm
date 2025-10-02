---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_ipam_pool_static_cidr"
description: |-
  Manages a Network Manager IPAM Pool Static CIDR.
---

# azurerm_network_manager_ipam_pool_static_cidr

Manages a Network Manager IPAM Pool Static CIDR.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_subscription" "current" {}

resource "azurerm_network_manager" "example" {
  name                = "example-nm"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
}

resource "azurerm_network_manager_ipam_pool" "example" {
  name               = "example-ipampool"
  network_manager_id = azurerm_network_manager.example.id
  location           = azurerm_resource_group.example.location
  display_name       = "ipampool1"
  address_prefixes   = ["10.0.0.0/24"]
}

resource "azurerm_network_manager_ipam_pool_static_cidr" "example" {
  name             = "example-ipsc"
  ipam_pool_id     = azurerm_network_manager_ipam_pool.example.id
  address_prefixes = ["10.0.0.0/26", "10.0.0.128/27"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Network Manager IPAM Pool Static CIDR. Changing this forces a new Network Manager IPAM Pool Static CIDR to be created.

* `ipam_pool_id` - (Required) The ID of the Network Manager IP Address Management (IPAM) Pool. Changing this forces a new Network Manager IPAM Pool Static CIDR to be created.

---

* `address_prefixes` - (Optional) Specifies a list of IPv4 or IPv6 IP address prefixes which will be allocated to the Static CIDR.

-> **Note:** Exactly one of `address_prefixes` or `number_of_ip_addresses_to_allocate` must be specified.

* `number_of_ip_addresses_to_allocate` - (Optional) The number of IP addresses to allocate to the Static CIDR. The value must be a string representing a positive integer which is a positive power of 2, e.g., `"16"`.

-> **Note:** Exactly one of `address_prefixes` or `number_of_ip_addresses_to_allocate` must be specified.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager IPAM Pool Static CIDR.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Manager IPAM Pool Static CIDR.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager IPAM Pool Static CIDR.
* `update` - (Defaults to 30 minutes) Used when updating the Network Manager IPAM Pool Static CIDR.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Manager IPAM Pool Static CIDR.

## Import

Network Manager IPAM Pool Static CIDRs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_ipam_pool_static_cidr.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/networkManagers/manager1/ipamPools/pool1/staticCidrs/cidr1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network` - 2024-05-01
