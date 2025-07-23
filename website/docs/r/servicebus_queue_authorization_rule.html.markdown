---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_queue_authorization_rule"
description: |-
  Manages an Authorization Rule for a ServiceBus Queue.
---

# azurerm_servicebus_queue_authorization_rule

Manages an Authorization Rule for a ServiceBus Queue.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "terraform-servicebus"
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

resource "azurerm_servicebus_queue" "example" {
  name         = "tfex_servicebus_queue"
  namespace_id = azurerm_servicebus_namespace.example.id

  enable_partitioning = true
}

resource "azurerm_servicebus_queue_authorization_rule" "example" {
  name     = "examplerule"
  queue_id = azurerm_servicebus_queue.example.id

  listen = true
  send   = true
  manage = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Authorization Rule. Changing this forces a new resource to be created.

* `queue_id` - (Required) Specifies the ID of the ServiceBus Queue. Changing this forces a new resource to be created.

~> **Note:** At least one of the 3 permissions below needs to be set.

* `listen` - (Optional) Does this Authorization Rule have Listen permissions to the ServiceBus Queue? Defaults to `false`.

* `send` - (Optional) Does this Authorization Rule have Send permissions to the ServiceBus Queue? Defaults to `false`.

* `manage` - (Optional) Does this Authorization Rule have Manage permissions to the ServiceBus Queue? When this property is `true` - both `listen` and `send` must be too. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Authorization Rule.

* `primary_key` - The Primary Key for the Authorization Rule.

* `primary_connection_string` - The Primary Connection String for the Authorization Rule.

* `secondary_key` - The Secondary Key for the Authorization Rule.

* `secondary_connection_string` - The Secondary Connection String for the Authorization Rule.

* `primary_connection_string_alias` - The alias Primary Connection String for the ServiceBus Namespace, if the namespace is Geo DR paired.

* `secondary_connection_string_alias` - The alias Secondary Connection String for the ServiceBus Namespace

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the ServiceBus Queue Authorization Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the ServiceBus Queue Authorization Rule.
* `update` - (Defaults to 30 minutes) Used when updating the ServiceBus Queue Authorization Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the ServiceBus Queue Authorization Rule.

## Import

ServiceBus Queue Authorization Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_queue_authorization_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ServiceBus/namespaces/namespace1/queues/queue1/authorizationRules/rule1
```
