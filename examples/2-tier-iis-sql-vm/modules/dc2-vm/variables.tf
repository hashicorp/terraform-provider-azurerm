variable "resource_group_name" {
  description = "The name of the Resource Group where the Domain Controllers resources will be created"
}

variable "location" {
  description = "The Azure Region in which the Resource Group exists"
}

variable dcavailability_set_id {}
variable "prefix" {
  description = "The Prefix used for the Domain Controller's resources"
}

variable "subnet_id" {
  description = "The Subnet ID which the Domain Controller's NIC should be created in"
}

variable "dc2private_ip_address" {}

variable "active_directory_domain" {
  description = "The name of the Active Directory domain, for example `consoto.local`"
}

variable "active_directory_username" {
  description = "The username of an account with permissions to bind machines to the Active Directory Domain"
}

variable "active_directory_password" {
  description = "The password of the account with permissions to bind machines to the Active Directory Domain"
}

variable "active_directory_netbios_name" {
  description = "The netbios name of the Active Directory domain, for example `consoto`"
}

variable "admin_username" {
  description = "The username associated with the local administrator account on the virtual machine"
}

variable "admin_password" {
  description = "The password associated with the local administrator account on the virtual machine"
}

variable "domainadmin_username" {
  description = "The password associated with the local administrator account on the virtual machine"
}