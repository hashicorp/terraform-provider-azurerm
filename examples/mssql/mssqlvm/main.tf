provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-VN"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "${var.prefix}-SN"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.0.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "example" {
  subnet_id                 = azurerm_subnet.example.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_public_ip" "vm" {
  name                = "${var.prefix}-PIP"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Dynamic"
}

resource "azurerm_network_security_group" "example" {
  name                = "${var.prefix}-NSG"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_network_security_rule" "RDPRule" {
  name                        = "RDPRule"
  resource_group_name         = azurerm_resource_group.example.name
  priority                    = 1000
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = 3389
  source_address_prefix       = "167.220.255.0/25"
  destination_address_prefix  = "*"
  network_security_group_name = azurerm_network_security_group.example.name
}

resource "azurerm_network_security_rule" "MSSQLRule" {
  name                        = "MSSQLRule"
  resource_group_name         = azurerm_resource_group.example.name
  priority                    = 1001
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = 1433
  source_address_prefix       = "167.220.255.0/25"
  destination_address_prefix  = "*"
  network_security_group_name = azurerm_network_security_group.example.name
}

resource "azurerm_network_interface" "example" {
  name                = "${var.prefix}-NIC"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "exampleconfiguration1"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.vm.id
  }
}

resource "azurerm_network_interface_security_group_association" "example" {
  network_interface_id      = azurerm_network_interface.example.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_virtual_machine" "example" {
  name                  = "${var.prefix}-VM"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  network_interface_ids = [azurerm_network_interface.example.id]
  vm_size               = "Standard_DS14_v2"

  storage_image_reference {
    publisher = "MicrosoftSQLServer"
    offer     = "SQL2017-WS2016"
    sku       = "SQLDEV"
    version   = "laexample"
  }

  storage_os_disk {
    name              = "${var.prefix}-OSDisk"
    caching           = "ReadOnly"
    create_option     = "FromImage"
    managed_disk_type = "Premium_LRS"
  }

  os_profile {
    computer_name  = "winhost01"
    admin_username = "exampleadmin"
    admin_password = "Password1234!"
  }

  os_profile_windows_config {
    timezone                  = "Pacific Standard Time"
    provision_vm_agent        = true
    enable_automatic_upgrades = true
  }
}

resource "azurerm_mssql_virtual_machine" "example" {
  virtual_machine_id = azurerm_virtual_machine.example.id
  sql_license_type   = "PAYG"
}
