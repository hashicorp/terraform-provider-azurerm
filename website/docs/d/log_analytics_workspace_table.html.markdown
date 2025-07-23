---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_workspace_table"
description: |-
  Gets information about an existing log_analytics_workspace.
---

# Data Source: azurerm_log_analytics_workspace_table

Use this data source to access information about an existing log_analytics_workspace table.

## Example Usage

```hcl

data "azurerm_log_analytics_workspace_table" "this" {
  name                = "InsightsMetrics"
  workspace_id        = azurerm_log_analytics_workspace.this.id
  resource_group_name = "test-resource-group"
}

output "retention_in_days" {
  value = data.azurerm_log_analytics_workspace_table.this.retention_in_days
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this log analytics workspace table.

* `resource_group_name` - (Required) The name of the Resource Group where the log analytics workspace exists.

* `workspace_id` - (Required) The ID of the log analytics workspace the table belongs to.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `retention_in_days` - The table's data retention in days.

* `total_retention_in_days` - The table's total data retention in days.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Workspace Table.
