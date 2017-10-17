---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_scheduled_time"
sidebar_current: "docs-azurerm-datasource-scheduled-time"
description: |-
  Get information about the specified scheduled time.
---

# azurerm\_scheduled\_time

Use this data source to access the properties of an existing Scheduled Time.

## Example Usage

```hcl
data "azurerm_scheduled_time" "test" {
        "frequency" = "Day"
        "hour" = "14"
        "minute" = "25"
        "minimum_delay_from_now_in_minutes" = "6"
}
```

## Argument Reference

* `frequency` - (Required) Specifies the frequency of the schedule. - can be either `OneTime`, `Day`, `Hour`, `Week`, or `Month`.
* `second` - (Optional) Specifies the second part of the schedule time.
* `minute` - (Optional) Specifies the minute part of the schedule time.
* `hour` - (Optional) Specifies the hour part of the schedule time.
* `day_of_week` - (Optional) Specifies the day of the week when the schedule time will be triggered. Conflicts with `day_of_week`.
* `day_of_month` - (Optional) Specifies the day of the month when the schedule time will be triggered. Conflicts with `day_of_month`.
* `minimum_delay_from_now_in_minutes` - (Optional) Minimum delay for the calculation of the first run time.
* `timezone` - (Optional) Timezone of the scheduled time.

## Attributes Reference

The following attributes are exported:
* `next_run_time` - The value when should this schedule next run. This is computed from the previous parameters.
