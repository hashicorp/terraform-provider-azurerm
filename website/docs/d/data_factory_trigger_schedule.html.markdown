---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_data_factory_trigger_schedule"
description: |-
  Gets information about a trigger schedule in Azure Data Factory.
---

# Data Source: azurerm_data_factory_trigger_schedule

Use this data source to access information about a trigger schedule in Azure Data Factory.

## Example Usage

```hcl
data "azurerm_data_factory_trigger_schedule" "example" {
  name            = "example_trigger"
  data_factory_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.DataFactory/factories/datafactory1"
}

output "id" {
  value = data.azurerm_data_factory_trigger_schedule.example.id
}
```

## Arguments Reference

The following arguments are supported:

- `name` - (Required) The name of the trigger schedule.

- `data_factory_id` - (Required) The ID of the Azure Data Factory to fetch trigger schedule from.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `id` - The ID of the Azure Data Factory trigger schedule.

* `description` - The Schedule Trigger's description.

* `schedule` - A `schedule` block as described below, which further specifies the recurrence schedule for the trigger.

* `start_time` - The time the Schedule Trigger will start. The time will be represented in UTC.

* `time_zone` - The timezone of the start/end time.

* `end_time` - The time the Schedule Trigger should end. The time will be represented in UTC.

* `interval` - The interval for how often the trigger occurs.

* `frequency` - The trigger frequency.

* `activated` - Specifies if the Data Factory Schedule Trigger is activated.

* `pipeline_name` - The Data Factory Pipeline name that the trigger will act on.

* `annotations` - List of tags that can be used for describing the Data Factory Schedule Trigger.

---

A `schedule` block exports the following:

* `days_of_month` - Day(s) of the month on which the trigger is scheduled.

* `days_of_week` - Day(s) of the week on which the trigger is scheduled.

* `hours` - Hours of the day on which the trigger is scheduled.

* `minutes` - Minutes of the hour on which the trigger is scheduled.

* `monthly` - A `monthly` block as documented below, which specifies the days of the month on which the trigger is scheduled.

---

A `monthly` block exports the following:

* `weekday` - The day of the week on which the trigger runs.

* `week` - The occurrence of the specified day during the month.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Data Factory trigger schedule.
