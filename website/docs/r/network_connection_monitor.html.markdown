---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_connection_monitor"
description: |-
  Manages a Network Connection Monitor.
---

# azurerm_network_connection_monitor

Manages a Network Connection Monitor.

~> **NOTE:** As `test_frequency_sec` has default value, so terraform cannot make `test_frequency_sec` compatible with `interval_in_seconds` while flattens `test_frequency_sec`.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-watcher-resources"
  location = "eastus2"
}

resource "azurerm_network_watcher" "example" {
  name                = "example-watcher"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "example" {
  name                = "example-nic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "example" {
  name                  = "example-vm"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  network_interface_ids = [azurerm_network_interface.example.id]
  vm_size               = "Standard_D2s_v3"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osdisk-example01"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "hostnametest01"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_virtual_machine_extension" "example" {
  name                       = "example-vmextension"
  virtual_machine_id         = azurerm_virtual_machine.example.id
  publisher                  = "Microsoft.Azure.NetworkWatcher"
  type                       = "NetworkWatcherAgentLinux"
  type_handler_version       = "1.4"
  auto_upgrade_minor_version = true
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "pergb2018"
}

resource "azurerm_network_connection_monitor" "example" {
  name                 = "example-monitor"
  network_watcher_name = azurerm_network_watcher.example.name
  resource_group_name  = azurerm_resource_group.example.name
  location             = azurerm_network_watcher.example.location

  endpoint {
    name               = "source"
    virtual_machine_id = azurerm_virtual_machine.example.id

    filter {
      item {
        address = azurerm_virtual_machine.example.id
        type    = "AgentAddress"
      }

      type = "Include"
    }
  }

  endpoint {
    name    = "destination"
    address = "terraform.io"
  }

  test_configuration {
    name               = "tcpName"
    protocol           = "Tcp"
    test_frequency_sec = 60

    tcp_configuration {
      port = 80
    }
  }

  test_group {
    name                = "exampletg"
    destinations        = ["destination"]
    sources             = ["source"]
    test_configurations = ["tcpName"]
    disable             = false
  }

  notes = "examplenote"

  output_workspace_resource_ids = [azurerm_log_analytics_workspace.example.id]

  depends_on = [azurerm_virtual_machine_extension.example]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Network Connection Monitor. Changing this forces a new Network Connection Monitor to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Network Connection Monitor should exist. Changing this forces a new Network Connection Monitor to be created.

* `location` - (Required) The Azure Region where the Network Connection Monitor should exist. Changing this forces a new Network Connection Monitor to be created.

* `network_watcher_name` - (Required) The name of the Network Watcher. Changing this forces a new Network Connection Monitor to be created.

---

* `auto_start` - (Optional / **Deprecated**) Will the connection monitor start automatically once created?

~> **NOTE:** The field `auto_start` has been deprecated in new api version 2020-05-01.

* `destination` - (Optional / **Deprecated**) A `destination` block as defined below.

~> **NOTE:** The field `destination` has been deprecated in favor of `endpoint`.

* `endpoint` - (Optional) A `endpoint` block as defined below.

* `interval_in_seconds` - (Optional / **Deprecated**) Monitoring interval in seconds.

~> **NOTE:** The field `interval_in_seconds` has been deprecated in favor of `test_frequency_sec`.

* `notes` - (Optional) The notes to be associated with the connection monitor.

* `output_workspace_resource_ids` - (Optional) A list of the log analytics workspace id.

* `source` - (Optional / **Deprecated**) A `source` block as defined below.

~> **NOTE:** The field `source` has been deprecated in favor of `endpoint`.

* `test_configuration` - (Optional) A `test_configuration` block as defined below.

* `test_group` - (Optional) A `test_group` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Network Connection Monitor.

---

A `destination` block supports the following:

* `port` - (Required) The destination port used by connection monitor.

* `address` - (Optional) The address of the connection monitor destination (IP or domain name). Conflicts with `destination.0.virtual_machine_id`

* `virtual_machine_id` - (Optional) The ID of the virtual machine used as the destination by connection monitor. Conflicts with `destination.0.address`

---

A `endpoint` block supports the following:

* `name` - (Required) The name of the connection monitor endpoint.

* `address` - (Optional) The address of the connection monitor endpoint (IP or domain name).

* `filter` - (Optional) A `filter` block as defined below.

* `virtual_machine_id` - (Optional) The ID of the virtual machine used as the endpoint by connection monitor.

---

A `filter` block supports the following:

* `type` - (Optional) The behavior of the endpoint filter. Possible value is `Include`. Defaults to `Include`.

* `item` - (Optional) A `item` block as defined below.

---

A `item` block supports the following:

* `type` - (Optional) The type of item included in the filter. Possible value is `AgentAddress`. Defaults to `AgentAddress`.

* `address` - (Optional) The address of the filter item.

---

A `source` block supports the following:

* `virtual_machine_id` - (Required) The ID of the virtual machine used as the source by connection monitor.

* `port` - (Optional) The source port used by connection monitor. Defaults to `0`.

---

A `test_configuration` block supports the following:

* `name` - (Required) The name of the connection monitor test configuration.

* `protocol` - (Required) The protocol to use in test evaluation. Possible values are `Tcp`, `Http` and `Icmp`.

* `test_frequency_sec` - (Optional) The frequency of test evaluation, in seconds. Defaults to `60`.

* `http_configuration` - (Optional) A `http_configuration` block as defined below.

* `icmp_configuration` - (Optional) A `icmp_configuration` block as defined below.

* `preferred_ip_version` - (Optional) The preferred IP version to use in test evaluation. Possible values are `IPv4` and `IPv6`. 

* `success_threshold` - (Optional) A `success_threshold` block as defined below.

* `tcp_configuration` - (Optional) A `tcp_configuration` block as defined below.

---

A `http_configuration` block supports the following:

* `method` - (Optional) The HTTP method to use. Possible values are `Get` and `Post`. Defaults to `Get`.

* `port` - (Optional) The port to connect to.

* `path` - (Optional) The path component of the URI. For instance, "/dir1/dir2".

* `prefer_https` - (Optional) Will https be preferred over http in cases where the choice is not explicit? Defaults to `false`.

* `request_header` - (Optional) A `request_header` block as defined below.

* `valid_status_code_ranges` - (Optional) The http status codes to consider successful. For instance, "2xx, 301-304, 418".

---

A `request_header` block supports the following:

* `name` - (Required) The name in HTTP header.

* `value` - (Required) The value in HTTP header.

---

A `icmp_configuration` block supports the following:

* `disable_trace_route` - (Optional) Will path evaluation with trace route be disabled? Defaults to `false`.

---

A `success_threshold` block supports the following:

* `checks_failed_percent` - (Optional) The maximum percentage of failed checks permitted for a test to evaluate as successful.

* `round_trip_time_ms` - (Optional) The maximum round-trip time in milliseconds permitted for a test to evaluate as successful.

---

A `tcp_configuration` block supports the following:

* `port` - (Required) The port to connect to.

* `disable_trace_route` - (Optional) Should path evaluation with trace route be disabled? Defaults to `false`.

---

A `test_group` block supports the following:

* `name` - (Required) The name of the connection monitor test group.

* `destinations` - (Required) A list of destination endpoint names.

* `sources` - (Required) A list of source endpoint names.

* `test_configurations` - (Required) A list of test configuration names.

* `enabled` - (Optional) Should the test group be enabled? Defaults to `true`.

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
