---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_topic_authorization_rule"
sidebar_current: "docs-azurerm-resource-servicebus-topic-authorization-rule"
description: |-
  Manages a ServiceBus Topic authorization Rule within a ServiceBus Topic.
---

# azurerm_servicebus_topic_authorization_rule

Manages a ServiceBus Topic authorization Rule within a ServiceBus Topic.

## Example Usage

```hcl
variable "location" {
  description = "Azure datacenter to deploy to."
  default = "West US"
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
}

resource "azurerm_servicebus_topic_authorization_rule" "test" {
  name                = "examplerule"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name          = "${azurerm_servicebus_topic.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  listen              = true
  send                = false
  manage              = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the erviceBus Topic Authorization Rule resource. Changing this forces a new resource to be created.

* `namespace_name` - (Required) Specifies the name of the ServiceBus Namespace. Changing this forces a new resource to be created.

* `topic_name` - (Required) Specifies the name of the ServiceBus Topic. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the ServiceBus Namespace exists. Changing this forces a new resource to be created.

~> **NOTE** At least one of the 3 permissions below needs to be set.

* `listen` - (Optional) Does this Authorization Rule have permissions to Listen to the ServiceBus Topic? Defaults to `false`.

* `send` - (Optional) Does this Authorization Rule have permissions to Send to the ServiceBus Topic? Defaults to `false`.

* `manage` - (Optional) Does this Authorization Rule have permissions to Manage to the ServiceBus Topic? When this property is `true` - both `listen` and `send` must be too. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ServiceBus Topic ID.

* `primary_key` - The Primary Key for the ServiceBus Topic authorization Rule.

* `primary_connection_string` - The Primary Connection String for the ServiceBus Topic authorization Rule.

* `secondary_key` - The Secondary Key for the ServiceBus Topic authorization Rule.

* `secondary_connection_string` - The Secondary Connection String for the ServiceBus Topic authorization Rule.

## Import

ServiceBus Topic authorization rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_topic_authorization_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ServiceBus/namespaces/namespace1/topics/topic1/authorizationRules/rule1
```
