---
subcategory: "Durable Task"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_durable_task_scheduler"
description: |-
  Lists Durable Task Scheduler resources.
---

# List resource: azurerm_durable_task_scheduler

Lists Durable Task Scheduler resources.

## Example Usage

### List all Durable Task Schedulers in the subscription

```hcl
list "azurerm_durable_task_scheduler" "example" {
  provider = azurerm
  config {}
}
```

### List all Durable Task Schedulers in a specific resource group

```hcl
list "azurerm_durable_task_scheduler" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the Resource Group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.
