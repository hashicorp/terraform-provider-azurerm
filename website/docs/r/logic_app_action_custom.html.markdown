---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_action_custom"
description: |-
  Manages a Custom Action within a Logic App Workflow
---

# azurerm_logic_app_action_custom

Manages a Custom Action within a Logic App Workflow

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "workflow-resources"
  location = "West Europe"
}

resource "azurerm_logic_app_workflow" "example" {
  name                = "workflow1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_logic_app_action_custom" "example" {
  name         = "example-action"
  logic_app_id = azurerm_logic_app_workflow.example.id

  body = <<BODY
{
    "description": "A variable to configure the auto expiration age in days. Configured in negative number. Default is -30 (30 days old).",
    "inputs": {
        "variables": [
            {
                "name": "ExpirationAgeInDays",
                "type": "Integer",
                "value": -30
            }
        ]
    },
    "runAfter": {},
    "type": "InitializeVariable"
}
BODY

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the HTTP Action to be created within the Logic App Workflow. Changing this forces a new resource to be created.

-> **NOTE:** This name must be unique across all Actions within the Logic App Workflow.

* `logic_app_id` - (Required) Specifies the ID of the Logic App Workflow. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the JSON Blob defining the Body of this Custom Action.

-> **NOTE:** To make the Action more readable, you may wish to consider using HEREDOC syntax (as shown above) or [the `local_file` resource](https://www.terraform.io/docs/providers/local/d/file.html) to load the schema from a file on disk.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Action within the Logic App Workflow.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Logic App Custom Action.
* `update` - (Defaults to 30 minutes) Used when updating the Logic App Custom Action.
* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Custom Action.
* `delete` - (Defaults to 30 minutes) Used when deleting the Logic App Custom Action.

## Import

Logic App Custom Actions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_action_custom.custom1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Logic/workflows/workflow1/actions/custom1
```

-> **NOTE:** This ID is unique to Terraform and doesn't directly match to any other resource. To compose this ID, you can take the ID Logic App Workflow and append `/actions/{name of the action}`.
