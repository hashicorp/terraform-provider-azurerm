---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_watcher_flow_log"
description: |-
  Manages a Network Watcher Flow Log.

---

# azurerm_network_watcher_flow_log

Manages a Network Watcher Flow Log.

~> **Note:** The `azurerm_network_watcher_flow_log` creates a new storage lifecyle management rule that overwrites existing rules. Please make sure to use a `storage_account` with no existing management rules, until the [issue](https://github.com/hashicorp/terraform-provider-azurerm/issues/6935) is fixed.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_network_watcher" "test" {
  name                = "acctestnw"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsa"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  account_tier               = "Standard"
  account_kind               = "StorageV2"
  account_replication_type   = "LRS"
  https_traffic_only_enabled = true
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestlaw"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.example.name
  name                 = "example-log"

  target_resource_id = azurerm_network_security_group.test.id
  storage_account_id = azurerm_storage_account.test.id
  enabled            = true

  retention_policy {
    enabled = true
    days    = 7
  }

  traffic_analytics {
    enabled               = true
    workspace_id          = azurerm_log_analytics_workspace.test.workspace_id
    workspace_region      = azurerm_log_analytics_workspace.test.location
    workspace_resource_id = azurerm_log_analytics_workspace.test.id
    interval_in_minutes   = 10
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Network Watcher Flow Log. Changing this forces a new resource to be created.

* `network_watcher_name` - (Required) The name of the Network Watcher. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Network Watcher was deployed. Changing this forces a new resource to be created.

* `target_resource_id` - (Required) The ID of the Resource for which to enable flow logs for. Changing this forces a new resource to be created.

* `storage_account_id` - (Required) The ID of the Storage Account where flow logs are stored.

* `enabled` - (Required) Should Network Flow Logging be Enabled?

* `retention_policy` - (Required) A `retention_policy` block as documented below.

* `location` - (Optional) The location where the Network Watcher Flow Log resides. Changing this forces a new resource to be created. Defaults to the `location` of the Network Watcher.

* `traffic_analytics` - (Optional) A `traffic_analytics` block as documented below.

* `version` - (Optional) The version (revision) of the flow log. Possible values are `1` and `2`. Defaults to `1`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Network Watcher Flow Log.

---

The `retention_policy` block supports the following:

* `enabled` - (Required) Boolean flag to enable/disable retention.
* `days` - (Required) The number of days to retain flow log records.
 
---

The `traffic_analytics` block supports the following:

* `enabled` - (Required) Boolean flag to enable/disable traffic analytics.
* `workspace_id` - (Required) The resource GUID of the attached workspace.
* `workspace_region` - (Required) The location of the attached workspace.
* `workspace_resource_id` - (Required) The resource ID of the attached workspace.
* `interval_in_minutes` - (Optional) How frequently service should do flow analytics in minutes. Defaults to `60`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Watcher.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Watcher Flow Log.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Watcher Flow Log.
* `update` - (Defaults to 30 minutes) Used when updating the Network Watcher Flow Log.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Watcher Flow Log.

## Import

Network Watcher Flow Logs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_watcher_flow_log.watcher1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/networkWatchers/watcher1/flowLogs/log1
```
