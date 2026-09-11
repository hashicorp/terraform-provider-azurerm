---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_subscription"
description: |-
    Lists ServiceBus Subscription resources.
---

# List resource: azurerm_servicebus_subscription

Lists ServiceBus Subscription resources.

## Example Usage

### List all ServiceBus Subscriptions in a ServiceBus Topic

```hcl
list "azurerm_servicebus_subscription" "example" {
  provider = azurerm
  config {
    topic_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.ServiceBus/namespaces/example-namespace/topics/example-topic"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `topic_id` - (Required) The ID of the ServiceBus Topic to query.
