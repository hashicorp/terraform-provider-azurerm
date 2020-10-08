	provider "azurerm" {
		features {}
	  }
	  
	  resource "azurerm_resource_group" "rg" {
		name     = "acctestJIT_RG-4"
		location = "NorthEurope"
	  }
	  
	  # create vm
	  resource "azurerm_virtual_network" "vnet" {
		name                = "acctestJIT_network"
		address_space       = ["10.0.0.0/16"]
		location            = azurerm_resource_group.rg.location
		resource_group_name = azurerm_resource_group.rg.name
	  }
	  
	  resource "azurerm_subnet" "subnet" {
		name                 = "internal"
		resource_group_name  = azurerm_resource_group.rg.name
		virtual_network_name = azurerm_virtual_network.vnet.name
		address_prefixes     = ["10.0.2.0/24"]
	  }
	  
	  resource "azurerm_subnet_network_security_group_association" "example" {
  		subnet_id                 = azurerm_subnet.subnet.id
  		network_security_group_id = azurerm_network_security_group.vmnsg.id
		}

	  resource "azurerm_network_interface" "inet" {
		name                = "acctestJIT-nic"
		location            = azurerm_resource_group.rg.location
		resource_group_name = azurerm_resource_group.rg.name
	  
		ip_configuration {
		  name                          = "testconfiguration1"
		  subnet_id                     = azurerm_subnet.subnet.id
		  private_ip_address_allocation = "Dynamic"
      public_ip_address_id          = azurerm_public_ip.public_vm_ip.id

		}
	  }
	  
	  resource "azurerm_virtual_machine" "main" {
		name                             = "acctestJIT-vm1"
		location                         = azurerm_resource_group.rg.location
		resource_group_name              = azurerm_resource_group.rg.name
		network_interface_ids            = [azurerm_network_interface.inet.id]
		vm_size                          = "Standard_A1_v2"
		delete_os_disk_on_termination    = true
		delete_data_disks_on_termination = true

		storage_image_reference {
		  publisher = "Canonical"
		  offer     = "UbuntuServer"
		  sku       = "16.04-LTS"
		  version   = "latest"
		}
		storage_os_disk {
		  name              = "disk1"
		  caching           = "ReadWrite"
		  create_option     = "FromImage"
		  managed_disk_type = "Standard_LRS"
		}
		os_profile {
		  computer_name  = "hostname"
		  admin_username = "tester"
		  admin_password = "12345&tester"
		}
		os_profile_linux_config {
		  disable_password_authentication = false
		}
	  }

resource "azurerm_network_security_group" "vmnsg"{
	name = "acctestnsg"
	location            = azurerm_resource_group.rg.location
	resource_group_name = azurerm_resource_group.rg.name
    security_rule {
	name						="acctest_secrule"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
    }

}

resource "azurerm_public_ip" "public_vm_ip" {
  name                = "acctest_vm_ip"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  allocation_method   = "Dynamic"

}

resource "azurerm_jit_network_access_policies" "jit-4" {
	name = "default-4"
	asc_location = azurerm_resource_group.rg.location
	resource_group_name = azurerm_resource_group.rg.name
	kind = "Basic"
	virtual_machines {
		name = azurerm_virtual_machine.main.name
		ports {
			port = 22
			protocol = "*"
			allowed_source_address_prefix = "*"
			max_request_access_duration = "PT3H"
		}
		
	}
	depends_on = [azurerm_resource_group.rg,azurerm_network_security_group.vmnsg]
  }