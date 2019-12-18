---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_watcher"
sidebar_current: "docs-azurerm-datasource-network-watcher"
description: |-
  Gets information about an existing Network Watcher.
---

# Data Source: azurerm_network_watcher

Use this data source to access information about an existing Network Watcher.

## Example Usage

```hcl
data "azurerm_network_watcher" "example" {
  name                = "${azurerm_network_watcher.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

output "network_watcher_id" {
  value = "${data.azurerm_network_watcher.example.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the Name of the Network Watcher.
* `resource_group_name` - (Required) Specifies the Name of the Resource Group within which the Network Watcher exists.


## Attributes Reference

* `id` - The ID of the Network Watcher.

* `location` - The supported Azure location where the resource exists.

* `tags` - A mapping of tags assigned to the resource.
