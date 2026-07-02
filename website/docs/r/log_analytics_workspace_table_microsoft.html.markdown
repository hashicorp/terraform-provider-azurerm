---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_workspace_table_microsoft"
description: |-
  Manages a Microsoft Table in a Log Analytics (formally Operational Insights) Workspace.
---

# azurerm_log_analytics_workspace_table_microsoft

Manages a Microsoft Table in a Log Analytics (formally Operational Insights) Workspace.

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
  retention_in_days   = 30
}

resource "azurerm_log_analytics_workspace_table_microsoft" "example" {
  name         = "AppCenterError"
  workspace_id = azurerm_log_analytics_workspace.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Log Analytics Workspace Table Microsoft. Possible values are `Alert`, `AppCenterError`, `ComputerGroup`, `InsightsMetrics`, `Operation` and `Usage`. Changing this forces a new Log Analytics Workspace Table Microsoft to be created.

* `workspace_id` - (Required) The ID of the Log Analytics Workspace. Changing this forces a new Log Analytics Workspace Table Microsoft to be created.

---

* `column` - (Optional) One or more `column` blocks as defined below.

* `description` - (Optional) A description of the table.

* `display_name` - (Optional) The display name of the table.

* `labels` - (Optional) Specifies a list of table labels.

* `retention_in_days` - (Optional) The table retention in days, between `4` and `730`. 

-> **Note:** `retention_in_days` must be less than or equal to `total_retention_in_days`.

* `total_retention_in_days` - (Optional) The table total retention in days, between `4` and `4383`. 

---

A `column` block supports the following:

* `name` - (Required) The name which should be used for this column.

* `type` - (Required) The column data type. Possible values are `string`, `int`, `long`, `real`, `boolean`, `dateTime`, `guid`, `dynamic`.

* `description` - (Optional) The description of the column.

* `display_by_default` - (Optional) Whether the column defaults to being displayed. Defaults to `true`.

* `display_name` - (Optional) The display name of the column.

* `hidden` - (Optional) Whether the column is hidden. Defaults to `false`.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Log Analytics Workspace Table Microsoft.

* `solutions` - The list of solutions associated with this table.

* `standard_column` - A `standard_column` block as defined below.

---

A `standard_column` block exports the following attributes:

* `description` - The description of the column.

* `display_by_default` - Whether the column defaults to being displayed. Defaults to `true`.

* `display_name` - The display name of the column.

* `hidden` - Is the column hidden? Defaults to `false`.

* `name` - The name of the column.

* `type` -  The type of the column.

* `type_hint` - The type hint of the column.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Log Analytics Workspace Table Microsoft.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Workspace Table Microsoft.
* `update` - (Defaults to 5 minutes) Used when updating the Log Analytics Workspace Table Microsoft.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Workspace Table Microsoft.

## Import

Log Analytics Workspace Table Microsofts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_workspace_table_microsoft.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/tables/table1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.OperationalInsights` - 2022-10-01
