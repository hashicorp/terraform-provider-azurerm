---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_packet_capture"
description: |-
  Configures Packet Capturing against a Virtual Machine using a Network Watcher.

---

# azurerm_packet_capture

Configures Packet Capturing against a Virtual Machine using a Network Watcher.

~> **NOTE:** This resource has been deprecated in favour of the `azurerm_network_connection_monitor` resource and will be removed in the next major version of the AzureRM Provider. The new resource shares the same fields as this one, and information on migrating across [can be found in this guide](../guides/migrating-between-renamed-resources.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "packet-capture-rg"
  location = "West Europe"
}

resource "azurerm_network_watcher" "example" {
  name                = "network-watcher"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "production-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name               = "internal"
  virtual_network_id = azurerm_virtual_network.example.id
  address_prefixes   = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "example" {
  name                = "pctest-nic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "example" {
  name                  = "pctest-vm"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  network_interface_ids = [azurerm_network_interface.example.id]
  vm_size               = "Standard_F2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osdisk"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "pctest-vm"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_virtual_machine_extension" "example" {
  name                       = "network-watcher"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  virtual_machine_name       = azurerm_virtual_machine.example.name
  publisher                  = "Microsoft.Azure.NetworkWatcher"
  type                       = "NetworkWatcherAgentLinux"
  type_handler_version       = "1.4"
  auto_upgrade_minor_version = true
}

resource "azurerm_storage_account" "example" {
  name                     = "pctestsa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_packet_capture" "example" {
  name                 = "pctestcapture"
  network_watcher_name = azurerm_network_watcher.example.name
  resource_group_name  = azurerm_resource_group.example.name
  target_resource_id   = azurerm_virtual_machine.example.id

  storage_location {
    storage_account_id = azurerm_storage_account.example.id
  }

  depends_on = [azurerm_virtual_machine_extension.example]
}
```

~> **NOTE:** This Resource requires that [the Network Watcher Virtual Machine Extension](https://docs.microsoft.com/azure/network-watcher/network-watcher-packet-capture-manage-portal#before-you-begin) is installed on the Virtual Machine before capturing can be enabled which can be installed via [the `azurerm_virtual_machine_extension` resource](virtual_machine_extension.html).

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to use for this Packet Capture. Changing this forces a new resource to be created.

* `network_watcher_name` - (Required) The name of the Network Watcher. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Network Watcher exists. Changing this forces a new resource to be created.

* `target_resource_id` - (Required) The ID of the Resource to capture packets from. Changing this forces a new resource to be created.

~> **NOTE:** Currently only Virtual Machines ID's are supported.

* `maximum_bytes_per_packet` - (Optional) The number of bytes captured per packet. The remaining bytes are truncated. Defaults to `0` (Entire Packet Captured). Changing this forces a new resource to be created.

* `maximum_bytes_per_session` - (Optional) Maximum size of the capture in Bytes. Defaults to `1073741824` (1GB). Changing this forces a new resource to be created.

* `maximum_capture_duration` - (Optional) The maximum duration of the capture session in seconds. Defaults to `18000` (5 hours). Changing this forces a new resource to be created.

* `storage_location` - (Required) A `storage_location` block as defined below. Changing this forces a new resource to be created.

* `filter` - (Optional) One or more `filter` blocks as defined below. Changing this forces a new resource to be created.

---

A `storage_location` block contains:

* `file_path` - (Optional) A valid local path on the targeting VM. Must include the name of the capture file (*.cap). For linux virtual machine it must start with `/var/captures`.

* `storage_account_id` - (Optional) The ID of the storage account to save the packet capture session

~> **NOTE:** At least one of `file_path` or `storage_account_id` must be specified.

A `filter` block contains:

* `local_ip_address` - (Optional) The local IP Address to be filtered on. Notation: "127.0.0.1" for single address entry. "127.0.0.1-127.0.0.255" for range. "127.0.0.1;127.0.0.5" for multiple entries. Multiple ranges not currently supported. Mixing ranges with multiple entries not currently supported. Changing this forces a new resource to be created.

* `local_port` - (Optional) The local port to be filtered on. Notation: "80" for single port entry."80-85" for range. "80;443;" for multiple entries. Multiple ranges not currently supported. Mixing ranges with multiple entries not currently supported. Changing this forces a new resource to be created.

* `protocol` - (Required) The Protocol to be filtered on. Possible values include `Any`, `TCP` and `UDP`. Changing this forces a new resource to be created.

* `remote_ip_address` - (Optional) The remote IP Address to be filtered on. Notation: "127.0.0.1" for single address entry. "127.0.0.1-127.0.0.255" for range. "127.0.0.1;127.0.0.5;" for multiple entries. Multiple ranges not currently supported. Mixing ranges with multiple entries not currently supported.. Changing this forces a new resource to be created.

* `remote_port` - (Optional) The remote port to be filtered on. Notation: "80" for single port entry."80-85" for range. "80;443;" for multiple entries. Multiple ranges not currently supported. Mixing ranges with multiple entries not currently supported. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The Packet Capture ID.

* `storage_location` - (Required) A `storage_location` block as defined below.

---

A `storage_location` block contains:

* `storage_path` - The URI of the storage path to save the packet capture.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Packet Capture.
* `update` - (Defaults to 30 minutes) Used when updating the Packet Capture.
* `read` - (Defaults to 5 minutes) Used when retrieving the Packet Capture.
* `delete` - (Defaults to 30 minutes) Used when deleting the Packet Capture.

## Import

Packet Captures can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_packet_capture.capture1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/networkWatchers/watcher1/packetCaptures/capture1
```
