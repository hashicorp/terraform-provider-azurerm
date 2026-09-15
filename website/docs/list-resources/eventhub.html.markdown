---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub"
description: |-
    Lists Event Hub resources.
---

# List resource: azurerm_eventhub

Lists Event Hub resources.

## Example Usage

### List all Event Hubs in a Namespace

```hcl
list "azurerm_eventhub" "example" {
  provider = azurerm
  config {
    namespace_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.EventHub/namespaces/example-namespace"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `namespace_id` - (Required) The ID of the Namespace to query.
