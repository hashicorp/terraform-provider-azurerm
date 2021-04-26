variable "location" {
  description = "The Azure location where all resources in this example should be created."
}

variable "prefix" {
  description = "The prefix used for all resources used by this NetApp Volume"
}

variable "password" {
  description = "User Password to create ActiveDirectory object"
}

variable "subnet_id" {
  description = "Subnet id (delegated to Microsoft.Netapp/volumes) to attach the volume to"
}
