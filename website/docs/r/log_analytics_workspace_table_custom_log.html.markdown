---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_workspace_table_custom_log"
description: |-
  Manages a Custom Log Table in a Log Analytics (formally Operational Insights) Workspace.
---

# azurerm_log_analytics_workspace_table_custom_log

Manages a Custom Log Table in a Log Analytics (formally Operational Insights) Workspace.

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

resource "azurerm_log_analytics_workspace_table_custom_log" "example" {
  name         = "example_CL"
  workspace_id = azurerm_log_analytics_workspace.example.id
  display_name = "Example Custom Log"
  description  = "Custom log table for example data"
  plan         = "Analytics"

  column {
    name = "TimeGenerated"
    type = "datetime"
  }

  column {
    name = "Application"
    type = "string"
  }

  column {
    name = "RawData"
    type = "string"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Log Analytics Workspace Table Custom Log. The table name must end with `_CL` for Custom Log tables. Changing this forces a new resource to be created.

* `workspace_id` - (Required) The object ID of the Log Analytics Workspace that contains the table. Changing this forces a new resource to be created.

---

* `column` - (Required) One or more `column` blocks as defined below.

* `description` - (Optional) The description of the table.

* `display_name` - (Optional) The display name of the table.

* `labels` - (Optional) Specifies a list of labels to assign to the table.

* `plan` - (Optional) Specify the system how to handle and charge the logs ingested to the table. Possible values are `Analytics` and `Basic`. Defaults to `Analytics`.

~> **Note:** Changing the table's `plan` is limited to once a week.

* `retention_in_days` - (Optional) The table's retention in days. Possible values range between `4` and `730`.

~> **Note:** `retention_in_days` cannot be set when `plan` is set to `Basic` because the retention is fixed at eight days on the Basic plan.

* `total_retention_in_days` - (Optional) The table's total retention in days. Possible values range between `4` and `730`; or `1095`, `1460`, `1826`, `2191`, `2556`, `2922`, `3288`, `3653`, `4018`, or `4383`.

---

A `column` block supports the following:

* `name` - (Required) Specifies the name of the column.

* `type` - (Required) The data type of the column. Possible values are `boolean`, `datetime`, `dynamic`, `guid`, `int`, `long`, `real`, and `string`.

* `description` - (Optional) The description of the column.

* `display_by_default` - (Optional) Specifies whether the column should be displayed by default in the query results. Defaults to `true`.

* `display_name` - (Optional) The display name of the column.

* `hidden` - (Optional) Specifies whether the column is hidden in the query results. Defaults to `false`.

* `type_hint` - (Optional) The type hint for the column data type. Can only be set for columns of type `string`. Possible values are `armPath`, `guid`, `ip`, and `uri`.

~> **Note:** `type_hint` can only be set for columns of type `string`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Log Analytics Workspace Table Custom Log.

* `solutions` - A list of solutions associated with the table.

* `standard_column` - One or more `standard_column` blocks as defined below.

---

A `standard_column` block exports the following:

* `description` - The description of the standard column.

* `display_by_default` - Whether the standard column is displayed by default. Defaults to `true`.

* `display_name` - The display name of the standard column.

* `hidden` - Whether the standard column is hidden.

* `name` - (Required) The name of the standard column.

* `type` - (Required) The data type of the standard column.

* `type_hint` - The type hint of the standard column.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics Workspace Table Custom Log.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Workspace Table Custom Log.
* `update` - (Defaults to 30 minutes) Used when updating the Log Analytics Workspace Table Custom Log.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Workspace Table Custom Log.

## Import

Log Analytics Workspace Table Custom Logs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_workspace_table_custom_log.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/tables/table1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.OperationalInsights` - 2022-10-01
