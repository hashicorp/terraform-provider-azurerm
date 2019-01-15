variable "resource_group" {
  description = "The name of the resource group in which to create the virtual network."
}
variable "rg_prefix" {
  description = "The shortened abbreviation to represent your resource group that will go on the front of some resources."
}
variable "hostname" {
  description = "VM name referenced also in storage-related names."
}
variable "dns_name" {
  description = " Label for the Domain Name. Will be used to make up the FQDN. If a domain name label is specified, an A DNS record is created for the public IP in the Microsoft Azure DNS system."
}
variable "location" {
  description = "The location/region where the virtual network is created. Changing this forces a new resource to be created."
}
variable "virtual_network_name" {
  description = "The name for the virtual network."
}
variable "address_space" {
  description = "The address space that is used by the virtual network. You can supply more than one address space. Changing this forces a new resource to be created."
}
variable "subnet_prefix" {
  description = "The address prefix to use for the subnet."
}
variable "vm_size" {
  description = "Specifies the size of the virtual machine."
}
variable "admin_username" {
  description = "administrator user name"
}
variable "admin_password" {
  description = "administrator password (recommended to disable password auth)"
}