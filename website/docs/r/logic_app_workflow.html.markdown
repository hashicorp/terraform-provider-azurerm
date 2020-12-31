---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_workflow"
description: |-
  Manages a Logic App Workflow.
---

# azurerm_logic_app_workflow

Manages a Logic App Workflow.

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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Logic App Workflow. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Logic App Workflow should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the Logic App Workflow exists. Changing this forces a new resource to be created.

* `integration_service_environment_id` - (Optional) The ID of the Integration Service Environment to which this Logic App Workflow belongs.  Changing this forces a new Logic App Workflow to be created.

* `logic_app_integration_account_id` - (Optional) The ID of the integration account linked by this Logic App Workflow.

* `workflow_schema` - (Optional) Specifies the Schema to use for this Logic App Workflow. Defaults to `https://schema.management.azure.com/providers/Microsoft.Logic/schemas/2016-06-01/workflowdefinition.json#`. Changing this forces a new resource to be created.

* `workflow_version` - (Optional) Specifies the version of the Schema used for this Logic App Workflow. Defaults to `1.0.0.0`. Changing this forces a new resource to be created.

* `parameters` - (Optional) A map of Key-Value pairs.

-> **NOTE:** Any parameters specified must exist in the Schema defined in `workflow_schema`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Logic App Workflow ID.

* `access_endpoint` - The Access Endpoint for the Logic App Workflow.

* `connector_endpoint_ip_addresses` - The list of access endpoint ip addresses of connector.

* `connector_outbound_ip_addresses` - The list of outgoing ip addresses of connector.

* `workflow_endpoint_ip_addresses` - The list of access endpoint ip addresses of workflow.

* `workflow_outbound_ip_addresses` - The list of outgoing ip addresses of workflow.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Logic App Workflow.
* `update` - (Defaults to 30 minutes) Used when updating the Logic App Workflow.
* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Workflow.
* `delete` - (Defaults to 30 minutes) Used when deleting the Logic App Workflow.

## Import

Logic App Workflows can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_workflow.workflow1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Logic/workflows/workflow1
```
