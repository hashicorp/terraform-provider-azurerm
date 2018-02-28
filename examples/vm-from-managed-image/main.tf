provider "azurerm" {
   subscription_id = "${var.subscription_id}"
   client_id       = "${var.client_id}"
   client_secret   = "${var.client_secret}"
   tenant_id       = "${var.tenant_id}"
 }

#take a pointer to the custom image from our subscription
data "azurerm_image" "search" {
        #name of the existing Image
  name                = "${var.GoldenImage}"   
        #name of the existing RG where the Image is - must be in the same region!!!
  resource_group_name = "${var.RgOfGoldenImage}" 
}

output "image_id" {
  value = "${data.azurerm_image.search.id}"
}


# Create an Azure resource group
resource "azurerm_resource_group" "rg" {
  name     = "${var.resource_group}"
  location = "${var.location}"
}

# Create a virtual network in the resource group
resource "azurerm_virtual_network" "vNet" {
  name                = "${var.customer_name}"
  address_space       = ["172.16.0.0/16"]
  resource_group_name = "${azurerm_resource_group.rg.name}"
  location = "${azurerm_resource_group.rg.location}"
}

# Create Subnet 
resource "azurerm_subnet" "Subnet" {
  name                 = "Subnet"
  virtual_network_name = "${azurerm_virtual_network.vNet.name}"
  resource_group_name  = "${azurerm_resource_group.rg.name}"
  address_prefix       = "172.16.1.0/24"
}  
    
# Create a public IP resource for VMs
resource "azurerm_public_ip" "vm-pip" {
  count 					   = 3
  name                         = "VM-PIP-${count.index}"
  location                     = "${azurerm_resource_group.rg.location}"
  resource_group_name          = "${azurerm_resource_group.rg.name}"
  public_ip_address_allocation = "dynamic" 
} 

# Create a network secuirty group with some rules
resource "azurerm_network_security_group" "nsg" {
  name                = "${var.customer_name}-NSG"
  location            = "${azurerm_resource_group.rg.location}"
  resource_group_name = "${azurerm_resource_group.rg.name}"

  security_rule {
    name                       = "allow_SSH"
    description                = "Allow SSH access"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
  
    security_rule {
    name                       = "allow_RDP"
    description                = "Allow RDP access"
    priority                   = 110
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "3389"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

# Create a network interface for VMs and attach the PIP and the NSG
resource "azurerm_network_interface" "vm-Nic" {
  count				  = 3
  name                = "vm-Nic-${count.index}"
  location            = "${azurerm_resource_group.rg.location}"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  network_security_group_id     = "${azurerm_network_security_group.nsg.id}"
  
  ip_configuration {
    name                          = "Nic-config-${count.index}"
    subnet_id                     = "${azurerm_subnet.Subnet.id}"
    private_ip_address_allocation = "dynamic"
    public_ip_address_id          = "${element(azurerm_public_ip.vm-pip.*.id, count.index)}"
	}
}


##### Create new virtual machine - 3 vms 
resource "azurerm_virtual_machine" "vm" {
  count 			        	= 3
  name                  = "VM-${count.index}"
  location              = "${azurerm_resource_group.rg.location}"
  resource_group_name   = "${azurerm_resource_group.rg.name}"
  network_interface_ids = ["${element(azurerm_network_interface.vm-Nic.*.id, count.index)}"]
  vm_size               = "Standard_E32s_v3"
  delete_os_disk_on_termination = true
  delete_data_disks_on_termination = true

  storage_image_reference {
    id = "${data.azurerm_image.search.id}" 
  }

  storage_os_disk {
    name              = "vm-OS-${count.index}"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Premium_LRS"
    disk_size_gb = "40"
  }

  storage_data_disk {
	name = "vm-Data-Disk-${count.index}"
	managed_disk_type = "Premium_LRS"
	create_option = "Empty"
	lun = 0
	disk_size_gb = "1024"
} 

  os_profile {
    computer_name  = "VM-${count.index}"
    admin_username = "${var.admin_username}"
    admin_password = "${var.admin_password}"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}