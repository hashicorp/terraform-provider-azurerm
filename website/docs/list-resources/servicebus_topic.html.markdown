---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_topic"
description: |-
    Lists ServiceBus Topic resources.
---

# List resource: azurerm_servicebus_topic

Lists ServiceBus Topic resources.

## Example Usage

### List all ServiceBus Topics in the subscription

```hcl
list "azurerm_servicebus_topic" "example" {
  provider = azurerm
  config {}
}
```

### List all ServiceBus Topics in a specific resource group

```hcl
list "azurerm_servicebus_topic" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the resource group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.
