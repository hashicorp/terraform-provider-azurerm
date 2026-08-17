---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_scheduled_query_rules_alert_v2"
description: |-
  Lists Scheduled Query Rules Alert V2 resources.
---

# List resource: azurerm_monitor_scheduled_query_rules_alert_v2

Lists Scheduled Query Rules Alert V2 resources.

## Example Usage

### List all Scheduled Query Rules Alert V2 resources in the subscription

```hcl
list "azurerm_monitor_scheduled_query_rules_alert_v2" "example" {
  provider = azurerm
  config {}
}
```

### List all Scheduled Query Rules Alert V2 resources in a specific resource group

```hcl
list "azurerm_monitor_scheduled_query_rules_alert_v2" "example" {
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
