---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_commit"
description: |-
  Manages a Network Manager Commit.
---

# azurerm_network_manager_commit

Manages a Network Manager Commit.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_subscription" "current" {
}

resource "azurerm_network_manager" "example" {
  name                = "example-network-manager"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Connectivity", "SecurityAdmin"]
  description    = "example network manager"
}

resource "azurerm_network_manager_network_group" "example" {
  name               = "example-group"
  network_manager_id = azurerm_network_manager.example.id
}

resource "azurerm_virtual_network" "example" {
  name                    = "example-net"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  address_space           = ["10.0.0.0/16"]
  flow_timeout_in_minutes = 10
}

resource "azurerm_network_manager_connectivity_configuration" "example" {
  name                  = "example-connectivity-conf"
  network_manager_id    = azurerm_network_manager.example.id
  connectivity_topology = "HubAndSpoke"
  applies_to_group {
    group_connectivity = "None"
    network_group_id   = azurerm_network_manager_network_group.example.id
  }
  hub {
    resource_id   = azurerm_virtual_network.example.id
    resource_type = "Microsoft.Network/virtualNetworks"
  }
}

resource "azurerm_network_manager_commit" "example" {
  network_manager_id = azurerm_network_manager.example.id
  location           = "eastus"
  scope_access       = "Connectivity"
  configuration_ids  = [azurerm_network_manager_connectivity_configuration.example.id]
}
```

## example usage (Triggers)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_subscription" "current" {
}

resource "azurerm_network_manager" "example" {
  name                = "example-network-manager"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Connectivity", "SecurityAdmin"]
  description    = "example network manager"
}

resource "azurerm_network_manager_network_group" "example" {
  name               = "example-group"
  network_manager_id = azurerm_network_manager.example.id
}

resource "azurerm_virtual_network" "example" {
  name                    = "example-net"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  address_space           = ["10.0.0.0/16"]
  flow_timeout_in_minutes = 10
}

resource "azurerm_network_manager_security_admin_configuration" "example" {
  name               = "example-nmsac"
  network_manager_id = azurerm_network_manager.example.id
}

resource "azurerm_network_manager_admin_rule_collection" "example" {
  name                            = "example-nmarc"
  security_admin_configuration_id = azurerm_network_manager_security_admin_configuration.example.id
  network_group_ids               = [azurerm_network_manager_network_group.example.id]
}

resource "azurerm_network_manager_admin_rule" "example" {
  name                     = "example-nmar"
  admin_rule_collection_id = azurerm_network_manager_admin_rule_collection.example.id
  action                   = "Deny"
  description              = "example"
  direction                = "Inbound"
  priority                 = 1
  protocol                 = "Tcp"
  source_port_ranges       = ["80"]
  destination_port_ranges  = ["80"]
  source {
    address_prefix_type = "ServiceTag"
    address_prefix      = "Internet"
  }
  destination {
    address_prefix_type = "IPPrefix"
    address_prefix      = "*"
  }
}

resource "azurerm_network_manager_commit" "example" {
  network_manager_id = azurerm_network_manager.example.id
  location           = "eastus"
  scope_access       = "SecurityAdmin"
  configuration_ids  = [azurerm_network_manager_security_admin_configuration.example.id]
  depends_on         = [azurerm_network_manager_admin_rule.example]
  triggers = {
    source_port_ranges = join(",", azurerm_network_manager_admin_rule.example.source_port_ranges)
  }
}
```

## Arguments Reference

The following arguments are supported:

* `network_manager_id` - (Required) Specifies the ID of the Network Manager. Changing this forces a new Network Manager Commit to be created.

* `location` - (Required) Specifies the location which the configurations will be deployed to. Changing this forces a new Network Manager Commit to be created.

* `scope_access` - (Required) Specifies the configuration deployment type. Possible values are `Connectivity` and `SecurityAdmin`.

* `configuration_ids` - (Required) A list of Network Manager Configuration IDs which should be aligned with `scope_access`.

* `triggers` - (Optional) A mapping of key values pairs that can be used to keep the deployment up with the Network Manager configurations and rules.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager Admin Rule Collection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 hours) Used when creating the Network Manager Commit.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager Commit.
* `update` - (Defaults to 2 hours) Used when updating the Network Manager Commit.
* `delete` - (Defaults to 1 hour) Used when deleting the Network Manager Commit.

## Import

Network Manager Commit can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_commit.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/networkManagers/networkManager1/commit|eastus|Connectivity
```
