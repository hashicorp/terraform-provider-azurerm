---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_workflow"
description: |-
  Gets information about an existing Logic App Workflow.
---

# Data Source: azurerm_logic_app_workflow

Use this data source to access information about an existing Logic App Workflow.

## Example Usage

```hcl
data "azurerm_logic_app_workflow" "example" {
  name                = "workflow1"
  resource_group_name = "my-resource-group"
}

output "access_endpoint" {
  value = data.azurerm_logic_app_workflow.example.access_endpoint
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Logic App Workflow.

* `resource_group_name` - The name of the Resource Group in which the Logic App Workflow exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Logic App Workflow ID.

* `location` - The Azure location where the Logic App Workflow exists.

* `logic_app_integration_account_id` - The ID of the integration account linked by this Logic App Workflow.

* `workflow_schema` - The Schema used for this Logic App Workflow.

* `workflow_version` - The version of the Schema used for this Logic App Workflow. Defaults to `1.0.0.0`.

* `parameters` - A map of Key-Value pairs.

* `tags` - A mapping of tags assigned to the resource.

* `access_endpoint` - The Access Endpoint for the Logic App Workflow

* `connector_endpoint_ip_addresses` - The list of access endpoint IP addresses of connector.

* `connector_outbound_ip_addresses` - The list of outgoing IP addresses of connector.

* `workflow_endpoint_ip_addresses` - The list of access endpoint IP addresses of workflow.

* `workflow_outbound_ip_addresses` - The list of outgoing IP addresses of workflow.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Logic App Workflow.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Logic App Workflow.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Logic App Workflow.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Logic App Workflow.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Workflow.
