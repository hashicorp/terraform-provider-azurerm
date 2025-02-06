---
subcategory: "Databricks"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_databricks_workspace_private_endpoint_connection"
description: |-
  Gets information on an existing Databricks Workspace private endpoint connection state
---

# Data Source: azurerm_databricks_workspace_private_endpoint_connection

Use this data source to access information on an existing Databricks Workspace private endpoint connection state.

## Example Usage

```hcl
data "azurerm_databricks_workspace_private_endpoint_connection" "example" {
  workspace_id        = azurerm_databricks_workspace.example.id
  private_endpoint_id = azurerm_private_endpoint.example.id
}

output "databricks_workspace_private_endpoint_connection_status" {
  value = data.azurerm_databricks_workspace_private_endpoint_connection.example.connections[0].status
}
```

## Argument Reference

* `name` - The name of the Databricks Workspace.
* `resource_group_name` - The Name of the Resource Group where the Databricks Workspace exists.

## Attributes Reference

* `workspace_id` - The resource ID of the Databricks Workspace.

* `private_endpoint_id` - The resource ID of the Private Endpoint.

* `connections` - A `connections` block as documented below.

---

A `connections` block exports the following:

* `workspace_private_endpoint_id` - The Databricks Workspace resource ID for the private link endpoint.

* `name` - The name of the private endpoint connection.

* `status` - The status of a private endpoint connection. Possible values are `Pending`, `Approved`, `Rejected` or `Disconnected`.

* `description` - The description for the current state of a private endpoint connection.

* `action_required` - Actions required for a private endpoint connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Databricks Workspace Private Endpoint Connection.
