---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_ipam_pool"
description: |-
  Manages a Network Manager IP Address Management (IPAM) Pool.
---

# azurerm_network_manager_ipam_pool

Manages a Network Manager IP Address Management (IPAM) Pool.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_subscription" "current" {}

resource "azurerm_network_manager" "example" {
  name                = "example-network-manager"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Connectivity", "SecurityAdmin"]
}

resource "azurerm_network_manager_ipam_pool" "example" {
  name               = "example-ipam-pool"
  location           = "West Europe"
  network_manager_id = azurerm_network_manager.example.id
  display_name       = "example-pool"
  address_prefixes   = ["10.0.0.0/24"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Network Manager IPAM Pool. Changing this forces a new Network Manager IPAM Pool to be created.

* `network_manager_id` - (Required) The ID of the parent Network Manager. Changing this forces a new Network Manager IPAM Pool to be created.

* `address_prefixes` - (Required) Specifies a list of IPv4 or IPv6 IP address prefixes. Changing this forces a new Network Manager IPAM Pool to be created.

* `display_name` - (Required) The display name for the Network Manager IPAM Pool.

* `location` - (Required) The Azure Region where the Network Manager IPAM Pool should exist. Changing this forces a new Network Manager IPAM Pool to be created.

---

* `description` - (Optional) The description of the Network Manager IPAM Pool.

* `parent_pool_name` - (Optional) The name of the parent IPAM Pool. Changing this forces a new Network Manager IPAM Pool to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Network Manager IPAM Pool.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager IPAM Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Manager IPAM Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager IPAM Pool.
* `update` - (Defaults to 30 minutes) Used when updating the Network Manager IPAM Pool.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Manager IPAM Pool.

## Import

Network Manager IPAM Pools can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_ipam_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/networkManagers/manager1/ipamPools/pool1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network`: 2024-05-01
