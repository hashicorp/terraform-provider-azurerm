variable "subscription_id" {}
variable "client_id" {}
variable "client_secret" {}
variable "tenant_id" {}


variable "GoldenImage" {
  description ="name of the existing Golden Image"
}

variable "RgOfGoldenImage" {
  description = "name of the existing RG where the Image is - must be in the same region!!!"
}

variable "resource_group" {
  description = "The name of the resource group in which to create the virtual network."
}

variable "location" {
  description = "The location/region where the virtual network is created. Changing this forces a new resource to be created."
}

variable "customer_name" {
  description = "The name of the customer. will be the name of the vNet."
}

variable "admin_username" {
  description = "the local user name of the VM"
}

variable "admin_password" {
  description = "the local password name of the VM"
}
