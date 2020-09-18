---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_datasource_linux_performance_collection"
description: |-
  Manages a Log Analytics Linux Performance Collection DataSource.
---

# azurerm_log_analytics_datasource_linux_performance_collection

Manages a Log Analytics Linux Performance Collection DataSource.

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

resource "azurerm_log_analytics_datasource_linux_performance_collection" "example" {
  name                = "example-lad-wpc"
  resource_group_name = azurerm_resource_group.example.name
  workspace_name      = azurerm_log_analytics_workspace.example.name
  state               = "Enabled"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Log Analytics Linux Performance Collection DataSource. Changing this forces a new Log Analytics Linux Performance Collection DataSource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Log Analytics Linux Performance Collection DataSource should exist. Changing this forces a new Log Analytics Linux Performance Collection DataSource to be created.

* `workspace_name` - (Required) The name of the Log Analytics Workspace where the Log Analytics Linux Performance Collection DataSource should exist. Changing this forces a new Log Analytics Linux Performance Collection DataSource to be created.

* `state` - (Required) Specifies the start/stop state for collection of performance data from connected linux computers. Possible values include `Enabled` and `Disabled`.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Log Analytics Linux Performance Collection DataSource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics Linux Performance Collection DataSource.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Linux Performance Collection DataSource.
* `update` - (Defaults to 30 minutes) Used when updating the Log Analytics Linux Performance Collection DataSource.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Linux Performance Collection DataSource.

## Import

Log Analytics Linux Performance Collection DataSources can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_datasource_linux_performance_collection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/datasources/datasource1
```
