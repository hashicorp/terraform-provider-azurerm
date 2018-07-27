resource "azurerm_managed_disk" "external" {
  count                = "${var.number_of_disks}"
  name                 = "${var.prefix}-disk${count.index+1}"
  location             = "${azurerm_resource_group.main.location}"
  resource_group_name  = "${azurerm_resource_group.main.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "10"
}

resource "azurerm_virtual_machine_data_disk_attachment" "external" {
  count              = "${var.number_of_disks}"
  managed_disk_id    = "${azurerm_managed_disk.external.*.id[count.index]}"
  virtual_machine_id = "${azurerm_virtual_machine.main.id}"
  lun                = "${10+count.index}"
  caching            = "ReadWrite"
}
