---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_watcher_flow_log"
description: |-
    Lists Network Watcher Flow Log resources.
---

# List resource: azurerm_network_watcher_flow_log

Lists Network Watcher Flow Log resources.

## Example Usage

### List all Network Watcher Flow Logs in a Network Watcher

```hcl
list "azurerm_network_watcher_flow_log" "example" {
  provider = azurerm
  config {
    network_watcher_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/networkWatchers/watcher1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `network_watcher_id` - (Required) The ID of the Network Watcher to query.
