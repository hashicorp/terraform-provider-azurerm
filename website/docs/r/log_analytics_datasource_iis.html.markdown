---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_datasource_iis"
description: |-
  Manages a Log Analytics IIS DataSource.
---

# azurerm_log_analytics_datasource_iis

Manages a Log Analytics IIS DataSource.

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

resource "azurerm_log_analytics_datasource_iis" "example" {
  name                = "example-lad-iis"
  resource_group_name = azurerm_resource_group.example.name
  workspace_name      = azurerm_log_analytics_workspace.example.name
  on_premise_enabled  = true
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Log Analytics IIS DataSource. Changing this forces a new Log Analytics IIS DataSource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Log Analytics IIS DataSource should exist. Changing this forces a new Log Analytics IIS DataSource to be created.

* `workspace_name` - (Required) The name of the Log Analytics Workspace where the Log Analytics IIS DataSource should exist. Changing this forces a new Log Analytics Windows Event DataSource to be created.

* `on_premise_enabled` - (Required) Should collecting IIS log files be enabled?

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Log Analytics IIS DataSource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics IIS DataSource.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics IIS DataSource.
* `update` - (Defaults to 30 minutes) Used when updating the Log Analytics IIS DataSource.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics IIS DataSource.

## Import

Log Analytics IIS DataSources can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_datasource_iis.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/datasources/datasource1
```
