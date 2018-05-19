---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_workflow"
sidebar_current: "docs-azurerm-resource-logic-app-workflow"
description: |-
  Manages a Logic App Workflow.
---

# azurerm_logic_app_workflow

Manages a Logic App Workflow.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "workflow-resources"
  location = "East US"
}

resource "azurerm_logic_app_workflow" "test" {
  name = "workflow1"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  action {
    http {
      name = "example"
      method = "GET"
      uri = "http://example.com/foo"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Logic App Workflow. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Logic App Workflow should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the Logic App Workflow exists. Changing this forces a new resource to be created.

* `action` - (Optional) A `action` block as defined below.

* `trigger` - (Optional) A `trigger` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `action` block contains:

* `http` - (Optional) An `http` block as defined below.

* `function` - (Optional) An `http` block as defined below.

---

A `trigger` block contains:

* `recurrence` - (Optional) One or more `recurrence` blocks as defined below.

---

A `http` block contains:

* `name` - (Required) The name of the HTTP Action. This needs to be unique within all Actions in the `action` block

* `method` - (Required) The HTTP Method which should be used for this action, such as `GET`, `POST` etc.

* `uri` - (Required) The HTTP URI that should be accessed via this action.

---

A `function` block contains:

* `name` - (Required) The name of the Function Action. This needs to be unique within all Actions in the `action` block

* `function_id` - (Required) The ID of the Function which should be run by this action.

* `body` - (Required) The HTTP Body to be posted to the Function.

---

A `recurrence` block contains:

* `name` - (Required) The Name of the Recurrence Trigger. This needs to be unique within all triggers in the `trigger` block.

* `frequency` - (Required) The Frequency of the Recurrence Trigger. Possible values include `Month`, `Week`, `Day`, `Hour`, `Minute` and `Second`.

* `interval` - (Required) The Interval of the Recurrence Trigger.


## Attributes Reference

The following attributes are exported:

* `id` - The Logic App Workflow ID.

* `access_endpoint` - The Access Endpoint for the Logic App Workflow


## Import

Logic App Workflows can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_workflow.workflow1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Logic/workflows/workflow1
```
