---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_workspace_table"
description: |-
  Manages a Table in a Log Analytics (formally Operational Insights) Workspace.
---

# azurerm_log_analytics_workspace_table

Manages a Table in a Log Analytics (formally Operational Insights) Workspace.

~> **Note:** This resource does not create or destroy tables. This resource is used to update attributes (currently only retention_in_days) of the tables created when a Log Analytics Workspace is created. Deleting an azurerm_log_analytics_workspace_table resource will not delete the table. Instead, the table's retention_in_days field will be set to the value of azurerm_log_analytics_workspace retention_in_days

## Example Usage

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
  retention_in_days       = 60
  total_retention_in_days = 180
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of a table in a Log Analytics Workspace.

* `workspace_id` - (Required) The object ID of the Log Analytics Workspace that contains the table.

* `plan` - (Optional) Specify the system how to handle and charge the logs ingested to the table. Possible values are `Analytics` and `Basic`. Defaults to `Analytics`.

-> **Note:** The `name` of tables currently supported by the `Basic` plan can be found [here](https://learn.microsoft.com/en-us/azure/azure-monitor/logs/basic-logs-azure-tables).

* `retention_in_days` - (Optional) The table's retention in days. Possible values are either `8` (Basic Tier only) or range between `4` and `730`.

* `total_retention_in_days` - (Optional) The table's total retention in days. Possible values range between `4` and `730`; or `1095`, `1460`, `1826`, `2191`, `2556`, `2922`, `3288`, `3653`, `4018`, or `4383`.

-> **Note:** `retention_in_days` and `total_retention_in_days` will revert back to the value of azurerm_log_analytics_workspace retention_in_days when a azurerm_log_analytics_workspace_table is deleted.

-> **Note:** The `retention_in_days` cannot be specified when `plan` is `Basic` because the retention is fixed at eight days.

## Attributes Reference

The following attributes are exported:

* `id` - The Log Analytics Workspace Table ID.

* `workspace_id` - The Workspace (or Customer) ID for the Log Analytics Workspace.

* `retention_in_days` - The table's data retention in days.

* `total_retention_in_days` - The table's total data retention in days.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Log Analytics Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Workspace.
* `update` - (Defaults to 5 minutes) Used when updating the Log Analytics Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Workspace.
