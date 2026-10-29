---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_queue"
description: |-
    Lists ServiceBus Queue resources.
---

# List resource: azurerm_servicebus_queue

Lists ServiceBus Queue resources.

## Example Usage

### List all ServiceBus Queues in a ServiceBus Namespace

```hcl
list "azurerm_servicebus_queue" "example" {
  provider = azurerm
  config {
    servicebus_namespace_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.ServiceBus/namespaces/example-namespace"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `servicebus_namespace_id` - (Required) The ID of the ServiceBus Namespace to query.
