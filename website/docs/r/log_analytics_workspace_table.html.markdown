---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_workspace_table"
description: |-
  Manages a Table in a Log Analytics (formally Operational Insights) Workspace.
---

# azurerm_log_analytics_workspace_table

Manages a Table in a Log Analytics (formally Operational Insights) Workspace.

~> **Note:** This resource does not create or destroy tables. This resource is used to update attributes of the tables created when a Log Analytics Workspace is created. Deleting an `azurerm_log_analytics_workspace_table` resource does not delete the table. Instead, the table's `retention_in_days` is reset to inherit the workspace retention, and `total_retention_in_days` is reset to the table's retention.

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
  name                    = "AppMetrics"
  workspace_id            = azurerm_log_analytics_workspace.example.id
  retention_in_days       = 60
  total_retention_in_days = 180
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of a table in a Log Analytics Workspace.

* `workspace_id` - (Required) The resource ID of the Log Analytics Workspace that contains the table.

* `plan` - (Optional) Specify the system how to handle and charge the logs ingested to the table. Possible values are `Analytics` and `Basic`. Defaults to `Analytics`.

-> **Note:** The `name` of tables currently supported by the `Basic` plan can be found [here](https://learn.microsoft.com/azure/azure-monitor/logs/basic-logs-azure-tables).

* `retention_in_days` - (Optional) The number of days that data is retained for interactive queries (analytics retention). Possible values are between `4` and `730`.

-> **Note:** `retention_in_days` is fixed at `30` days when `plan` is `Basic`. See more on [data retention](https://learn.microsoft.com/azure/azure-monitor/logs/data-retention-configure#analytics-long-term-and-total-retention).

-> **Note:** When `retention_in_days` is omitted for an `Analytics` table, the table inherits the workspace's default retention setting, subject to table-specific defaults. See more on [configuring the default analytics retention period](https://learn.microsoft.com/azure/azure-monitor/logs/data-retention-configure#configure-the-default-analytics-retention-period-of-analytics-tables).

* `total_retention_in_days` - (Optional) The total number of days that data is retained, including analytics retention and any additional long-term retention. Possible values range between `4` and `730`; or `1095`, `1460`, `1826`, `2191`, `2556`, `2922`, `3288`, `3653`, `4018`, or `4383`.

-> **Note:** When `total_retention_in_days` is omitted, the total retention equals the table's analytics retention period, with no additional long-term retention. See more on [analytics, long-term, and total retention](https://learn.microsoft.com/azure/azure-monitor/logs/data-retention-configure#analytics-long-term-and-total-retention).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Log Analytics Workspace Table ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Log Analytics Workspace Table.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Workspace Table.
* `update` - (Defaults to 5 minutes) Used when updating the Log Analytics Workspace Table.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Workspace Table.

## Import

Log Analytics Workspace Tables can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_workspace_table.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/tables/AppMetrics
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.OperationalInsights` - 2022-10-01
