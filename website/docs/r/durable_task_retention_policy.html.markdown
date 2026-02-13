---
subcategory: "Durable Task"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_durable_task_retention_policy"
description: |-
  Manages a Durable Task Retention Policy.
---

# azurerm_durable_task_retention_policy

Manages a Durable Task Retention Policy for a Durable Task Scheduler.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_durable_task_scheduler" "example" {
  name                = "example-scheduler"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Consumption"
  ip_allow_list       = ["0.0.0.0/0"]
}

resource "azurerm_durable_task_retention_policy" "example" {
  scheduler_id = azurerm_durable_task_scheduler.example.id

  retention_policy {
    retention_period_in_days = 30
    orchestration_state      = "Completed"
  }

  retention_policy {
    retention_period_in_days = 7
    orchestration_state      = "Failed"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `scheduler_id` - (Required) The ID of the Durable Task Scheduler where the Retention Policy should be applied. Changing this forces a new resource to be created.

* `retention_policy` - (Required) One or more `retention_policy` blocks as defined below.

---

A `retention_policy` block supports the following:

* `retention_period_in_days` - (Required) The number of days to retain orchestration data.

* `orchestration_state` - (Optional) The orchestration state to which this retention policy applies. Possible values are `Completed`, `Failed`, `Terminated`, and `Canceled`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Durable Task Retention Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Durable Task Retention Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Durable Task Retention Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Durable Task Retention Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Durable Task Retention Policy.

## Import

Durable Task Retention Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_durable_task_retention_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DurableTask/schedulers/scheduler1/retentionPolicies/default
```
