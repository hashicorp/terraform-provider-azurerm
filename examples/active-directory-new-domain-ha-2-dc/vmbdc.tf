#Create Backup Domain Controller NIC
resource "azurerm_network_interface" "bdc_network_interface" {
  depends_on          = ["azurerm_virtual_machine_extension.create_ad_forest_extension", "azurerm_virtual_network.adha_vnet_with_dns"]
  name                = "${var.config["bdc_nic_name"]}"
  resource_group_name = "${azurerm_resource_group.quickstartad.name}"
  location            = "${azurerm_resource_group.quickstartad.location}"
  dns_servers         = ["${var.config["pdc_nic_ip_address"]}"]
  
  ip_configuration {
    name                                    = "ipconfig1"
    subnet_id                               = "${azurerm_subnet.ad_subnet.id}"
    private_ip_address_allocation           = "Static"
    private_ip_address                      = "${var.config["bdc_nic_ip_address"]}"
    load_balancer_backend_address_pools_ids = ["${azurerm_lb_backend_address_pool.adha-lb.id}"]
    load_balancer_inbound_nat_rules_ids     = ["${azurerm_lb_nat_rule.bdc_rdp.id}"]
  }
}

#Create Backup Domain Controller
resource "azurerm_virtual_machine" "bdc_virtual_machine" {
  name                  = "${var.config["bdc_vm_name"]}"
  resource_group_name   = "${azurerm_resource_group.quickstartad.name}"
  location              = "${azurerm_resource_group.quickstartad.location}"
  availability_set_id   = "${azurerm_availability_set.adha_availability_set.id}"
  network_interface_ids = ["${azurerm_network_interface.bdc_network_interface.id}"]
  vm_size               = "${var.config["vm_size"]}"

  os_profile {
    computer_name  = "${var.config["bdc_vm_name"]}"
    admin_username = "${var.admin_username}"
    admin_password = "${var.admin_password}"
  }

  storage_image_reference {
    publisher = "${var.config["vm_image_publisher"]}"
    offer     = "${var.config["vm_image_offer"]}"
    sku       = "${var.config["vm_image_sku"]}"
    version   = "latest"
  }

  storage_os_disk {
    name          = "osdisk"
    vhd_uri       = "${azurerm_storage_account.adha_storage_account.primary_blob_endpoint}${azurerm_storage_container.bdc_storage_container.name}/osdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  storage_data_disk {
    name          = "${var.config["bdc_vm_name"]}-data-disk1"
    vhd_uri       = "${azurerm_storage_account.adha_storage_account.primary_blob_endpoint}${azurerm_storage_container.bdc_storage_container.name}/datadisk1.vhd"
    caching       = "ReadWrite"
    disk_size_gb  = "${var.config["ad_data_disk_size"]}"
    create_option = "Empty"
    lun           = 0
  }
}

#Prepare BDC
resource "azurerm_virtual_machine_extension" "prepare_bdc_extension" {
  name                       = "${format("%s-PrepareBDC", var.config["bdc_vm_name"])}"
  resource_group_name        = "${azurerm_resource_group.quickstartad.name}"
  location                   = "${azurerm_resource_group.quickstartad.location}"

  virtual_machine_name       = "${azurerm_virtual_machine.bdc_virtual_machine.name}"
  publisher                  = "Microsoft.Powershell"
  type                       = "DSC"
  type_handler_version       = "2.21"
  auto_upgrade_minor_version = "true"
  settings = <<SETTINGS
		{
			"ModulesUrl": "${var.config["asset_location"]}${var.config["prepare_bdc_script_path"]}",
			"ConfigurationFunction": "${var.config["ad_bdc_prepare_function"]}\\PrepareADBDC",
      "Properties": {
          "DNSServer": "${var.config["pdc_nic_ip_address"]}"
        }
		}
  SETTINGS
}

#Configure BDC: note that due to a limitation on VM Extensions (only one extension per VM), 
# the name of this extension must be identical to the "prepare_bdc_extension" script
resource "azurerm_virtual_machine_extension" "configure_bdc_extension" {
  depends_on                 = ["azurerm_virtual_machine_extension.prepare_bdc_extension"]
  name                       = "${format("%s-PrepareBDC", var.config["bdc_vm_name"])}"
  resource_group_name        = "${azurerm_resource_group.quickstartad.name}"
  location                   = "${azurerm_resource_group.quickstartad.location}"

  virtual_machine_name       = "${azurerm_virtual_machine.bdc_virtual_machine.name}"
  publisher                  = "Microsoft.Powershell"
  type                       = "DSC"
  type_handler_version       = "2.21"
  auto_upgrade_minor_version = "true"
  settings = <<SETTINGS
		{
			"ModulesUrl": "${var.config["asset_location"]}${var.config["configure_bdc_script_path"]}",
			"ConfigurationFunction": "${var.config["ad_bdc_config_function"]}\\ConfigureADBDC",
      "Properties": {
        "DomainName": "${var.domain_name}",
        "AdminCreds": {
            "UserName": "${var.admin_username}",
            "Password": "PrivateSettingsRef:AdminPassword"
        }
      }
		}
  SETTINGS

  protected_settings         = <<SETTINGS
    {
      "Items": {
        "AdminPassword": "${var.admin_password}"
      }
    }
  SETTINGS
}

#Finally, add BDC DNS to vNet
resource "azurerm_virtual_network" "adha_vnet_with_bdc_dns" {
  depends_on          = ["azurerm_virtual_machine_extension.configure_bdc_extension"]
  name                = "${var.config["vnet_name"]}"
  resource_group_name = "${azurerm_resource_group.quickstartad.name}"
  location            = "${azurerm_resource_group.quickstartad.location}"
  address_space       = ["${var.config["vnet_address_range"]}"]
  dns_servers         = ["${var.config["pdc_nic_ip_address"]}", "${var.config["bdc_nic_ip_address"]}"]
  subnet {
    name                = "${var.config["subnet_name"]}"
    address_prefix       = "${var.config["subnet_address_range"]}"
  }
}
