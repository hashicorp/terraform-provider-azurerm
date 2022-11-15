variable "prefix" {
  description = "The prefix which should be used for all resources in this example"
}

variable "shared_key" {
  description = "This is a shared key used to secure the connection between VNets. Make sure to set it in a secure way!"
}

variable "location_1" {
  description = "The Azure Region in which the primary VNet resources in this module are created."
}

variable "resource_group_name_1" {
  description = "The Azure Resouce Group in which the primary VNet resources in this module are created."
}

variable "vnet_name_1" {
  description = "The name of the primary VNet in which resources in this module are created."
}

variable "gateway_subnet_range_1" {
  description = "The range used in the primary Gateway Subnet."
}

variable "location_2" {
  description = "The Azure Region in which the secondary VNet resources in this module are created."
}

variable "resource_group_name_2" {
  description = "The Azure Resouce Group in which the secondary VNet resources in this module are created."
}

variable "vnet_name_2" {
  description = "The name of the secondary VNet in which resources in this module are created."
}

variable "gateway_subnet_range_2" {
  description = "The range used in the secondary Gateway Subnet."
}