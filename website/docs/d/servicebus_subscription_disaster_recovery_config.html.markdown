---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_subscription_disaster_recovery_config"
description: |-
  Gets information about an existing ServiceBus Subscription disaster recovery config.
---

# Data Source: azurerm_servicebus_subscription_disaster_recovery_config

Use this data source to access information about an existing ServiceBus Subscription disaster recovery config.

## Example Usage

```hcl
data "azurerm_servicebus_subscription_disaster_recovery_config" "example" {
  name         = "example"
  namespace_id = "example"
}

output "config" {
  value = data.azurerm_servicebus_subscription_disaster_recovery_config.example
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the ServiceBus Subscription disaster recovery config.

* `namespace_id` - (Required) The ID of the ServiceBus Namespace where the Service Bus Subscription disaster recovery config exists.

* `alias_authorization_rule_id` - (Optional) The Shared access policies used to access the connection string for the alias.

## Attributes Reference

* `partner_namespace_id` - The ID of the Service Bus Namespace being replicate to.

* `primary_connection_string_alias` - The alias Primary Connection String for the ServiceBus Namespace.

* `secondary_connection_string_alias` - The alias Secondary Connection String for the ServiceBus Namespace

* `default_primary_key` - The primary access key for the authorization rule `RootManageSharedAccessKey`.

* `default_secondary_key` - The secondary access key for the authorization rule `RootManageSharedAccessKey`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the ServiceBus Subscription disaster recovery config.
