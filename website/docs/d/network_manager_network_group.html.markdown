---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_network_group"
description: |-
  Get information about a Network Manager Network Group.
---

# azurerm_network_manager_network_group

Use this data source to access information about a Network Manager Network Group.

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

data "azurerm_network_manager_network_group" "example" {
  name               = azurerm_network_manager_network_group.example.name
  network_manager_id = azurerm_network_manager.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Network Manager Network Group.

* `network_manager_id` - (Required) Specifies the ID of the Network Manager.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager Network Group.

* `description` - A description of the Network Manager Network Group.
 
## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager Network Group.
