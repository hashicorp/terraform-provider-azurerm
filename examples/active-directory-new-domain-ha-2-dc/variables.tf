variable "admin_username" {
  type        = "string"
  description = "User Name for the Virtual Machine"
}

variable "admin_password" {
  type        = "string"
  description = "Password for the Virtual Machine."
}

variable "ad_vm_size" {
  type        = "string"
  description = "VM SKU for the AD Controller"
}

variable "dns_prefix" {
  type        = "string"
  description = "DNS Prefix for the cluster"
}

variable "public_ip_address_name" {
  type        = "string"
  description = "Public IP Address Name"
}

variable "domain_name" {
  type        = "string"
  description = "The Name of the Domain to Create"
}

variable "resource_group_name" {
  description = "Name of the resource group container for all resources"
}

variable "resource_group_location" {
  description = "Resource group location"
}


# Allowed versions of Windows Server
variable "windows_os_version" {
  type = "string"
  default = "2012-R2-Datacenter"
}

variable "namespace" {
  type = "string"
  default = "highlyavailablead"
}

variable "vnet_name" {
  type = "string"
  default = "adVNET"
}

variable "vnet_address_range" {
  type = "string"
  default = "10.0.0.0/16"
}

variable "subnet_name" {
  type = "string"
  default = "adSubnet"
}

variable "subnet_address_range" {
  type = "string"
  default = "10.0.0.0/24"
}

variable "pdc_nic_name" {
  type = "string"
  default = "adPDCNic"
}

variable "pdc_nic_ip_address" {
  type = "string"
  default = "10.0.0.4"
}

variable "bdc_nic_name" {
  type = "string"
  default = "adBDCNic"
}

variable "bdc_nic_ip_address" {
  type = "string"
  default = "10.0.0.5"
}

variable "network_public_ipaddress_name" {
  type = "string"
  default = "adpublicIP"
}

variable "network_public_ipaddress_type" {
  type = "string"
  default = "Static"
}

variable "pdc_vm_name" {
  type = "string"
  default = "adPDC"
}

variable "bdc_vm_name" {
  type = "string"
  default = "adBDC"
}

variable "vm_size" {
  type = "string"
  default = "Standard_DS2_v2"
}

variable "vm_image_publisher" {
  type = "string"
  default = "MicrosoftWindowsServer"
}

variable "vm_image_offer" {
  type = "string"
  default = "WindowsServer"
}

variable "vm_image_sku" {
  type = "string"
  default = "2012-R2-Datacenter"
}

variable "storage_account_type" {
  type = "string"
  default = "Premium_LRS"
}

variable "availability_set_name" {
  type = "string"
  default = "adAvailabiltySet"
}

variable "pdc_rdp_port" {
  type = "string"
  default = "3389"
}

variable "bdc_rdp_port" {
  type = "string"
  default = "13389"
}

variable "asset_location" {
  type = "string"
  default = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/active-directory-new-domain-ha-2-dc"
}

variable "create_pdc_script_path" {
  type = "string"
  default = "/DSC/CreateADPDC.zip"
}

variable "prepare_bdc_script_path" {
  type = "string"
  default = "/DSC/PrepareADBDC.zip"
}

variable "configure_bdc_script_path" {
  type = "string"
  default = "/DSC/ConfigureADBDC.zip"
}

variable "ad_pdc_config_function" {
  type = "string"
  default = "CreateADPDC.ps1"
}

variable "ad_bdc_prepare_function" {
  type = "string"
  default = "PrepareADBDC.ps1"
}

variable "ad_bdc_config_function" {
  type = "string"
  default = "ConfigureADBDC.ps1"
}

variable "ad_loadbalancer_frontend_name" {
  type = "string"
  default = "LBFE"
}

variable "ad_loadbalancer_backend_name" {
  type = "string"
  default = "LBBE"
}

variable "ad_pdc_rdp_nat_name" {
  type = "string"
  default = "adPDCRDP"
}

variable "ad_bdc_rdp_nat_name" {
  type = "string"
  default = "adBDCRDP"
}

variable "ad_data_disk_size" {
  type = "string"
  default = "1000"
}

variable "pdc_storage_container_name" {
  type = "string"
  default = "pdcvhds"
}

variable "bdc_storage_container_name" {
  type = "string"
  default = "bdcvhds"
}
