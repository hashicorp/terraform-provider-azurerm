---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_namespace_authorization_rule"
description: |-
  Manages a ServiceBus Namespace authorization Rule within a ServiceBus.
---

# azurerm_servicebus_namespace_authorization_rule

Manages a ServiceBus Namespace authorization Rule within a ServiceBus.

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

resource "azurerm_servicebus_namespace_authorization_rule" "example" {
  name         = "examplerule"
  namespace_id = azurerm_servicebus_namespace.example.id

  listen = true
  send   = true
  manage = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the ServiceBus Namespace Authorization Rule resource. Changing this forces a new resource to be created.

* `namespace_id` - (Required) Specifies the ID of the ServiceBus Namespace. Changing this forces a new resource to be created.

~> **Note:** At least one of the 3 permissions below needs to be set.

* `listen` - (Optional) Grants listen access to this Authorization Rule. Defaults to `false`.

* `send` - (Optional) Grants send access to this Authorization Rule. Defaults to `false`.

* `manage` - (Optional) Grants manage access to this Authorization Rule. When this property is `true` - both `listen` and `send` must be too. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ServiceBus Topic ID.

* `primary_key` - The Primary Key for the ServiceBus Namespace authorization Rule.

* `primary_connection_string` - The Primary Connection String for the ServiceBus Namespace authorization Rule.

* `secondary_key` - The Secondary Key for the ServiceBus Namespace authorization Rule.

* `secondary_connection_string` - The Secondary Connection String for the ServiceBus Namespace authorization Rule.

* `primary_connection_string_alias` - The alias Primary Connection String for the ServiceBus Namespace, if the namespace is Geo DR paired.

* `secondary_connection_string_alias` - The alias Secondary Connection String for the ServiceBus Namespace

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the ServiceBus Namespace Authorization Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the ServiceBus Namespace Authorization Rule.
* `update` - (Defaults to 30 minutes) Used when updating the ServiceBus Namespace Authorization Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the ServiceBus Namespace Authorization Rule.

## Import

ServiceBus Namespace authorization rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_namespace_authorization_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ServiceBus/namespaces/namespace1/authorizationRules/rule1
```
