---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_topic_authorization_rule"
description: |-
  Manages a ServiceBus Topic authorization Rule within a ServiceBus Topic.
---

# azurerm_servicebus_topic_authorization_rule

Manages a ServiceBus Topic authorization Rule within a ServiceBus Topic.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tfex-servicebus"
  location = "West Europe"
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "tfex-servicebus-namespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"

  tags = {
    source = "terraform"
  }
}

resource "azurerm_servicebus_topic" "example" {
  name         = "tfex_servicebus_topic"
  namespace_id = azurerm_servicebus_namespace.example.id
}

resource "azurerm_servicebus_topic_authorization_rule" "example" {
  name     = "tfex_servicebus_topic_sasPolicy"
  topic_id = azurerm_servicebus_topic.example.id
  listen   = true
  send     = false
  manage   = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the ServiceBus Topic Authorization Rule resource. Changing this forces a new resource to be created.

* `topic_id` - (Required) Specifies the ID of the ServiceBus Topic. Changing this forces a new resource to be created.

~> **Note:** At least one of the 3 permissions below needs to be set.

* `listen` - (Optional) Grants listen access to this this Authorization Rule. Defaults to `false`.

* `send` - (Optional) Grants send access to this this Authorization Rule. Defaults to `false`.

* `manage` - (Optional) Grants manage access to this this Authorization Rule. When this property is `true` - both `listen` and `send` must be too. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ServiceBus Topic ID.

* `primary_key` - The Primary Key for the ServiceBus Topic authorization Rule.

* `primary_connection_string` - The Primary Connection String for the ServiceBus Topic authorization Rule.

* `secondary_key` - The Secondary Key for the ServiceBus Topic authorization Rule.

* `secondary_connection_string` - The Secondary Connection String for the ServiceBus Topic authorization Rule.

* `primary_connection_string_alias` - The alias Primary Connection String for the ServiceBus Namespace, if the namespace is Geo DR paired.

* `secondary_connection_string_alias` - The alias Secondary Connection String for the ServiceBus Namespace

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the ServiceBus Topic Authorization Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the ServiceBus Topic Authorization Rule.
* `update` - (Defaults to 30 minutes) Used when updating the ServiceBus Topic Authorization Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the ServiceBus Topic Authorization Rule.

## Import

ServiceBus Topic authorization rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_topic_authorization_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ServiceBus/namespaces/namespace1/topics/topic1/authorizationRules/rule1
```
