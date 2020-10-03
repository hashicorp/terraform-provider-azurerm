---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_subscription_rule"
description: |-
  Manages a ServiceBus Subscription Rule.
---

# azurerm_servicebus_subscription_rule

Manages a ServiceBus Subscription Rule.

## Example Usage (SQL Filter)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tfex-servicebus-subscription-rule-sql"
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
  name                = "tfex_servicebus_topic"
  resource_group_name = azurerm_resource_group.example.name
  namespace_name      = azurerm_servicebus_namespace.example.name

  enable_partitioning = true
}

resource "azurerm_servicebus_subscription" "example" {
  name                = "tfex_servicebus_subscription"
  resource_group_name = azurerm_resource_group.example.name
  namespace_name      = azurerm_servicebus_namespace.example.name
  topic_name          = azurerm_servicebus_topic.example.name
  max_delivery_count  = 1
}

resource "azurerm_servicebus_subscription_rule" "example" {
  name                = "tfex_servicebus_rule"
  resource_group_name = azurerm_resource_group.example.name
  namespace_name      = azurerm_servicebus_namespace.example.name
  topic_name          = azurerm_servicebus_topic.example.name
  subscription_name   = azurerm_servicebus_subscription.example.name
  filter_type         = "SqlFilter"
  sql_filter          = "colour = 'red'"
}
```

## Example Usage (Correlation Filter)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tfex-servicebus-subscription-rule-cor"
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
  name                = "tfex_servicebus_topic"
  resource_group_name = azurerm_resource_group.example.name
  namespace_name      = azurerm_servicebus_namespace.example.name

  enable_partitioning = true
}

resource "azurerm_servicebus_subscription" "example" {
  name                = "tfex_servicebus_subscription"
  resource_group_name = azurerm_resource_group.example.name
  namespace_name      = azurerm_servicebus_namespace.example.name
  topic_name          = azurerm_servicebus_topic.example.name
  max_delivery_count  = 1
}

resource "azurerm_servicebus_subscription_rule" "example" {
  name                = "tfex_servicebus_rule"
  resource_group_name = azurerm_resource_group.example.name
  namespace_name      = azurerm_servicebus_namespace.example.name
  topic_name          = azurerm_servicebus_topic.example.name
  subscription_name   = azurerm_servicebus_subscription.example.name
  filter_type         = "CorrelationFilter"

  correlation_filter {
    correlation_id = "high"
    label          = "red"
    properties = {
      customProperty = "value"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the ServiceBus Subscription Rule. Changing this forces a new resource to be created.

* `namespace_name` - (Required) The name of the ServiceBus Namespace in which the ServiceBus Topic exists. Changing this forces a new resource to be created.

* `topic_name` - (Required) The name of the ServiceBus Topic in which the ServiceBus Subscription exists. Changing this forces a new resource to be created.

* `subscription_name` - (Required) The name of the ServiceBus Subscription in which this Rule should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in the ServiceBus Namespace exists. Changing this forces a new resource to be created.

* `filter_type` - (Required) Type of filter to be applied to a BrokeredMessage. Possible values are `SqlFilter` and `CorrelationFilter`.

* `sql_filter` - (Optional) Represents a filter written in SQL language-based syntax that to be evaluated against a BrokeredMessage. Required when `filter_type` is set to `SqlFilter`.

* `correlation_filter` - (Optional) A `correlation_filter` block as documented below to be evaluated against a BrokeredMessage. Required when `filter_type` is set to `CorrelationFilter`.

* `action` - (Optional) Represents set of actions written in SQL language-based syntax that is performed against a BrokeredMessage.

`correlation_filter` supports the following:

* `content_type` - (Optional) Content type of the message.

* `correlation_id` - (Optional) Identifier of the correlation.

* `label` - (Optional) Application specific label.

* `message_id` - (Optional) Identifier of the message.

* `reply_to` - (Optional) Address of the queue to reply to.

* `reply_to_session_id` - (Optional) Session identifier to reply to.

* `session_id` - (Optional) Session identifier.

* `to` - (Optional) Address to send to.

* `properties` - (Optional) A list of user defined properties to be included in the filter. Specified as a map of name/value pairs.

~> **NOTE:** When creating a subscription rule of type `CorrelationFilter` at least one property must be set in the `correlation_filter` block.


## Attributes Reference

The following attributes are exported:

* `id` - The ServiceBus Subscription Rule ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the ServiceBus Subscription Rule.
* `update` - (Defaults to 30 minutes) Used when updating the ServiceBus Subscription Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the ServiceBus Subscription Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the ServiceBus Subscription Rule.

## Import

Service Bus Subscription Rule can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_subscription_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.servicebus/namespaces/sbns1/topics/sntopic1/subscriptions/sbsub1/rules/sbrule1
```
