---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_topic_authorization_rule"
description: |-
    Lists ServiceBus Topic Authorization Rule resources.
---

# List resource: azurerm_servicebus_topic_authorization_rule

Lists ServiceBus Topic Authorization Rule resources.

## Example Usage

### List all ServiceBus Topic Authorization Rules in a Topic

```hcl
list "azurerm_servicebus_topic_authorization_rule" "example" {
  provider = azurerm
  config {
    servicebus_topic_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.ServiceBus/namespaces/example-namespace/topics/example-topic"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `servicebus_topic_id` - (Required) The ID of the Topic to query.
