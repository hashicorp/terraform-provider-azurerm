---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_log_analytics_workspace_tables"
description: |-
  Gets all tables from an existing Log Analytics Workspace.
---

# Data Source: azurerm_log_analytics_workspace_tables

Use this data source to access information about Tables within an existing Log Analytics Workspace.

## Example Usage

```hcl
data "azurerm_log_analytics_workspace_tables" "example" {
  workspace_id = azurerm_log_analytics_workspace.example.id
}

output "table_names" {
  value = [for table in data.azurerm_log_analytics_workspace_tables.example.tables : table.name]
}
```

## Arguments Reference

The following arguments are supported:

* `workspace_id` - (Required) The ID of the Log Analytics Workspace from which to retrieve table information.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Log Analytics Workspace.

* `tables` - A list of `tables` blocks as defined below.

---

Each element in `tables` block exports the following:

* `name` - The name of the table in the Log Analytics Workspace.

* `plan` - The plan type for the table. Possible values are `Analytics` and `Basic`.

* `retention_in_days` - The table's data retention in days.

* `total_retention_in_days` - The table's total data retention in days.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Workspace Tables.
