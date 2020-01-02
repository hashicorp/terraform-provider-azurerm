---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_watcher_flow_log"
sidebar_current: "docs-azurerm-resource-network-watcher-flow-log"
description: |-
  Manages a Network Watcher Flow Log.

---

# azurerm_network_watcher_flow_log

Manages a Network Watcher Flow Log.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "eastus"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_network_watcher" "test" {
    name                = "acctestnw"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsa"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  account_tier              = "Standard"
  account_replication_type  = "LRS"
  enable_https_traffic_only = true
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestlaw"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "PerGB2018"
}

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = "${azurerm_network_watcher.test.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"

  network_security_group_id = "${azurerm_network_security_group.test.id}"
  storage_account_id        = "${azurerm_storage_account.test.id}"
  enabled                   = true
  
  retention_policy {
    enabled = true
    days    = 7
  }

  traffic_analytics {
    enabled               = true
    workspace_id          = "${azurerm_log_analytics_workspace.test.workspace_id}"
    workspace_region      = "${azurerm_log_analytics_workspace.test.location}"
    workspace_resource_id = "${azurerm_log_analytics_workspace.test.id}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `network_watcher_name` - (Required) The name of the Network Watcher. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Network Watcher was deployed. Changing this forces a new resource to be created.

* `network_security_group_id` - (Required) The ID of the Network Security Group for which to enable flow logs for. Changing this forces a new resource to be created.
  
* `storage_account_id` - (Required) The ID of the Storage Account where flow logs are stored.
  
* `enabled` - (Required) Should Network Flow Logging be Enabled?

* `retention_policy` - (Required) A `retention_policy` block as documented below.

* `traffic_analytics` - (Optional) A `traffic_analytics` block as documented below.

---

* `retention_policy` supports the following:

* `enabled` - (Required) Boolean flag to enable/disable retention.
* `days` - (Required) The number of days to retain flow log records.

---

* `traffic_analytics` supports the following:

* `enabled` - (Required) Boolean flag to enable/disable traffic analytics.
* `workspace_id` - (Required) The resource guid of the attached workspace.
* `workspace_region` - (Required) The location of the attached workspace.
* `workspace_resource_id` - (Required) The resource ID of the attached workspace.

## Attributes Reference

The following attributes are exported:

* `id` - The Network Watcher ID.

## Import

Network Watcher Flow Logs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_watcher_flow_log.watcher1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/networkWatchers/watcher1
```
