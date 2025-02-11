---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_namespace_authorization_rule"
description: |-
  Gets information about an existing ServiceBus Namespace Authorization Rule.
---

# Data Source: azurerm_servicebus_namespace_authorization_rule

Use this data source to access information about an existing ServiceBus Namespace Authorization Rule.

## Example Usage

```hcl
data "azurerm_servicebus_namespace_authorization_rule" "example" {
  name         = "examplerule"
  namespace_id = "examplenamespace"
}

output "rule_id" {
  value = data.azurerm_servicebus_namespace_authorization_rule.example.id
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the ServiceBus Namespace Authorization Rule.

* `namespace_id` - (Required) Specifies the ID of the ServiceBus Namespace where the Service Bus Namespace Authorization Rule exists.

## Attributes Reference

* `id` - The id of the ServiceBus Namespace Authorization Rule.

* `primary_connection_string` - The primary connection string for the authorization rule.

* `primary_key` - The primary access key for the authorization rule.

* `secondary_connection_string` - The secondary connection string for the authorization rule.

* `secondary_key` - The secondary access key for the authorization rule.

* `primary_connection_string_alias` - The alias Primary Connection String for the ServiceBus Namespace, if the namespace is Geo DR paired.

* `secondary_connection_string_alias` - The alias Secondary Connection String for the ServiceBus Namespace

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the ServiceBus Namespace Authorization Rule.
