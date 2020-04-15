---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_datasource_windows_performance_counter"
description: |-
  Manages a Log Analytics (formally Operational Insights) Windows Performance Counter DataSource.
---

# azurerm_log_analytics_datasource_windows_performance_counter

Manages a Log Analytics (formally Operational Insights) Windows Performance Counter DataSource.

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

resource "azurerm_log_analytics_datasource_windows_performance_counter" "example" {
  name                = "example-lad-wpc"
  resource_group_name = azurerm_resource_group.example.name
  workspace_name      = azurerm_log_analytics_workspace.example.name
  object_name         = "CPU"
  instance_name       = "*"
  counter_name        = "CPU"
  interval_seconds    = 10
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The Name which should be used for this Log Analytics Windows Performance Counter DataSource. Changing this forces a new Log Analytics Windows Performance Counter DataSource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Log Analytics Windows Performance Counter DataSource should exist. Changing this forces a new Log Analytics Windows Performance Counter DataSource to be created.

* `workspace_name` - (Required) The name of the Log Analytics Workspace where the Log Analytics Windows Performance Counter DataSource should exist. Changing this forces a new Log Analytics Windows Performance Counter DataSource to be created.

* `object_name` - (Required) The object name of the Log Analytics Windows Performance Counter DataSource.

* `instance_name` - (Required) The name of the virtual machine instance to which the Windows Performance Counter DataSource be applied. Specify a `*` will apply to all instances.

* `counter_name` - (Required) The friendly name of the performance counter.

* `interval_seconds` - (Required) The time of sample interval in seconds. Supports values between 10 and 2147483647.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Log Analytics Windows Performance Counter DataSource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics Windows Performance Counter DataSource.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Windows Performance Counter DataSource.
* `update` - (Defaults to 30 minutes) Used when updating the Log Analytics Windows Performance Counter DataSource.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Windows Performance Counter DataSource.

## Import

Log Analytics Windows Performance Counter DataSources can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_datasource_windows_performance_counter.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/datasources/datasource1
```
