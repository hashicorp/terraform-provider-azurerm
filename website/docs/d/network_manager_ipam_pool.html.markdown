---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_network_manager_ipam_pool"
description: |-
  Gets information about an existing Network Manager IPAM Pool.
---

# Data Source: azurerm_network_manager_ipam_pool

Use this data source to access information about an existing Network Manager IPAM Pool.

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
}

resource "azurerm_network_manager_ipam_pool" "example" {
  name               = "example-ipam-pool"
  location           = "West Europe"
  network_manager_id = azurerm_network_manager.example.id
  display_name       = "example-pool"
  address_prefixes   = ["10.0.0.0/24"]
}

data "azurerm_network_manager_ipam_pool" "example" {
  name               = azurerm_network_manager_ipam_pool.example.name
  network_manager_id = azurerm_network_manager.example.id
}

output "id" {
  value = data.azurerm_network_manager_ipam_pool.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Network Manager IPAM Pool.

* `network_manager_id` - (Required) The ID of the parent Network Manager.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager IPAM Pool.

* `address_prefixes` - A list of IPv4 or IPv6 IP address prefixes assigned to the Network Manager IPAM Pool.

* `description` - The description of the Network Manager IPAM Pool.

* `display_name` - The display name of the Network Manager IPAM Pool.

* `location` - The Azure Region where the Network Manager IPAM Pool exists.

* `parent_pool_name` - The name of the parent IPAM Pool.

* `tags` - A mapping of tags assigned to the Network Manager IPAM Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager IPAM Pool.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Network`: 2024-05-01
