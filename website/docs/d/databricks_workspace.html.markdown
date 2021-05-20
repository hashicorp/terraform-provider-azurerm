---
subcategory: "Databricks"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_databricks_workspace"
description: |-
  Gets information on an existing Databricks Workspace
---

# Data Source: azurerm_databricks_workspace

Use this data source to access information about an existing Databricks workspace.

## Example Usage

```hcl
data "azurerm_databricks_workspace" "example" {
  name                = "example-workspace"
  resource_group_name = "example-rg"
}

output "databricks_workspace_id" {
  value = data.azurerm_databricks_workspace.example.workspace_id
}
```

## Argument Reference

* `name` - The name of the Databricks Workspace.
* `resource_group_name` - The Name of the Resource Group where the Databricks Workspace exists.

## Attributes Reference

* `id` - The ID of the Databricks Workspace.

* `location` - The Azure location where the Databricks Workspace exists.

* `sku` - SKU of this Databricks Workspace.

* `workspace_id` - Unique ID of this Databricks Workspace in Databricks management plane.

* `workspace_url` - URL this Databricks Workspace is accessible on.

* `tags` - A mapping of tags to assign to the Databricks Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Databricks Workspace.
