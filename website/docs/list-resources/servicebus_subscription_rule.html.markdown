---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_subscription_rule"
description: |-
    Lists ServiceBus Subscription Rule resources.
---

# List resource: azurerm_servicebus_subscription_rule

Lists ServiceBus Subscription Rule resources.

## Example Usage

### List all ServiceBus Subscription Rules in a ServiceBus Subscription

```hcl
list "azurerm_servicebus_subscription_rule" "example" {
  provider = azurerm
  config {
    servicebus_subscription_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.ServiceBus/namespaces/example-namespace/topics/example-topic/subscriptions/example-subscription"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `servicebus_subscription_id` - (Required) The ID of the ServiceBus Subscription to query.
