# Install the Network Watcher agent on the Virtual Machine
resource "azurerm_virtual_machine_extension" "example" {
  name                       = "network-watcher"
  location                   = "${azurerm_resource_group.example.location}"
  resource_group_name        = "${azurerm_resource_group.example.name}"
  virtual_machine_name       = "${azurerm_virtual_machine.example.name}"
  publisher                  = "Microsoft.Azure.NetworkWatcher"
  type                       = "NetworkWatcherAgentLinux"
  type_handler_version       = "1.4"
  auto_upgrade_minor_version = true
}

# Which allows for Packets to be captured
resource "azurerm_packet_capture" "example" {
  name                 = "${var.prefix}capture"
  network_watcher_name = "${azurerm_network_watcher.example.name}"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  target_resource_id   = "${azurerm_virtual_machine.example.id}"

  storage_location {
    storage_account_id = "${azurerm_storage_account.example.id}"
  }

  depends_on = ["azurerm_virtual_machine_extension.example"]
}
