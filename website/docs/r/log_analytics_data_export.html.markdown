---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_data_export_rule"
description: |-
  Manages a log analytics Data Export Rule.
---

# azurerm_log_analytics_data_export

Manages a Log Analytics Data Export Rule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "exampleworkspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoracc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_log_analytics_data_export_rule" "example" {
  name                    = "dataExport1"
  resource_group_name     = azurerm_resource_group.example.name
  workspace_resource_id   = azurerm_log_analytics_workspace.example.id
  destination_resource_id = azurerm_storage_account.example.id
  table_names             = ["Heartbeat"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Log Analytics Data Export Rule. Changing this forces a new Log Analytics Data Export Rule to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Log Analytics Data Export should exist. Changing this forces a new Log Analytics Data Export Rule to be created.

* `workspace_resource_id` - (Required) The resource ID of the workspace. Changing this forces a new Log Analytics Data Export Rule to be created.

* `destination_resource_id` - (Required) The destination resource ID. It should be a storage account, an event hub namespace or an event hub. If the destination is an event hub namespace, an event hub would be created for each table automatically.

* `table_names` - (Required) A list of table names to export to the destination resource, for example: ["Heartbeat", "SecurityEvent"].

* `enabled` - (Optional) Is this Log Analytics Data Export Rule when enabled? Possible values include `true` or `false`. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Log Analytics Data Export Rule.

* `export_rule_id` - The ID of the created Data Export Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics Data Export Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Data Export Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Log Analytics Data Export Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Data Export Rule.

## Import

Log Analytics Data Export Rule can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_data_export_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/dataExports/dataExport1
```