---
subcategory: "Durable Task"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_durable_task_hub"
description: |-
  Manages a Durable Task Hub.
---

# azurerm_durable_task_hub

Manages a Durable Task Hub.

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

resource "azurerm_durable_task_hub" "example" {
  name         = "example-taskhub"
  scheduler_id = azurerm_durable_task_scheduler.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Durable Task Hub. Changing this forces a new resource to be created.

* `scheduler_id` - (Required) The ID of the Durable Task Scheduler where the Task Hub should be created. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Durable Task Hub.

* `dashboard_url` - The URL of the dashboard for the Durable Task Hub.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Durable Task Hub.
* `read` - (Defaults to 5 minutes) Used when retrieving the Durable Task Hub.
* `delete` - (Defaults to 30 minutes) Used when deleting the Durable Task Hub.

## Import

Durable Task Hubs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_durable_task_hub.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DurableTask/schedulers/scheduler1/taskHubs/taskHub1
```
