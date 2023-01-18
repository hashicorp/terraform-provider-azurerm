---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_scope_connection"
description: |-
  Manages a Network Manager Scope Connection.
---

# azurerm_network_manager_scope_connection

Manages a Network Manager Scope Connection which may cross tenants.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_client_config" "current" {
}

data "azurerm_subscription" "current" {
}

data "azurerm_subscription" "alt" {
  subscription_id = "00000000-0000-0000-0000-000000000000"
}

resource "azurerm_network_manager" "example" {
  name                = "example-networkmanager"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["SecurityAdmin"]
}

resource "azurerm_network_manager_scope_connection" "example" {
  name               = "example-nsc"
  network_manager_id = azurerm_network_manager.example.id
  tenant_id          = data.azurerm_client_config.current.tenant_id
  target_scope_id    = data.azurerm_subscription.alt.id
  description        = "example"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Manager Scope Connection. Changing this forces a new Network Manager Scope Connection to be created.

* `network_manager_id` - (Required) Specifies the ID of the Network Manager Scope Connection. Changing this forces a new Network Manager Scope Connection to be created.

* `target_scope_id` - (Required) Specifies the Resource ID of the target scope which the Network Manager is connected to. It should be either Subscription ID or Management Group ID.

* `tenant_id` - (Required) Specifies the Tenant ID of the Resource which the Network Manager is connected to.

* `description` - (Optional) A description of the Network Manager Scope Connection.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager Scope Connection.

* `connection_state` - The Connection state of the Network Manager Scope Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Manager Scope Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager Scope Connection.
* `update` - (Defaults to 30 minutes) Used when updating the Network Manager Scope Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Manager Scope Connection.

## Import

Network Manager Scope Connection can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_scope_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/networkManagers/networkManager1/scopeConnections/scopeConnection1
```
