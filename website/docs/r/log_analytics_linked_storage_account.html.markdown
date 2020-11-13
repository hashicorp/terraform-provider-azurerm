---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_linked_storage_account"
description: |-
  Manages a Log Analytics Linked Storage Account.
---

# azurerm_log_analytics_linked_storage_account

Manages a Log Analytics Linked Storage Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "exampleworkspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_log_analytics_linked_storage_account" "example" {
  data_source_type      = "customlogs"
  resource_group_name   = azurerm_resource_group.example.name
  workspace_resource_id = azurerm_log_analytics_workspace.example.id
  storage_account_ids   = [azurerm_storage_account.example.id]
}
```

## Arguments Reference

The following arguments are supported:

* `data_source_type` - (Required) The data source type which should be used for this Log Analytics Linked Storage Account. Possible values are "customlogs", "azurewatson", "query", "ingestion" and "alerts". Changing this forces a new Log Analytics Linked Storage Account to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Log Analytics Linked Storage Account should exist. Changing this forces a new Log Analytics Linked Storage Account to be created.

* `workspace_resource_id` - (Required) The resource ID of the Log Analytics Workspace. Changing this forces a new Log Analytics Linked Storage Account to be created.

* `storage_account_ids` - (Required) The storage account resource ids to be linked.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Log Analytics Linked Storage Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics Linked Storage Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Linked Storage Account.
* `update` - (Defaults to 30 minutes) Used when updating the Log Analytics Linked Storage Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Linked Storage Account.

## Import

Log Analytics Linked Storage Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_linked_storage_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/linkedStorageAccounts/{dataSourceType}
```
