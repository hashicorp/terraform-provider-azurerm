---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_workspace_table"
description: |-
  Manages a Table in a Log Analytics (formally Operational Insights) Workspace.
---

# azurerm_log_analytics_workspace_table

Manages a Table in a Log Analytics (formally Operational Insights) Workspace.

-> **Note:** Only `CustomLog` tables with a `DataCollectionRuleBased` sub-type are supported for creation and deletion with this resource. For `Microsoft` and `Classic` tables, they can be imported and managed, but if removed from the Terraform configuration they will have their retention policy set to the workspace default and be removed from state.

## Example Usage

Update the retention policy on an inbuilt table

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_log_analytics_workspace" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}
resource "azurerm_log_analytics_workspace_table" "example" {
  workspace_id            = azurerm_log_analytics_workspace.example.id
  name                    = "AppMetrics"
  type                    = "Microsoft"
  sub_type                = "DataCollectionRuleBased"
  retention_in_days       = 60
  total_retention_in_days = 180
}
```

Create a new Custom Log Table

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_log_analytics_workspace" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}
resource "azurerm_log_analytics_workspace_table" "example" {
  name         = "CustomTable_CL"
  type         = "CustomLog"
  sub_type     = "DataCollectionRuleBased"
  workspace_id = azurerm_log_analytics_workspace.example.id

  column {
    name = "ServiceName"
    type = "string"
  }
  column {
    name = "TimeGenerated"
    type = "dateTime"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Specifies the name of the table in a Log Analytics Workspace. Must end in `_CL` for custom tables.

- `workspace_id` - (Required) The object ID of the Log Analytics Workspace that will contain the table.

- `type` - (Required) The type of table. Must be either of `Microsoft` for inbuilt tables, or `CustomLog` for custom tables.

- `sub_type` - (Required) The sub type of table. Must be one of `Any`, `Classic`, or `DataCollectionRuleBased`.

- `display_name` - (Optional) The display name of the table in a Log Analytics Workspace.

- `description` - (Optional) The description of the table in a Log Analytics Workspace.

- `categories` - (Optional) The categories applied to the table.

- `column` - (Optional) One or more `column` blocks detailed below.

-> **Note:** The order of the columns will match the display order in Log Analytics.

- `labels` - (Optional) The labels applied to the table.

- `plan` - (Optional) Specify the system how to handle and charge the logs ingested to the table. Possible values are `Analytics` and `Basic`. Defaults to `Analytics`.

-> **Note:** The `name` of tables currently supported by the `Basic` plan can be found [here](https://learn.microsoft.com/en-us/azure/azure-monitor/logs/basic-logs-configure?tabs=portal-1#supported-tables).

- `retention_in_days` - (Optional) The table's retention in days. Possible values are either `8` (Basic Tier only) or range between `4` and `730`.

- `total_retention_in_days` - (Optional) The table's total retention in days. Possible values range between `4` and `730`; or `1095`, `1460`, `1826`, `2191`, `2556`, `2922`, `3288`, `3653`, `4018`, or `4383`.

-> **Note:** `retention_in_days` and `total_retention_in_days` will revert back to the value of azurerm_log_analytics_workspace retention_in_days when a Microsoft or Classic azurerm_log_analytics_workspace_table is deleted.

-> **Note:** The `retention_in_days` cannot be specified when `plan` is `Basic` because the retention is fixed at eight days.

---

A `column` block supports the following:

- `name` - (Required) The name of the column.

- `type` - (Required) The type of data stored in the column. Must be one of `boolean`, `dateTime`, `dynamic`, `guid"`, `int`, `long`, `real` , or `string`.

- `display_name` - (Optional) The display name of the column.

- `description` - (Optional) A description of the column.

- `display_by_default` - (Optional) Is the column displayed by default. Defaults to `true`.

- `type_hint` - (Optional) A hint as to what kind of data is stored in a `string` column. Must be one of `armpath`, `guid`, `ip`, or `uri`.

- `hidden` - (Optional) Is the column hidden? Defaults to `false`.

## Attributes Reference

The following attributes are exported:

- `id` - The Log Analytics Workspace Table ID.

- `workspace_id` - The Workspace (or Customer) ID for the Log Analytics Workspace.

- `retention_in_days` - The table's data retention in days.

- `solutions` - The list of solutions associated with this table.

- `standard_column` - The details of the standard columns in this table.

- `total_retention_in_days` - The table's total data retention in days.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

- `create` - (Defaults to 5 minutes) Used when creating the Log Analytics Workspace.
- `update` - (Defaults to 5 minutes) Used when updating the Log Analytics Workspace.
- `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Workspace.
- `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Workspace.
