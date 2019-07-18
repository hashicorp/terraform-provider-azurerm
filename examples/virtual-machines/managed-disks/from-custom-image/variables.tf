variable "prefix" {
  description = "The Prefix used for all resources in this example"
}

variable "location" {
  description = "The Azure Region in which the resources in this example should exist"
}

variable "custom_image_resource_group_name" {
  description = "The name of the Resource Group in which the Custom Image exists."
}

variable "custom_image_name" {
  description = "The name of the Custom Image to provision this Virtual Machine from."
}
