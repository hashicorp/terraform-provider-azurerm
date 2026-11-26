---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_consumer_group"
description: |-
    Lists EventHub Consumer Group resources.
---

# List resource: azurerm_eventhub_consumer_group

Lists EventHub Consumer Group resources.

## Example Usage

### List all EventHub Consumer Groups in an EventHub

```hcl
list "azurerm_eventhub_consumer_group" "example" {
  provider = azurerm
  config {
    eventhub_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/eventhub1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `eventhub_id` - (Required) The ID of the EventHub to query.
