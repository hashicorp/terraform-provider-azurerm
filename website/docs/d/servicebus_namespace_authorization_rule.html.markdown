---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_namespace_authorization_rule"
sidebar_current: "docs-azurerm-datasource-servicebus-namespace-authorization-rule"
description: |-
  Gets information about an existing ServiceBus Namespace Authorization Rule.
---

# Data Source: azurerm_servicebus_namespace_authorization_rule

Use this data source to access information about an existing ServiceBus Namespace Authorization Rule.

## Example Usage

```hcl
data "azurerm_servicebus_namespace_authorization_rule" "example" {
  name                = "examplerule"
  namespace_name      = "examplenamespace"
  resource_group_name = "example-resources"
}

output "rule_id" {
  value = "${data.azurerm_servicebus_namespace_authorization_rule.example.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the ServiceBus Namespace Authorization Rule.

* `namespace_name` - (Required) Specifies the name of the ServiceBus Namespace.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the ServiceBus Namespace exists.

## Attributes Reference

* `id` - The id of the ServiceBus Namespace Authorization Rule.

* `primary_connection_string` - The primary connection string for the authorization rule.
    
* `primary_key` - The primary access key for the authorization rule.

* `secondary_connection_string` - The secondary connection string for the authorization rule.

* `secondary_key` - The secondary access key for the authorization rule.
