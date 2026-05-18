---
subcategory: "Durable Task"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_durable_task_scheduler"
description: |-
  Manages a Durable Task Scheduler.
---

# azurerm_durable_task_scheduler

Manages a Durable Task Scheduler.

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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Durable Task Scheduler. Changing this forces a new resource to be created.

-> **Note:** `name` must be between `3` and `63` characters, start and end with an alphanumeric character, and only contain alphanumeric characters and hyphens.

* `resource_group_name` - (Required) The name of the Resource Group where the Durable Task Scheduler should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Durable Task Scheduler should exist. Changing this forces a new resource to be created.

* `ip_allow_list` - (Required) A list of IP addresses or CIDR ranges that are allowed to access the Durable Task Scheduler.

* `sku_name` - (Required) The SKU of the Durable Task Scheduler. Possible values are `Consumption` and `Dedicated`. Changing this forces a new resource to be created.

* `capacity` - (Optional) The capacity of the Durable Task Scheduler. Possible values range between `1` and `3`.

~> **Note:** The `capacity` argument must be configured when `sku_name` is set to `Dedicated` and must not be specified when `sku_name` is set to `Consumption`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Durable Task Scheduler.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Durable Task Scheduler.

* `endpoint` - The endpoint URL of the Durable Task Scheduler.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Durable Task Scheduler.
* `read` - (Defaults to 5 minutes) Used when retrieving the Durable Task Scheduler.
* `update` - (Defaults to 30 minutes) Used when updating the Durable Task Scheduler.
* `delete` - (Defaults to 30 minutes) Used when deleting the Durable Task Scheduler.

## Import

A Durable Task Scheduler can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_durable_task_scheduler.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DurableTask/schedulers/scheduler1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DurableTask` - 2025-11-01
