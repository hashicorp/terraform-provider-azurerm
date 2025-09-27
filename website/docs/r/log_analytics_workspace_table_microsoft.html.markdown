---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_workspace_table_microsoft"
description: |-
  Manages a Log Analytics Workspace Table Microsoft.
---

# azurerm_log_analytics_workspace_table_microsoft

Manages a Log Analytics Workspace Table Microsoft.

## Example Usage

```hcl
resource "azurerm_log_analytics_workspace" "example" {
  name = "example"
}
resource "azurerm_log_analytics_workspace_table_microsoft" "example" {
  name         = "example"
  workspace_id = azurerm_log_analytics_workspace.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Log Analytics Workspace Table Microsoft. Changing this forces a new Log Analytics Workspace Table Microsoft to be created.

* `workspace_id` - (Required) The ID of the log analytics workspace. Changing this forces a new Log Analytics Workspace Table Microsoft to be created.

---

* `categories` - (Optional) Specifies a list of table categories.

* `column` - (Optional) One or more `column` blocks as defined below.

* `description` - (Optional) a description of the table.

* `display_name` - (Optional) the display name of the table.

* `labels` - (Optional) Specifies a list of table labels.

* `plan` - (Optional) the plan by which the logs that are ingested to this table are handled and charged. Possible values are `Basic` and `Analytics`. Defaults to `Analytics`.

* `retention_in_days` - (Optional) the table retention in days, between 4 and 730. Setting this property to -1 will default to the workspace retention.

* `sub_type` - (Optional) the API and feature subtype of the table. Possible values are `Any`, `Classic`, and `DataCollectionRuleBased`. Defaults to `Any`. Changing this forces a new resource to be created.

* `total_retention_in_days` - (Optional) the table total retention in days, between 4 and 4383. Setting this property to -1 will default to table retention.

---

A `column` block supports the following:

* `name` - (Required) The name which should be used for this column.

* `type` - (Required) the column data type. Possible values are `string`, `int`, `long`, `real`, `boolean`, `dateTime`, `guid`, `dynamic`.

* `description` - (Optional) the description of the column.

* `display_by_default` - (Optional) whether the column defaults to being displayed. Defaults to `true`.

* `display_name` - (Optional) the display name of the column.

* `hidden` - (Optional) whether the column is hidden. Defaults to `false`.

* `type_hint` - (Optional) the column data type logical hint. Possible values are `uri`, `guid`, `armPath`, `ip`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Log Analytics Workspace Table Microsoft.

* `solutions` - The list of solutions associated with this table.

* `standard_column` - A `standard_column` block as defined below.

---

A `standard_column` block exports the same properties as a `column` block:

* `description` - the description of the column.

* `display_by_default` - whether the column defaults to being displayed. Defaults to `true`.

* `display_name` - the display name of the column.

* `hidden` - (Optional) Is the column hidden? Defaults to `false`.

* `name` - (Required) The name of the column.

* `type` - (Required) the type of the column.

* `type_hint` - the type hint of the column.

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
