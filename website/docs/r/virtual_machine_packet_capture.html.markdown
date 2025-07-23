---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_packet_capture"
description: |-
  Configures Packet Capturing against a Virtual Machine using a Network Watcher.

---

# azurerm_virtual_machine_packet_capture

Configures Network Packet Capturing against a Virtual Machine using a Network Watcher.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_network_watcher" "example" {
  name                = "example-nw"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
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
  vm_size               = "Standard_F2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
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
  virtual_machine_id         = azurerm_virtual_machine.example.id
  publisher                  = "Microsoft.Azure.NetworkWatcher"
  type                       = "NetworkWatcherAgentLinux"
  type_handler_version       = "1.4"
  auto_upgrade_minor_version = true
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_virtual_machine_packet_capture" "example" {
  name               = "example-pc"
  network_watcher_id = azurerm_network_watcher.example.id
  virtual_machine_id = azurerm_virtual_machine.example.id

  storage_location {
    storage_account_id = azurerm_storage_account.example.id
  }

  depends_on = [azurerm_virtual_machine_extension.example]
}
```

~> **Note:** This Resource requires that [the Network Watcher Virtual Machine Extension](https://docs.microsoft.com/azure/network-watcher/network-watcher-packet-capture-manage-portal#before-you-begin) is installed on the Virtual Machine before capturing can be enabled which can be installed via [the `azurerm_virtual_machine_extension` resource](virtual_machine_extension.html).

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to use for this Network Packet Capture. Changing this forces a new resource to be created.

* `network_watcher_id` - (Required) The resource ID of the Network Watcher. Changing this forces a new resource to be created.

* `virtual_machine_id` - (Required) The resource ID of the target Virtual Machine to capture packets from. Changing this forces a new resource to be created.

* `maximum_bytes_per_packet` - (Optional) The number of bytes captured per packet. The remaining bytes are truncated. Defaults to `0` (Entire Packet Captured). Changing this forces a new resource to be created.

* `maximum_bytes_per_session` - (Optional) Maximum size of the capture in Bytes. Defaults to `1073741824` (1GB). Changing this forces a new resource to be created.

* `maximum_capture_duration_in_seconds` - (Optional) The maximum duration of the capture session in seconds. Defaults to `18000` (5 hours). Changing this forces a new resource to be created.

* `storage_location` - (Required) A `storage_location` block as defined below. Changing this forces a new resource to be created.

* `filter` - (Optional) One or more `filter` blocks as defined below. Changing this forces a new resource to be created.

---

A `storage_location` block contains:

* `file_path` - (Optional) A valid local path on the target Virtual Machine. Must include the name of the capture file (*.cap). For Linux Virtual Machines it must start with `/var/captures`.

* `storage_account_id` - (Optional) The ID of the storage account where the packet capture sessions should be saved to.

~> **Note:** At least one of `file_path` or `storage_account_id` must be specified.

---

A `filter` block contains:

* `local_ip_address` - (Optional) The local IP Address to be filtered on. Specify `127.0.0.1` for a single address entry, `127.0.0.1-127.0.0.255` for a range and `127.0.0.1;127.0.0.5` for multiple entries. Multiple ranges and mixing ranges with multiple entries are currently not supported. Changing this forces a new resource to be created.

* `local_port` - (Optional) The local port to be filtered on. Specify `80` for single port entry, `80-85` for a range and `80;443;` for multiple entries. Multiple ranges and mixing ranges with multiple entries are currently not supported. Changing this forces a new resource to be created.

* `protocol` - (Required) The Protocol to be filtered on. Possible values include `Any`, `TCP` and `UDP`. Changing this forces a new resource to be created.

* `remote_ip_address` - (Optional) The remote IP Address to be filtered on. Specify `127.0.0.1` for a single address entry, `127.0.0.1-127.0.0.255` for a range and `127.0.0.1;127.0.0.5` for multiple entries. Multiple ranges and mixing ranges with multiple entries are currently not supported. Changing this forces a new resource to be created.

* `remote_port` - (Optional) The remote port to be filtered on. Specify `80` for single port entry, `80-85` for a range and `80;443;` for multiple entries. Multiple ranges and mixing ranges with multiple entries are currently not supported. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Virtual Machine Packet Capture ID.

* `storage_location` - (Required) A `storage_location` block as defined below.

---

A `storage_location` block contains:

* `storage_path` - The URI of the storage path where the packet capture sessions are saved to.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Machine Packet Capture.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Machine Packet Capture.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Machine Packet Capture.

## Import

Virtual Machine Packet Captures can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_machine_packet_capture.capture1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/networkWatchers/watcher1/packetCaptures/capture1
```
