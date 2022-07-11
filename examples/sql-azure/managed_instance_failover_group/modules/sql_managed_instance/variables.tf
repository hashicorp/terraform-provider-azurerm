variable "prefix" {
  description = "The prefix which is used for most resources in this module"
}

variable "location" {
  description = "The Azure Region in which all resources in this module should be created."
}

variable "resource_group_name" {
  description = "The Azure Resource Group in which all resources in this module should be created."
}

variable "vnet_range" {
  description = "The VNet range for the VNet containing the Managed Instance."

  default = "10.0.0.0/16"
}

variable "subnet_range" {
  description = "The subnet range for the subnet containing the Managed Instance."

  default = "10.0.0.0/24"
}

variable "dns_zone_partner_id" {
  description = "The possible DNS Zone Partner Managed Instance."

  default = ""
}