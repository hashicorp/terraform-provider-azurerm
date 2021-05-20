variable "admin" {
  description = "Virtual Machine Admin Username"
  default = "adminuser"
}

variable "adminPassword" {
  description = "VM Admin Password"
  default = "P@$$w0rd1234!"
}
variable "location" {
  description = "The Azure Region in which all resources in this example should be provisioned"
  default = "West US"
}

variable "storageaccount" {
  description = "Name of the storage account"
  default = "exampleblobname"
}
