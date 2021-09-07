---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_trigger_schedule"
description: |-
  Manages a Trigger Schedule inside a Azure Data Factory.
---

# azurerm_data_factory_trigger_schedule

Manages a Trigger Schedule inside a Azure Data Factory.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name                = "example"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
}

resource "azurerm_data_factory_trigger_schedule" "test" {
  name                = "example"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  pipeline_name       = azurerm_data_factory_pipeline.test.name

  interval  = 5
  frequency = "Day"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Schedule Trigger. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Factory Schedule Trigger. Changing this forces a new resource

* `data_factory_name` - (Required) The Data Factory name in which to associate the Schedule Trigger with. Changing this forces a new resource.

* `pipeline_name` - (Required) The Data Factory Pipeline name that the trigger will act on.

* `description` - (Optional) The Schedule Trigger's description.

* `schedule` - (Optional) The recurrence schedule for the trigger. A trigger with a specified frequency value alters its recurrence based on a recurrence schedule. The schedule property contains modifications for the recurrence that are based on minutes, hours, weekdays, month days, and week number. A `schedule` block as documented below.

* `start_time` - (Optional) The time the Schedule Trigger will start. This defaults to the current time. The time will be represented in UTC.

* `end_time` - (Optional) The time the Schedule Trigger should end. The time will be represented in UTC.

* `interval` - (Optional) The interval for how often the trigger occurs. This defaults to 1.

* `frequency` - (Optional) The trigger frequency. Valid values include `Minute`, `Hour`, `Day`, `Week`, `Month`. Defaults to `Minute`.

* `pipeline_parameters` - (Optional) The pipeline parameters that the trigger will act upon.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Schedule Trigger.

---

`schedule` supports the following:

* `days_of_month` - (Optional) 	Day of the month on which the trigger runs. The value can be specified with a monthly frequency only.

* `days_of_week` - (Optional) Days of the week the trigger runs. The value can be specified only with a weekly frequency.

* `hours` - (Optional) Hours of the day at which the trigger runs.

* `minutes` - (Optional) Minutes of the hour at which the trigger runs.

* `monthly_occurrence` - (Optional) Days of the month on which the trigger runs. The value can be specified with a monthly frequency only. A `monthly_occurrence` block as documented below.

---

`monthly_occurrence` supports the following:

* `day` - (Required) The day of the week on which the trigger runs. For example, a `monthly_occurrence` property with a day value of `Sunday` means every Sunday of the month.

* `occurrence` - (Optional) The occurrence of the specified day during the month. For example, a `monthly_occurrence` property with day and occurrence values of `Sunday, -1` means the last Sunday of the month.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Factory Schedule Trigger.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Schedule Trigger.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Schedule Trigger.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Schedule Trigger.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Schedule Trigger.

## Import

Data Factory Schedule Trigger can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_schedule_trigger.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/triggers/example
```
