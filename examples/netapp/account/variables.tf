variable "resource_group_name" {
  description = "The name of the resource group the NetApp Account is located in."
}

variable "location" {
  description = "The Azure location where all resources in this example should be created."
  default     = "eastus2"
}

variable "prefix" {
  description = "The prefix used for all resources used by this NetApp Account"
  default     = "example"
}

variable "username" {
  description = "The username of Active Directory"
  default     = "aduser"
}

variable "password" {
  description = "The password of Active Directory"
  default     = "aduserpwd"
}
