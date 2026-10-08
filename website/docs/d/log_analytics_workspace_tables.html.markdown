---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_log_analytics_workspace_tables"
description: |-
  Gets information about tables within an existing Log Analytics Workspace.
---

# Data Source: azurerm_log_analytics_workspace_tables

Gets information about tables within an existing Log Analytics Workspace.

## Example Usage

```hcl
data "azurerm_log_analytics_workspace" "example" {
  name                = "existing-log-analytics-workspace"
  resource_group_name = "existing-resource-group"
}

data "azurerm_log_analytics_workspace_tables" "example" {
  workspace_id = data.azurerm_log_analytics_workspace.example.id
}

output "table_names" {
  value = data.azurerm_log_analytics_workspace_tables.example.names
}
```

## Arguments Reference

The following arguments are supported:

* `workspace_id` - (Required) The ID of the Log Analytics Workspace from which to retrieve table information.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Log Analytics Workspace.

* `names` - A list containing names of tables that exist in the Log Analytics Workspace.

* `tables` - A list of `tables` blocks as defined below.

---

A `tables` block exports the following:

* `name` - The name of the table in the Log Analytics Workspace.

* `plan` - The billing plan for the table.

* `retention_in_days` - The table's data retention in days.

* `total_retention_in_days` - The table's total data retention in days.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Workspace Tables.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.OperationalInsights` - 2022-10-01
