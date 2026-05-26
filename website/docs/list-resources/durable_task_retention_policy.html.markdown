---
subcategory: "Durable Task"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_durable_task_retention_policy"
description: |-
  Lists Durable Task Retention Policy resources.
---

# List resource: azurerm_durable_task_retention_policy

Lists Durable Task Retention Policy resources.

## Example Usage

### List Durable Task Retention Policies in a Scheduler

```hcl
list "azurerm_durable_task_retention_policy" "example" {
  provider = azurerm
  config {
    durable_task_scheduler_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DurableTask/schedulers/scheduler1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `durable_task_scheduler_id` - (Required) The ID of the Durable Task Scheduler to query.
