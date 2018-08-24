locals {
  dc2virtual_machine_name = "${var.prefix}-dc2"
  dc2virtual_machine_fqdn = "${local.dc2virtual_machine_name}.${var.active_directory_domain}"
  dc2custom_data_params   = "Param($RemoteHostName = \"${local.dc2virtual_machine_fqdn}\", $ComputerName = \"${local.dc2virtual_machine_name}\")"
  dc2custom_data_content  = "${local.dc2custom_data_params} ${file("${path.module}/files/winrm.ps1")}"
}

resource "azurerm_virtual_machine" "domain-controller2" {
  name                          = "${local.dc2virtual_machine_name}"
  location                      = "${var.location}"
  availability_set_id           = "${var.dcavailability_set_id}"
  resource_group_name           = "${var.resource_group_name}"
  network_interface_ids         = ["${azurerm_network_interface.dc2primary.id}"]
  vm_size                       = "Standard_A1"
  delete_os_disk_on_termination = false

  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2012-R2-Datacenter"
    version   = "latest"
  }

  storage_os_disk {
    name              = "${local.dc2virtual_machine_name}-disk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "${local.dc2virtual_machine_name}"
    admin_username = "${var.admin_username}"
    admin_password = "${var.admin_password}"
    custom_data    = "${local.dc2custom_data_content}"
  }

  os_profile_windows_config {
    provision_vm_agent        = true
    enable_automatic_upgrades = false

    additional_unattend_config {
      pass         = "oobeSystem"
      component    = "Microsoft-Windows-Shell-Setup"
      setting_name = "AutoLogon"
      content      = "<AutoLogon><Password><Value>${var.admin_password}</Value></Password><Enabled>true</Enabled><LogonCount>1</LogonCount><Username>${var.admin_username}</Username></AutoLogon>"
    }

    # Unattend config is to enable basic auth in WinRM, required for the provisioner stage.
    additional_unattend_config {
      pass         = "oobeSystem"
      component    = "Microsoft-Windows-Shell-Setup"
      setting_name = "FirstLogonCommands"
      content      = "${file("${path.module}/files/FirstLogonCommands.xml")}"
    }
  }

  depends_on = ["azurerm_network_interface.dc2primary"]
}
