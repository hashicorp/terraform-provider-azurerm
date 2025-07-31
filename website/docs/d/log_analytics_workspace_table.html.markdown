---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_workspace_table"
description: |-
  Gets information about an existing Log Analytics Workspace Table.
---

# Data Source: azurerm_log_analytics_workspace_table

Use this data source to access information about an existing Log Analytics Workspace Table.

## Example Usage

```hcl
data "azurerm_log_analytics_workspace_table" "example" {
  name         = "InsightsMetrics"
  workspace_id = azurerm_log_analytics_workspace.example.id
}

output "retention_in_days" {
  value = data.azurerm_log_analytics_workspace_table.example.retention_in_days
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Log Analytics Workspace Table.

* `workspace_id` - (Required) The ID of the Log Analytics Workspace the table belongs to.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `retention_in_days` - The table's data retention in days.

* `total_retention_in_days` - The table's total data retention in days.

* `plan` - The billing plan information for the Log Analytics Workspace Table.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Workspace Table.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.OperationalInsights`: 2022-10-01
