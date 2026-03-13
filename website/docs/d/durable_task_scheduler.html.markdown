---
subcategory: "Durable Task"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_durable_task_scheduler"
description: |-
  Gets information about an existing Durable Task Scheduler.
---

# Data Source: azurerm_durable_task_scheduler

Use this data source to access information about an existing Durable Task Scheduler.

## Example Usage

```hcl
data "azurerm_durable_task_scheduler" "example" {
  name                = "existing-durable-task-scheduler"
  resource_group_name = "existing-resources"
}

output "endpoint" {
  value = data.azurerm_durable_task_scheduler.example.endpoint
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Durable Task Scheduler.

* `resource_group_name` - (Required) The name of the Resource Group where the Durable Task Scheduler exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Durable Task Scheduler.

* `capacity` - The capacity of the Durable Task Scheduler.

* `endpoint` - The endpoint URL of the Durable Task Scheduler.

* `ip_allow_list` - A list of IP addresses or CIDR ranges that are allowed to access the Durable Task Scheduler.

* `location` - The Azure Region where the Durable Task Scheduler exists.

* `redundancy_state` - The redundancy state of the Durable Task Scheduler.

* `sku_name` - The SKU of the Durable Task Scheduler.

* `tags` - A mapping of tags assigned to the Durable Task Scheduler.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Durable Task Scheduler.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DurableTask` - 2025-11-01
