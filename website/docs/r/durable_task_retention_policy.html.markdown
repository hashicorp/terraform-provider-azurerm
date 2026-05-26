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
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_durable_task_scheduler" "example" {
  name                = "example-durable-task-scheduler"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Consumption"
  ip_allow_list       = ["0.0.0.0/0"]
}

resource "azurerm_durable_task_retention_policy" "example" {
  durable_task_scheduler_id = azurerm_durable_task_scheduler.example.id

  retention_policy {
    canceled_retention_period_in_days   = 21
    completed_retention_period_in_days  = 30
    failed_retention_period_in_days     = 7
    terminated_retention_period_in_days = 14
  }
}
```

## Arguments Reference

The following arguments are supported:

* `durable_task_scheduler_id` - (Required) The ID of the Durable Task Scheduler where the Retention Policy should be applied. Changing this forces a new resource to be created.

* `retention_policy` - (Required) A `retention_policy` block as defined below.

---

A `retention_policy` block supports the following:

* `canceled_retention_period_in_days` - (Optional) The number of days to retain canceled orchestration data. Possible values range between `1` and `90`.

* `completed_retention_period_in_days` - (Optional) The number of days to retain completed orchestration data. Possible values range between `1` and `90`.

* `default_retention_period_in_days` - (Optional) The default number of days to retain orchestration data. Possible values range between `1` and `90`.

-> **Note:** `default_retention_period_in_days` cannot be configured together with the state-specific retention period fields.

* `failed_retention_period_in_days` - (Optional) The number of days to retain failed orchestration data. Possible values range between `1` and `90`.

* `terminated_retention_period_in_days` - (Optional) The number of days to retain terminated orchestration data. Possible values range between `1` and `90`.

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

A Durable Task Retention Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_durable_task_retention_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DurableTask/schedulers/scheduler1/retentionPolicies/default
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DurableTask` - 2025-11-01
