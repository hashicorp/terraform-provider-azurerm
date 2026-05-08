---
subcategory: "Durable Task"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_durable_task_hub"
description: |-
  Lists Durable Task Hub resources.
---

# List resource: azurerm_durable_task_hub

Lists Durable Task Hub resources.

## Example Usage

### List Durable Task Hubs in a Scheduler

```hcl
list "azurerm_durable_task_hub" "example" {
  provider = azurerm
  config {
    scheduler_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DurableTask/schedulers/scheduler1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `scheduler_id` - (Required) The ID of the Durable Task Scheduler to query.
