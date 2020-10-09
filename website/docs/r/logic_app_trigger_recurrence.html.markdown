---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_trigger_recurrence"
description: |-
  Manages a Recurrence Trigger within a Logic App Workflow
---

# azurerm_logic_app_trigger_recurrence

Manages a Recurrence Trigger within a Logic App Workflow

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "workflow-resources"
  location = "East US"
}

resource "azurerm_logic_app_workflow" "example" {
  name                = "workflow1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_logic_app_trigger_recurrence" "example" {
  name         = "run-every-day"
  logic_app_id = azurerm_logic_app_workflow.example.id
  frequency    = "Day"
  interval     = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Recurrence Triggers to be created within the Logic App Workflow. Changing this forces a new resource to be created.

-> **NOTE:** This name must be unique across all Triggers within the Logic App Workflow.

* `logic_app_id` - (Required) Specifies the ID of the Logic App Workflow. Changing this forces a new resource to be created.

* `frequency` - (Required) Specifies the Frequency at which this Trigger should be run. Possible values include `Month`, `Week`, `Day`, `Hour`, `Minute` and `Second`.

* `interval` - (Required) Specifies interval used for the Frequency, for example a value of `4` for `interval` and `hour` for `frequency` would run the Trigger every 4 hours.

* `start_time` - (Optional) Specifies the start date and time for this trigger in RFC3339 format: `2000-01-02T03:04:05Z`.

* `time_zone` - (Optional) Specifies the time zone for this trigger.  Supported time zone options are listed [here](https://support.microsoft.com/en-us/help/973627/microsoft-time-zone-index-values)

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Recurrence Trigger within the Logic App Workflow.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Logic App Recurrence Trigger.
* `update` - (Defaults to 30 minutes) Used when updating the Logic App Recurrence Trigger.
* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Recurrence Trigger.
* `delete` - (Defaults to 30 minutes) Used when deleting the Logic App Recurrence Trigger.

## Import

Logic App Recurrence Triggers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_trigger_recurrence.daily /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Logic/workflows/workflow1/triggers/daily
```

-> **NOTE:** This ID is unique to Terraform and doesn't directly match to any other resource. To compose this ID, you can take the ID Logic App Workflow and append `/triggers/{name of the trigger}`.
