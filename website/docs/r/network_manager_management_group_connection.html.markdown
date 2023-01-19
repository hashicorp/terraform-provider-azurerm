---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_management_group_connection"
description: |-
  Manages a Network Manager Management Group Connection.
---

# azurerm_network_manager_management_group_connection

Manages a Network Manager Management Group Connection which may cross tenants.

## Example Usage

```hcl
resource "azurerm_management_group" "example" {
}

resource "azurerm_management_group_subscription_association" "example" {
  management_group_id = azurerm_management_group.example.id
  subscription_id     = data.azurerm_subscription.alt.id
}

data "azurerm_subscription" "alt" {
  subscription_id = "00000000-0000-0000-0000-000000000000"
}

data "azurerm_subscription" "current" {
}

data "azurerm_client_config" "current" {
}

resource "azurerm_role_assignment" "network_contributor" {
  scope                = azurerm_management_group.example.id
  role_definition_name = "Network Contributor"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
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

resource "azurerm_network_manager_management_group_connection" "example" {
  name                = "example-nmmgc"
  management_group_id = azurerm_management_group.example.id
  network_manager_id  = azurerm_network_manager.example.id
  description         = "example"
  depends_on          = [azurerm_role_assignment.network_contributor]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Manager Management Group Connection. Changing this forces a new Network Manager Management Group Connection to be created.

* `management_group_id` - (Required) Specifies the ID of the target Management Group. Changing this forces a new resource to be created.

* `network_manager_id` - (Required) Specifies the ID of the Network Manager which the Management Group is connected to. Changing this forces a new resource to be created.

* `description` - (Optional) A description of the Network Manager Management Group Connection.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager Management Group Connection.

* `connection_state` - The Connection state of the Network Manager Management Group Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Manager Management Group Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager Management Group Connection.
* `update` - (Defaults to 30 minutes) Used when updating the Network Manager Management Group Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Manager Management Group Connection.

## Import

Network Manager Management Group Connection can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_management_group_connection.example /providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Network/networkManagerConnections/networkManagerConnection1
```
