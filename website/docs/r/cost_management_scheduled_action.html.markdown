---
subcategory: "Cost Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cost_management_scheduled_action"
description: |-
  Manages an Azure Cost Management Scheduled Action.
---

# azurerm_cost_management_scheduled_action

Manages an Azure Cost Management Scheduled Action.

## Example Usage

```hcl
resource "azurerm_cost_management_scheduled_action" "example" {
  name         = "examplescheduledaction"
  display_name = "Report Last 6 Months"

  view_id = "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.CostManagement/views/ms:CostByService"

  email_address_sender = "platformteam@test.com"
  email_subject        = "Cost Management Report"
  email_addresses      = ["example@example.com"]
  message              = "Hi all, take a look at last 6 months spending!"

  frequency  = "Daily"
  start_date = "2023-01-02T00:00:00Z"
  end_date   = "2023-02-02T00:00:00Z"
}
```

## Arguments Reference

The following arguments are supported:

* `display_name` - (Required) User visible input name of the Cost Management Scheduled Action.

* `email_address_sender` - (Required) Email address of the point of contact that should get the unsubscribe requests of Scheduled Action notification emails.

* `email_addresses` - (Required) Specifies a list of email addresses that will receive the Scheduled Action.

* `email_subject` - (Required) Subject of the email. Length is limited to 70 characters.

* `end_date` - (Required) The end date and time of the Scheduled Action (UTC).

* `frequency` - (Required) Frequency of the schedule. Possible values are `Daily`, `Monthly` and `Weekly`. Value `Monthly` requires either `weeks_of_month` and `days_of_week` or `day_of_month` to be specified. Value `Weekly` requires `days_of_week` to be specified.

* `name` - (Required) The name which should be used for this Azure Cost Management Scheduled Action. Changing this forces a new Azure Cost Management Scheduled Action to be created.

* `start_date` - (Required) The start date and time of the Scheduled Action (UTC).

* `view_id` - (Required) The ID of the Cost Management View that is used by the Scheduled Action. Changing this forces a new resource to be created.

---

* `day_of_month` - (Optional) UTC day on which cost analysis data will be emailed. Must be between `1` and `31`. This property is applicable when `frequency` is `Monthly`.

* `days_of_week` - (Optional) Specifies a list of day names on which cost analysis data will be emailed. This property is applicable when frequency is `Weekly` or `Monthly`. Possible values are `Friday`, `Monday`, `Saturday`, `Sunday`, `Thursday`, `Tuesday` and `Wednesday`.

* `hour_of_day` - (Optional) UTC time at which cost analysis data will be emailed. Must be between `0` and `23`.

* `message` - (Optional) Message to be added in the email. Length is limited to 250 characters.

* `weeks_of_month` - (Optional) Specifies a list of weeks in which cost analysis data will be emailed. This property is applicable when `frequency` is `Monthly` and used in combination with `days_of_week`. Possible values are `First`, `Fourth`, `Last`, `Second` and `Third`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Azure Cost Management Scheduled Action.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Cost Management Scheduled Action.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Cost Management Scheduled Action.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Cost Management Scheduled Action.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Cost Management Scheduled Action.

## Import

Azure Cost Management Scheduled Actions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cost_management_scheduled_action.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.CostManagement/scheduledActions/scheduledaction1
```
