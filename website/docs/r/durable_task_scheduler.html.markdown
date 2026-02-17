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
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_durable_task_scheduler" "example" {
  name                = "example-scheduler"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Consumption"
  ip_allow_list       = ["0.0.0.0/0"]

  tags = {
    environment = "production"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Durable Task Scheduler. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Durable Task Scheduler should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Durable Task Scheduler should exist. Changing this forces a new resource to be created.

* `sku_name` - (Required) The SKU of the Durable Task Scheduler. Possible values are `Consumption` and `Dedicated`.

* `ip_allow_list` - (Required) A list of IP addresses or CIDR ranges that are allowed to access the Durable Task Scheduler.

---

* `capacity` - (Optional) The capacity of the Durable Task Scheduler. This is only applicable when `sku_name` is set to `Dedicated`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Durable Task Scheduler.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Durable Task Scheduler.

* `endpoint` - The endpoint URL of the Durable Task Scheduler.

* `redundancy_state` - The redundancy state of the Durable Task Scheduler. Possible values are `None` and `Zone`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Durable Task Scheduler.
* `read` - (Defaults to 5 minutes) Used when retrieving the Durable Task Scheduler.
* `update` - (Defaults to 30 minutes) Used when updating the Durable Task Scheduler.
* `delete` - (Defaults to 30 minutes) Used when deleting the Durable Task Scheduler.

## Import

Durable Task Schedulers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_durable_task_scheduler.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DurableTask/schedulers/scheduler1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DurableTask` - 2025-11-01
