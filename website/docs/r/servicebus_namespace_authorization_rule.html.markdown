---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_namespace_authorization_rule"
sidebar_current: "docs-azurerm-resource-servicebus-namespace-authorization-rule"
description: |-
  Manages a ServiceBus Namespace authorization Rule within a ServiceBus Topic.
---

# azurerm_servicebus_namespace_authorization_rule

Manages a ServiceBus Namespace authorization Rule within a ServiceBus Topic.

## Example Usage

```hcl

resource "azurerm_resource_group" "example" {
  name     = "terraform-servicebus"
  location = "West US"
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "${var.servicebus_name}"
  location            = "${var.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "standard"

  tags {
    source = "terraform"
  }
}

resource "azurerm_servicebus_namespace_authorization_rule" "example" {
  name                = "examplerule"
  namespace_name      = "${azurerm_servicebus_namespace.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  rights              = ["Listen", "Send", "Manage"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the ServiceBus Namespace Authorization Rule resource. Changing this forces a new resource to be created.

* `namespace_name` - (Required) Specifies the name of the ServiceBus Namespace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the ServiceBus Namespace exists. Changing this forces a new resource to be created.

* `rights` - (Required) A list of rights grants by this authorization rule. Possible values are `Listen`, `Send` and `Manage`. The `Manage` right requires both `Listen` and `Send` rights.

## Attributes Reference

The following attributes are exported:

* `id` - The ServiceBus Topic ID.

* `primary_key` - The Primary Key for the ServiceBus Namespace authorization Rule.

* `primary_connection_string` - The Primary Connection String for the ServiceBus Namespace authorization Rule.

* `secondary_key` - The Secondary Key for the ServiceBus Namespace authorization Rule.

* `secondary_connection_string` - The Secondary Connection String for the ServiceBus Namespace authorization Rule.

## Import

ServiceBus Namespace authorization rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_namespace_authorization_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ServiceBus/namespaces/namespace1/AuthorizationRules/rule1
```
