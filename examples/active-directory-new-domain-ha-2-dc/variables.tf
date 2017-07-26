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
