---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_rule"
sidebar_current: "docs-azurerm-resource-servicebus-rule"
description: |-
  Create a ServiceBus Rule.
---

# azurerm\_servicebus\_rule

Create a ServiceBus Rule.

## Example Usage (SQL Filter)

```hcl
variable "location" {
  description = "Azure datacenter to deploy to."
  default = "West US"
}

variable "servicebus_name" {
  description = "Input your unique Azure service bus name"
}

resource "azurerm_resource_group" "test" {
  name     = "terraform-servicebus"
  location = "${var.location}"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "${var.servicebus_name}"
  location            = "${var.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "standard"

  tags {
    source = "terraform"
  }
}

resource "azurerm_servicebus_topic" "test" {
  name                = "testTopic"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"

  enable_partitioning = true
}

resource "azurerm_servicebus_subscription" "test" {
  name                = "testSubscription"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name          = "${azurerm_servicebus_topic.test.name}"
  max_delivery_count  = 1
}

resource "azurerm_servicebus_rule" "test" {
  name                = "testRule"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name          = "${azurerm_servicebus_topic.test.name}"
  subscription_name   = "${azurerm_servicebus_subscription.test.name}"
  filter_type         = "SqlFilter"
  sql_filter          = "color = 'red'"
}
```

## Example Usage (Correlation Filter)

```hcl
variable "location" {
  description = "Azure datacenter to deploy to."
  default = "West US"
}

variable "servicebus_name" {
  description = "Input your unique Azure service bus name"
}

resource "azurerm_resource_group" "test" {
  name     = "terraform-servicebus"
  location = "${var.location}"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "${var.servicebus_name}"
  location            = "${var.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "standard"

  tags {
    source = "terraform"
  }
}

resource "azurerm_servicebus_topic" "test" {
  name                = "testTopic"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"

  enable_partitioning = true
}

resource "azurerm_servicebus_subscription" "test" {
  name                = "testSubscription"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name          = "${azurerm_servicebus_topic.test.name}"
  max_delivery_count  = 1
}

resource "azurerm_servicebus_rule" "test" {
  name                = "testRule"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name          = "${azurerm_servicebus_topic.test.name}"
  subscription_name   = "${azurerm_servicebus_subscription.test.name}"
  filter_type         = "CorrelationFilter"
  correlation_filter  = {
    correlation_id = "high"
    label          = "red"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the ServiceBus Subscription resource. Changing this forces a new resource to be created.

* `namespace_name` - (Required) The name of the ServiceBus Namespace to create this Subscription in. Changing this forces a new resource to be created.

* `topic_name` - (Required) The name of the ServiceBus Topic to create this Subscription in. Changing this forces a new resource to be created.

* `subscription_name` - (Required) The name of the ServiceBus Subscription to create this Rule in. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the namespace. Changing this forces a new resource to be created.

* `filter_type` - (Required) Type of filter to be applied to a BrokeredMessage. Possible values are `SqlFilter` and `CorrelationFilter`.

* `sql_filter` - (Optional) Represents a filter written in SQL language-based syntax that to be evaluated against a BrokeredMessage. Must be set when `filter_type` is set to `SqlFilter`.

* `correlation_filter` - (Optional) A `correlation_filter` block as documented below to be evaluated against a BrokeredMessage. Must be set when `filter_type` is set to `CorrelationFilter`.

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

~> **NOTE:** When creating a rule of type `CorrelationFilter` at least one property must be set in the `correlation_filter` block.


## Attributes Reference

The following attributes are exported:

* `id` - The ServiceBus Rule ID.

## Import

Service Bus Rule can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_subscription.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.servicebus/namespaces/sbns1/topics/sntopic1/subscriptions/sbsub1/rules/sbrule1
```
