# Unique name for the storage account
resource "random_id" "storage_account_name" {
  keepers = {
    # Generate a new id each time a new resource group is defined
    resource_group = "${azurerm_resource_group.quickstartad.name}"
  }

  byte_length = 8
}

# Allowed versions of Windows Server
variable "windows_os_version_map" {
  type = "map"

  default = {
    "2012-R2-Datacenter" = "2012-R2-Datacenter"
  }
}

variable "config" {
  type = "map"

  default = {
    "namespace"                     = "highlyavailablead"
    "vnet_name"                     = "adVNET"
    "vnet_address_range"            = "10.0.0.0/16"
    "subnet_name"                   = "adSubnet"
    "subnet_address_range"          = "10.0.0.0/24"
    "pdc_nic_name"                  = "adPDCNic"
    "pdc_nic_ip_address"            = "10.0.0.4"
    "bdc_nic_name"                  = "adBDCNic"
    "bdc_nic_ip_address"            = "10.0.0.5"
    "network_public_ipaddress_name" = "adpublicIP" 
    "network_public_ipaddress_type" = "Static"
    "pdc_vm_name"                   = "adPDC"
    "bdc_vm_name"                   = "adBDC"
    "vm_size"                       = "Standard_DS2_v2"
    "vm_image_publisher"            = "MicrosoftWindowsServer"
    "vm_image_offer"                = "WindowsServer"
    "vm_image_sku"                  = "2012-R2-Datacenter"
    "storage_account_type"          = "Premium_LRS"
    "availability_set_name"         = "adAvailabiltySet"
    "pdc_rdp_port"                  = "3389"
    "bdc_rdp_port"                  = "13389"
    "asset_location"                = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/active-directory-new-domain-ha-2-dc"
    "create_pdc_script_path"        = "/DSC/CreateADPDC.zip"
    "prepare_bdc_script_path"       = "/DSC/PrepareADBDC.zip"
    "configure_bdc_script_path"     = "/DSC/ConfigureADBDC.zip"
    "ad_pdc_config_function"        = "CreateADPDC.ps1"
    "ad_bdc_prepare_function"       = "PrepareADBDC.ps1"
    "ad_bdc_config_function"        = "ConfigureADBDC.ps1"
    
    #the variables section from the original template
    "ad_loadbalancer_frontend_name" = "LBFE"
    "ad_loadbalancer_backend_name" = "LBBE"
    "ad_pdc_rdp_nat_name"           = "adPDCRDP"
    "ad_bdc_rdp_nat_name"           = "adBDCRDP"
    "ad_data_disk_size"             = "1000"
    "pdc_storage_container_name"    = "pdcvhds"
    "bdc_storage_container_name"    = "bdcvhds"
  }
}

resource "azurerm_resource_group" "quickstartad" {
  name     = "${var.resource_group_name}"
  location = "${var.resource_group_location}"

  tags {
    Source = "Azure Quickstarts for Terraform"
  }
}

# Using a storage account until managed disks are supported by terraform provider
resource "azurerm_storage_account" "adha_storage_account" {
  name                = "haad${random_id.storage_account_name.hex}"
  resource_group_name = "${azurerm_resource_group.quickstartad.name}"
  location            = "${azurerm_resource_group.quickstartad.location}"

  account_type = "${var.config["storage_account_type"]}"

  tags {
    Source = "Azure Quickstarts for Terraform"
  }
}

# Need a storage containers until managed disks supported by terraform provider
resource "azurerm_storage_container" "pdc_storage_container" {
  name                  = "${var.config["pdc_storage_container_name"]}"
  resource_group_name   = "${azurerm_resource_group.quickstartad.name}"
  storage_account_name  = "${azurerm_storage_account.adha_storage_account.name}"
  container_access_type = "private"
}

resource "azurerm_storage_container" "bdc_storage_container" {
  name                  = "${var.config["bdc_storage_container_name"]}"
  resource_group_name   = "${azurerm_resource_group.quickstartad.name}"
  storage_account_name  = "${azurerm_storage_account.adha_storage_account.name}"
  container_access_type = "private"
}

resource "azurerm_availability_set" "adha_availability_set" {
  name                = "${var.config["availability_set_name"]}"
  resource_group_name = "${azurerm_resource_group.quickstartad.name}"
  location            = "${azurerm_resource_group.quickstartad.location}"
}
