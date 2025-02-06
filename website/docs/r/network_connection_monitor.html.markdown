---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_connection_monitor"
description: |-
  Manages a Network Connection Monitor.
---

# azurerm_network_connection_monitor

Manages a Network Connection Monitor.

~> **NOTE:** Any Network Connection Monitor resource created with API versions 2019-06-01 or earlier (v1) are now incompatible with Terraform, which now only supports v2.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-Watcher-resources"
  location = "West Europe"
}

resource "azurerm_network_watcher" "example" {
  name                = "example-Watcher"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "example-Vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-Subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "example" {
  name                = "example-Nic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "example" {
  name                  = "example-VM"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  network_interface_ids = [azurerm_network_interface.example.id]
  vm_size               = "Standard_D2s_v3"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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
  name                       = "example-VMExtension"
  virtual_machine_id         = azurerm_virtual_machine.example.id
  publisher                  = "Microsoft.Azure.NetworkWatcher"
  type                       = "NetworkWatcherAgentLinux"
  type_handler_version       = "1.4"
  auto_upgrade_minor_version = true
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-Workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_network_connection_monitor" "example" {
  name               = "example-Monitor"
  network_watcher_id = azurerm_network_watcher.example.id
  location           = azurerm_network_watcher.example.location

  endpoint {
    name               = "source"
    target_resource_id = azurerm_virtual_machine.example.id

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
    name                      = "tcpName"
    protocol                  = "Tcp"
    test_frequency_in_seconds = 60

    tcp_configuration {
      port = 80
    }
  }

  test_group {
    name                     = "exampletg"
    destination_endpoints    = ["destination"]
    source_endpoints         = ["source"]
    test_configuration_names = ["tcpName"]
  }

  notes = "examplenote"

  output_workspace_resource_ids = [azurerm_log_analytics_workspace.example.id]

  depends_on = [azurerm_virtual_machine_extension.example]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Network Connection Monitor. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Network Connection Monitor should exist. Changing this forces a new resource to be created.

* `network_watcher_id` - (Required) The ID of the Network Watcher. Changing this forces a new resource to be created.

* `endpoint` - (Required) A `endpoint` block as defined below.

* `test_configuration` - (Required) A `test_configuration` block as defined below.

* `test_group` - (Required) A `test_group` block as defined below.

---

* `notes` - (Optional) The description of the Network Connection Monitor.

* `output_workspace_resource_ids` - (Optional) A list of IDs of the Log Analytics Workspace which will accept the output from the Network Connection Monitor.

* `tags` - (Optional) A mapping of tags which should be assigned to the Network Connection Monitor.

---

A `endpoint` block supports the following:

* `name` - (Required) The name of the endpoint for the Network Connection Monitor .

* `address` - (Optional) The IP address or domain name of the Network Connection Monitor endpoint.

* `coverage_level` - (Optional) The test coverage for the Network Connection Monitor endpoint. Possible values are `AboveAverage`, `Average`, `BelowAverage`, `Default`, `Full` and `Low`.

* `excluded_ip_addresses` - (Optional) A list of IPv4/IPv6 subnet masks or IPv4/IPv6 IP addresses to be excluded to the Network Connection Monitor endpoint.

* `included_ip_addresses` - (Optional) A list of IPv4/IPv6 subnet masks or IPv4/IPv6 IP addresses to be included to the Network Connection Monitor endpoint.

* `target_resource_id` - (Optional) The resource ID which is used as the endpoint by the Network Connection Monitor.

* `filter` - (Optional) A `filter` block as defined below.

* `target_resource_type` - (Optional) The endpoint type of the Network Connection Monitor. Possible values are `AzureArcVM`, `AzureSubnet`, `AzureVM`, `AzureVNet`, `ExternalAddress`, `MMAWorkspaceMachine` and `MMAWorkspaceNetwork`.

---

A `filter` block supports the following:

* `type` - (Optional) The behaviour type of this endpoint filter. Currently the only allowed value is `Include`. Defaults to `Include`.

* `item` - (Optional) A `item` block as defined below.

---

A `item` block supports the following:

* `type` - (Optional) The type of items included in the filter. Possible values are `AgentAddress`. Defaults to `AgentAddress`.

* `address` - (Optional) The address of the filter item.

---

A `test_configuration` block supports the following:

* `name` - (Required) The name of test configuration for the Network Connection Monitor.

* `protocol` - (Required) The protocol used to evaluate tests. Possible values are `Tcp`, `Http` and `Icmp`.

* `test_frequency_in_seconds` - (Optional) The time interval in seconds at which the test evaluation will happen. Defaults to `60`.

* `http_configuration` - (Optional) A `http_configuration` block as defined below.

* `icmp_configuration` - (Optional) A `icmp_configuration` block as defined below.

* `preferred_ip_version` - (Optional) The preferred IP version which is used in the test evaluation. Possible values are `IPv4` and `IPv6`.

* `success_threshold` - (Optional) A `success_threshold` block as defined below.

* `tcp_configuration` - (Optional) A `tcp_configuration` block as defined below.

---

A `http_configuration` block supports the following:

* `method` - (Optional) The HTTP method for the HTTP request. Possible values are `Get` and `Post`. Defaults to `Get`.

* `port` - (Optional) The port for the HTTP connection.

* `path` - (Optional) The path component of the URI. It only accepts the absolute path.

* `prefer_https` - (Optional) Should HTTPS be preferred over HTTP in cases where the choice is not explicit? Defaults to `false`.

* `request_header` - (Optional) A `request_header` block as defined below.

* `valid_status_code_ranges` - (Optional) The HTTP status codes to consider successful. For instance, `2xx`, `301-304` and `418`.

---

A `request_header` block supports the following:

* `name` - (Required) The name of the HTTP header.

* `value` - (Required) The value of the HTTP header.

---

A `icmp_configuration` block supports the following:

* `trace_route_enabled` - (Optional) Should path evaluation with trace route be enabled? Defaults to `true`.

---

A `success_threshold` block supports the following:

* `checks_failed_percent` - (Optional) The maximum percentage of failed checks permitted for a test to be successful.

* `round_trip_time_ms` - (Optional) The maximum round-trip time in milliseconds permitted for a test to be successful.

---

A `tcp_configuration` block supports the following:

* `port` - (Required) The port for the TCP connection.

* `trace_route_enabled` - (Optional) Should path evaluation with trace route be enabled? Defaults to `true`.

* `destination_port_behavior` - (Optional) The destination port behavior for the TCP connection. Possible values are `None` and `ListenIfAvailable`.

---

A `test_group` block supports the following:

* `name` - (Required) The name of the test group for the Network Connection Monitor.

* `destination_endpoints` - (Required) A list of destination endpoint names.

* `source_endpoints` - (Required) A list of source endpoint names.

* `test_configuration_names` - (Required) A list of test configuration names.

* `enabled` - (Optional) Should the test group be enabled? Defaults to `true`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Connection Monitor.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Connection Monitor.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Connection Monitor.
* `update` - (Defaults to 30 minutes) Used when updating the Network Connection Monitor.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Connection Monitor.

## Import

Network Connection Monitors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_connection_monitor.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/networkWatchers/watcher1/connectionMonitors/connectionMonitor1
```
