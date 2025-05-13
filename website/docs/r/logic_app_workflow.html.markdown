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
  location = "West Europe"
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

* `access_control` - (Optional) A `access_control` block as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `integration_service_environment_id` - (Optional) The ID of the Integration Service Environment to which this Logic App Workflow belongs. Changing this forces a new Logic App Workflow to be created.

* `logic_app_integration_account_id` - (Optional) The ID of the integration account linked by this Logic App Workflow.

* `enabled` - (Optional) Is the Logic App Workflow enabled? Defaults to `true`.

* `workflow_parameters` - (Optional) Specifies a map of Key-Value pairs of the Parameter Definitions to use for this Logic App Workflow. The key is the parameter name, and the value is a JSON encoded string of the parameter definition (see: <https://docs.microsoft.com/azure/logic-apps/logic-apps-workflow-definition-language#parameters>).
  
* `workflow_schema` - (Optional) Specifies the Schema to use for this Logic App Workflow. Defaults to `https://schema.management.azure.com/providers/Microsoft.Logic/schemas/2016-06-01/workflowdefinition.json#`. Changing this forces a new resource to be created.

* `workflow_version` - (Optional) Specifies the version of the Schema used for this Logic App Workflow. Defaults to `1.0.0.0`. Changing this forces a new resource to be created.

* `parameters` - (Optional) A map of Key-Value pairs.

-> **Note:** Any parameters specified must exist in the Schema defined in `workflow_parameters`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `access_control` block supports the following:

* `action` - (Optional) A `action` block as defined below.

* `content` - (Optional) A `content` block as defined below.

* `trigger` - (Optional) A `trigger` block as defined below.

* `workflow_management` - (Optional) A `workflow_management` block as defined below.

---

A `action` block supports the following:

* `allowed_caller_ip_address_range` - (Required) A list of the allowed caller IP address ranges.

---

A `content` block supports the following:

* `allowed_caller_ip_address_range` - (Required) A list of the allowed caller IP address ranges.

---

A `trigger` block supports the following:

* `allowed_caller_ip_address_range` - (Required) A list of the allowed caller IP address ranges.

* `open_authentication_policy` - (Optional) A `open_authentication_policy` block as defined below.

---

A `workflow_management` block supports the following:

* `allowed_caller_ip_address_range` - (Required) A list of the allowed caller IP address ranges.

---

A `open_authentication_policy` block supports the following:

* `name` - (Required) The OAuth policy name for the Logic App Workflow.

* `claim` - (Required) A `claim` block as defined below.

---

A `claim` block supports the following:

* `name` - (Required) The name of the OAuth policy claim for the Logic App Workflow.

* `value` - (Required) The value of the OAuth policy claim for the Logic App Workflow.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Logic App Workflow. Possible values are `SystemAssigned`, `UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Logic App Workflow.

~> **Note:** This is required when `type` is set to `UserAssigned`

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Logic App Workflow ID.

* `access_endpoint` - The Access Endpoint for the Logic App Workflow.

* `connector_endpoint_ip_addresses` - The list of access endpoint IP addresses of connector.

* `connector_outbound_ip_addresses` - The list of outgoing IP addresses of connector.

* `identity` - An `identity` block as defined below.

* `workflow_endpoint_ip_addresses` - The list of access endpoint IP addresses of workflow.

* `workflow_outbound_ip_addresses` - The list of outgoing IP addresses of workflow.

---

The `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this Logic App Workflow.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this Logic App Workflow.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Logic App Workflow.
* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Workflow.
* `update` - (Defaults to 30 minutes) Used when updating the Logic App Workflow.
* `delete` - (Defaults to 30 minutes) Used when deleting the Logic App Workflow.

## Import

Logic App Workflows can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_workflow.workflow1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Logic/workflows/workflow1
```
