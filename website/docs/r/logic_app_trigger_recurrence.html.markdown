---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_trigger_recurrence"
sidebar_current: "docs-azurerm-resource-logic-app-trigger-recurrence"
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
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_logic_app_trigger_recurrence" "example" {
  name         = "run-every-day"
  logic_app_id = "${azurerm_logic_app_workflow.example.id}"
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

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Recurrence Trigger within the Logic App Workflow.

## Import

Logic App Recurrence Triggers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_trigger_recurrence.daily /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Logic/workflows/workflow1/triggers/daily
```

-> **NOTE:** This ID is unique to Terraform and doesn't directly match to any other resource. To compose this ID, you can take the ID Logic App Workflow and append `/triggers/{name of the trigger}`.
