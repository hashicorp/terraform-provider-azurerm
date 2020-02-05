resource "azurerm_public_ip" "static" {
  name                = "${var.prefix}-client-pip"
  location            = var.location
  resource_group_name = var.resource_group_name
  allocation_method   = "Static"
}

resource "azurerm_network_interface" "primary" {
  name                = "${var.prefix}-client-nic"
  location            = var.location
  resource_group_name = var.resource_group_name
  ip_configuration {
    name                          = "primary"
    subnet_id                     = var.subnet_id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.static.id
  }
}

resource "azurerm_windows_virtual_machine" "domain-member" {
  name                     = local.virtual_machine_name
  resource_group_name      = var.resource_group_name
  location                 = var.location
  size                     = "Standard_F2"
  admin_username           = var.admin_username
  admin_password           = var.admin_password
  provision_vm_agent       = true
  enable_automatic_updates = true

  network_interface_ids = [
    azurerm_network_interface.primary.id,
  ]

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  provisioner "remote-exec" {
    connection {
      type     = "winrm"
      user     = var.admin_username
      password = var.admin_password
      host     = azurerm_public_ip.static.ip_address
    }

  }
}

// NOTE: this is a hack.
// the AD Domain takes ~7m to provision, so we don't try and join an non-existant domain we sleep
// unfortunately we can't depend on the Domain Creation VM Extension since there's a reboot.
// We sleep for 12 minutes here to give Azure some breathing room.
//resource "null_resource" "wait-for-domain-to-provision" {
//  provisioner "local-exec" {
//    command = "sleep 720"
//  }
//
//  depends_on = [azurerm_windows_virtual_machine.domain-member]
//}

// Waits for up to 20 minutes for the Domain to become available. Will return an error 1 if unsuccessful preventing the member attempting to join.

resource "azurerm_virtual_machine_extension" "wait-for-domain-to-provision" {
  name                 = "TestConnectionDomain"
  publisher            = "Microsoft.Compute"
  type                 = "CustomScriptExtension"
  type_handler_version = "1.9"
  virtual_machine_id   = azurerm_windows_virtual_machine.domain-member.id
  settings             = <<SETTINGS
  {
    "commandToExecute": "powershell.exe -Command \"Test-Connection -TargetName ${var.active_directory_domain_name} -Count 120 -Delay 10 -Quiet \""
  }
SETTINGS
}

resource "azurerm_virtual_machine_extension" "join-domain" {
  name                 = azurerm_windows_virtual_machine.domain-member.name
  publisher            = "Microsoft.Compute"
  type                 = "JsonADDomainExtension"
  type_handler_version = "1.3"
  virtual_machine_id   = azurerm_windows_virtual_machine.domain-member.id

  settings = <<SETTINGS
    {
        "Name": "${var.active_directory_domain_name}",
        "OUPath": "",
        "User": "${var.active_directory_username}@${var.active_directory_domain_name}",
        "Restart": "true",
        "Options": "3"
    }
SETTINGS

  protected_settings = <<SETTINGS
    {
        "Password": "${var.active_directory_password}"
    }
SETTINGS

//  depends_on = [null_resource.wait-for-domain-to-provision]
  depends_on = [azurerm_virtual_machine_extension.wait-for-domain-to-provision]
}
