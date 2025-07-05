---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_routing_configuration"
description: |-
  Manages a Network Manager Routing Configuration.
---

# azurerm_network_manager_routing_configuration

Manages a Network Manager Routing Configuration.

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
  scope_accesses = ["Routing"]
}

resource "azurerm_network_manager_routing_configuration" "example" {
  name               = "example-routing-configuration"
  network_manager_id = azurerm_network_manager.example.id
  description        = "example routing configuration"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Network Manager Routing Configuration. Changing this forces a new Network Manager Routing Configuration to be created.

* `network_manager_id` - (Required) The ID of the Network Manager. Changing this forces a new Network Manager Routing Configuration to be created.

---

* `description` - (Optional) The description of the Network Manager.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager Routing Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Manager Routing Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager Routing Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the Network Manager Routing Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Manager Routing Configuration.

## Import

Network Manager Routing Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_routing_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/networkManagers/manager1/routingConfigurations/conf1
```


## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network`: 2024-05-01
