# Unique name for the diagnostics storage account
resource "random_id" "diagnostics_storage_account_name" {
  keepers = {
    # Generate a new id each time a new resource group is defined
    resource_group = "${var.resource_group_name}"
  }

  byte_length = 8
}

variable "ubuntu_os_version_allowedvalues" {
  type = "map"

  default = {
    "16.04.0-LTS" = "16.04.0-LTS"
    "16.10"       = "16.10"
  }
}

variable "config" {
  type = "map"

  default = {
    "namespace"                         = "docker"
    "vm_size"                           = "Standard_DS2_v2"
    "vm_image_publisher"                = "Canonical"
    "vm_image_offer"                    = "UbuntuServer"
    "diagnostics_storage_account_type"  = "Standard_LRS"
    "managed_disk_storage_account_type" = "Premium_LRS"
    "network_address_prefix"            = "10.0.0.0/16"
    "network_subnet_prefix"             = "10.0.0.0/24"
    "network_public_ipaddress_type"     = "Static"
  }
}

resource "azurerm_resource_group" "resource_group" {
  name     = "${random_id.diagnostics_storage_account_name.keepers.resource_group}"
  location = "${var.resource_group_location}"

  tags {
    Source = "Azure Quickstarts for Terraform"
  }
}

resource "azurerm_storage_account" "diagnostics_storage_account" {
  name                = "docker${random_id.diagnostics_storage_account_name.hex}"
  resource_group_name = "${azurerm_resource_group.resource_group.name}"
  location            = "${var.resource_group_location}"

  account_type = "${var.config["diagnostics_storage_account_type"]}"

  tags {
    Source = "Azure Quickstarts for Terraform"
  }
}

resource "azurerm_public_ip" "public_ip" {
  name                = "${var.config["namespace"]}-publicip"
  resource_group_name = "${azurerm_resource_group.resource_group.name}"
  location            = "${var.resource_group_location}"

  public_ip_address_allocation = "${var.config["network_public_ipaddress_type"]}"
  domain_name_label            = "${var.dns_name_for_public_ip}"

  tags {
    Source = "Azure Quickstarts for Terraform"
  }
}

resource "azurerm_network_security_group" "network_security_group" {
  name                = "${var.config["namespace"]}-nsg"
  resource_group_name = "${azurerm_resource_group.resource_group.name}"
  location            = "${var.resource_group_location}"

  security_rule {
    name                       = "allow_ssh"
    description                = "Allow inbound traffic on default ssh port 22."
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  tags {
    Source = "Azure Quickstarts for Terraform"
  }
}

resource "azurerm_virtual_network" "virtual_network" {
  name                = "${var.config["namespace"]}-vnet"
  resource_group_name = "${azurerm_resource_group.resource_group.name}"
  location            = "${var.resource_group_location}"

  address_space = ["${var.config["network_address_prefix"]}"]

  subnet {
    name           = "subnet1"
    address_prefix = "${var.config["network_subnet_prefix"]}"
    security_group = "${azurerm_network_security_group.network_security_group.id}"
  }

  tags {
    Source = "Azure Quickstarts for Terraform"
  }
}

resource "azurerm_network_interface" "network_interface" {
  name                = "${var.config["namespace"]}-nic"
  resource_group_name = "${azurerm_resource_group.resource_group.name}"
  location            = "${var.resource_group_location}"

  ip_configuration {
    name                          = "ipconfig1"
    subnet_id                     = "${azurerm_resource_group.resource_group.id}/providers/Microsoft.Network/virtualNetworks/${var.config["namespace"]}-vnet/subnets/subnet1"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = "${azurerm_public_ip.public_ip.id}"
  }

  # Required due to non-Terraform referencing in ip_configuration of subnet from virtual network resource
  depends_on = ["azurerm_virtual_network.virtual_network"]

  tags {
    Source = "Azure Quickstarts for Terraform"
  }
}

resource "azurerm_virtual_machine" "virtual_machine" {
  name                = "${var.config["namespace"]}-vm"
  resource_group_name = "${azurerm_resource_group.resource_group.name}"
  location            = "${var.resource_group_location}"

  network_interface_ids = ["${azurerm_network_interface.network_interface.id}"]
  vm_size               = "${var.config["vm_size"]}"

  os_profile {
    computer_name  = "${var.config["namespace"]}-vm"
    admin_username = "${var.admin_username}"
    admin_password = ""
  }

  os_profile_linux_config {
    disable_password_authentication = true

    ssh_keys {
      path     = "/home/${var.admin_username}/.ssh/authorized_keys"
      key_data = "${var.admin_ssh_publickey}"
    }
  }

  storage_image_reference {
    publisher = "${var.config["vm_image_publisher"]}"
    offer     = "${var.config["vm_image_offer"]}"
    sku       = "${var.ubuntu_os_version_allowedvalues[var.ubuntu_os_version]}"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osdisk1"
    managed_disk_type = "${var.config["managed_disk_storage_account_type"]}"
    create_option     = "FromImage"
  }

  boot_diagnostics {
    enabled     = "true"
    storage_uri = "${azurerm_storage_account.diagnostics_storage_account.primary_blob_endpoint}"
  }

  tags {
    Source = "Azure Quickstarts for Terraform"
  }
}

resource "azurerm_virtual_machine_extension" "docker_extension" {
  name                = "${var.config["namespace"]}-vmextension"
  resource_group_name = "${azurerm_resource_group.resource_group.name}"
  location            = "${var.resource_group_location}"

  virtual_machine_name       = "${azurerm_virtual_machine.virtual_machine.name}"
  publisher                  = "Microsoft.Azure.Extensions"
  type                       = "DockerExtension"
  type_handler_version       = "1.1"
  auto_upgrade_minor_version = "true"
  settings                   = "{}"

  tags {
    Source = "Azure Quickstarts for Terraform"
  }
}

output "ssh_command" {
  value = "ssh ${var.admin_username}@${azurerm_public_ip.public_ip.fqdn}"
}
