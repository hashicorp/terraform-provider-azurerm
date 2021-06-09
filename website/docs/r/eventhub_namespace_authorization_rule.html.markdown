---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_namespace_authorization_rule"
description: |-
  Manages an Authorization Rule for an Event Hub Namespace.
---

# azurerm_eventhub_namespace_authorization_rule

Manages an Authorization Rule for an Event Hub Namespace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "resourcegroup"
  location = "West Europe"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "acceptanceTestEventHubNamespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Basic"
  capacity            = 2

  tags = {
    environment = "Production"
  }
}

resource "azurerm_eventhub_namespace_authorization_rule" "example" {
  name                = "navi"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name

  listen = true
  send   = false
  manage = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Authorization Rule. Changing this forces a new resource to be created.

* `namespace_name` - (Required) Specifies the name of the EventHub Namespace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the EventHub Namespace exists. Changing this forces a new resource to be created.

~> **NOTE** At least one of the 3 permissions below needs to be set.

* `listen` - (Optional) Grants listen access to this this Authorization Rule. Defaults to `false`.

* `send` - (Optional) Grants send access to this this Authorization Rule. Defaults to `false`.

* `manage` - (Optional) Grants manage access to this this Authorization Rule. When this property is `true` - both `listen` and `send` must be too. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The EventHub Namespace Authorization Rule ID.

* `primary_connection_string_alias` - The alias of the Primary Connection String for the Authorization Rule, which is generated when disaster recovery is enabled.

* `secondary_connection_string_alias` - The alias of the Secondary Connection String for the Authorization Rule, which is generated when disaster recovery is enabled.

* `primary_connection_string` - The Primary Connection String for the Authorization Rule.

* `primary_key` - The Primary Key for the Authorization Rule.

* `secondary_connection_string` - The Secondary Connection String for the Authorization Rule.

* `secondary_key` - The Secondary Key for the Authorization Rule.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventHub Namespace Authorization Rule.
* `update` - (Defaults to 30 minutes) Used when updating the EventHub Namespace Authorization Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventHub Namespace Authorization Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the EventHub Namespace Authorization Rule.

## Import

EventHub Namespace Authorization Rules can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_eventhub_namespace_authorization_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/authorizationRules/rule1
```
