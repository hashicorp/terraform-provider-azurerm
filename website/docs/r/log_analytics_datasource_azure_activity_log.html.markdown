---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_datasource_azure_activity_log"
description: |-
  Manages a Log Analytics Azure Activity Log DataSource.
---

# azurerm_log_analytics_datasource_azure_activity_log

Manages a Log Analytics Azure Activity Log DataSource.

## Example Usage

```hcl
resource "azurerm_log_analytics_datasource_azure_activity_log" "example" {
  name                = "example"
  resource_group_name = "example"
  workspace_name      = "example"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Log Analytics Azure Activity Log DataSource. Changing this forces a new Log Analytics Azure Activity Log DataSource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Log Analytics Azure Activity Log DataSource should exist. Changing this forces a new Log Analytics Azure Activity Log DataSource to be created.

* `workspace_name` - (Required) The name of the Log Analytics Workspace where the Log Analytics Azure Activity Log DataSource should exist. Changing this forces a new Log Analytics Azure Activity Log DataSource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Log Analytics Azure Activity Log DataSource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics Azure Activity Log DataSource.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Azure Activity Log DataSource.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Azure Activity Log DataSource.

## Import

Log Analytics Azure Activity Log DataSources can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_data_source_azure_activity_log.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/datasources/datasource1
```
