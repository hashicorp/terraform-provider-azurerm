---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_datasource_windows_event"
description: |-
  Manages a Log Analytics Windows Event DataSource.
---

# azurerm_log_analytics_datasource_windows_event

Manages a Log Analytics Windows Event DataSource.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-law"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_log_analytics_datasource_windows_event" "example" {
  name                = "example-lad-wpc"
  resource_group_name = azurerm_resource_group.example.name
  workspace_name      = azurerm_log_analytics_workspace.example.name
  event_log_name      = "Application"
  event_types         = ["error"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Log Analytics Windows Event DataSource. Changing this forces a new Log Analytics Windows Event DataSource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Log Analytics Windows Event DataSource should exist. Changing this forces a new Log Analytics Windows Event DataSource to be created.

* `workspace_name` - (Required) The name of the Log Analytics Workspace where the Log Analytics Windows Event DataSource should exist. Changing this forces a new Log Analytics Windows Event DataSource to be created.

* `event_log_name` - (Required) Specifies the name of the Windows Event Log to collect events from.

* `event_types` - (Required) Specifies an array of event types applied to the specified event log. Possible values include `error`, `warning` and `information`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Log Analytics Windows Event DataSource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics Windows Event DataSource.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Windows Event DataSource.
* `update` - (Defaults to 30 minutes) Used when updating the Log Analytics Windows Event DataSource.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Windows Event DataSource.

## Import

Log Analytics Windows Event DataSources can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_datasource_windows_event.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/datasources/datasource1
```
