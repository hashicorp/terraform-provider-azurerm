---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_connection_monitor"
description: |-
  Manages a Network Connection Monitor.
---

# azurerm_network_connection_monitor

Manages a Network Connection Monitor.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_resource_group" "example" {
  name = "example-resources"
}

resource "azurerm_network_watcher" "example" {
  name                = "example-nw"
  location            = data.azurerm_resource_group.example.location
  resource_group_name = data.azurerm_resource_group.example.name
}

data "azurerm_virtual_machine" "src" {
  name                = "example-vm"
  resource_group_name = data.azurerm_resource_group.example.name
}

resource "azurerm_virtual_machine_extension" "src" {
  name                       = "network-watcher"
  virtual_machine_id         = data.azurerm_virtual_machine.src.id
  publisher                  = "Microsoft.Azure.NetworkWatcher"
  type                       = "NetworkWatcherAgentLinux"
  type_handler_version       = "1.4"
  auto_upgrade_minor_version = true
}

resource "azurerm_network_connection_monitor" "example" {
  name                 = "example-ncm"
  network_watcher_name = azurerm_network_watcher.example.name
  resource_group_name  = data.azurerm_resource_group.example.name
  location             = azurerm_network_watcher.example.location

  auto_start          = false
  interval_in_seconds = 30

  source {
    virtual_machine_id = data.azurerm_virtual_machine.src.id
    port               = 20020
  }

  destination {
    address = "terraform.io"
    port    = 443
  }

  tags = {
    foo = "bar"
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Network Connection Monitor. Changing this forces a new Network Connection Monitor to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Network Connection Monitor should exist. Changing this forces a new Network Connection Monitor to be created.

* `location` - (Required) The Azure Region where the Network Connection Monitor should exist. Changing this forces a new Network Connection Monitor to be created.

* `destination` - (Required) A `destination` block as defined below.

* `network_watcher_name` - (Required) The name of the Network Watcher. Changing this forces a new Network Connection Monitor to be created.

* `source` - (Required) A `source` block as defined below.

---

* `auto_start` - (Optional) Will the connection monitor start automatically once created? Changing this forces a new Network Connection Monitor to be created.

* `interval_in_seconds` - (Optional) Monitoring interval in seconds.

* `tags` - (Optional) A mapping of tags which should be assigned to the Network Connection Monitor.

---

A `destination` block supports the following:

* `port` - (Required) The destination port used by connection monitor.

* `address` - (Optional) The address of the connection monitor destination (IP or domain name). Conflicts with `destination.0.virtual_machine_id`

* `virtual_machine_id` - (Optional) The ID of the virtual machine used as the destination by connection monitor. Conflicts with `destination.0.address`

---

A `source` block supports the following:

* `virtual_machine_id` - (Required) The ID of the virtual machine used as the source by connection monitor.

* `port` - (Optional) The source port used by connection monitor.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Network Connection Monitor.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Connection Monitor.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Connection Monitor.
* `update` - (Defaults to 30 minutes) Used when updating the Network Connection Monitor.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Connection Monitor.

## Import

Network Connection Monitors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_connection_monitor.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/networkWatchers/watcher1/connectionMonitors/connectionMonitor1
```
