resource "azurerm_availability_set" "sqlavailabilityset" {
  name                         = "sqlavailabilityset"
  resource_group_name          = "${var.resource_group_name}"
  location                     = "${var.location}"
  platform_fault_domain_count  = 3
  platform_update_domain_count = 5
  managed                      = true
}

resource "azurerm_virtual_machine" "sql" {
  name                          = "${var.prefix}-sql${1 + count.index}"
  location                      = "${var.location}"
  availability_set_id           = "${azurerm_availability_set.sqlavailabilityset.id}"
  resource_group_name           = "${var.resource_group_name}"
  network_interface_ids         = ["${element(azurerm_network_interface.primary.*.id, count.index)}"]
  vm_size                       = "Standard_B1s"
  delete_os_disk_on_termination = true
  count                         = "${var.sqlvmcount}"

  storage_image_reference {
    publisher = "MicrosoftSQLServer"
    offer     = "SQL2014SP2-WS2012R2"
    sku       = "Enterprise"
    version   = "latest"
  }

  storage_os_disk {
    name              = "${var.prefix}-sql${1 + count.index}-disk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "${var.prefix}-sql${1 + count.index}"
    admin_username = "${var.admin_username}"
    admin_password = "${var.admin_password}"
  }

  os_profile_windows_config {
    provision_vm_agent        = true
    enable_automatic_upgrades = false
  }

   storage_data_disk {
    name              = "${var.prefix}-sql${1 + count.index}-data-disk1"
    disk_size_gb      = "2000"
    caching           = "ReadWrite"
    create_option     = "Empty"
    managed_disk_type = "Standard_LRS"
    lun               = "2"
  }

   storage_data_disk {
    name              = "${var.prefix}-sql${1 + count.index}-log-disk1"
    disk_size_gb      = "500"
    caching           = "ReadWrite"
    create_option     = "Empty"
    managed_disk_type = "Standard_LRS"
    lun               = "3"
  }

  depends_on = ["azurerm_network_interface.primary"]
}
