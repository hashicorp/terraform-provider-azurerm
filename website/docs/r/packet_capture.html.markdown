---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_packet_capture"
sidebar_current: "docs-azurerm-resource-network-packet-capture"
description: |-
  Configures Packet Capturing against a Virtual Machine using a Network Watcher.

---

# azurerm_packet_capture

Configures Packet Capturing against a Virtual Machine using a Network Watcher.

## Example Usage

A complete example of how to use the `azurerm_packet_capture` resource can be found [in the `./examples/packet-capture` folder within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/packet-capture)


```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_storage_account" "example" {
  # ...
}

resource "azurerm_virtual_machine" "example" {
  # ...
}

resource "azurerm_network_watcher" "example" {
  # ...
}

resource "azurerm_packet_capture" "example" {
  name                 = "example-capture"
  network_watcher_name = "${azurerm_network_watcher.test.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  target_resource_id   = "${azurerm_virtual_machine.test.id}"

  storage_location {
    storage_account_id = "${azurerm_storage_account.test.id}"
  }
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

## Import

Packet Captures can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_packet_capture.capture1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/networkWatchers/watcher1/packetCaptures/capture1
```
