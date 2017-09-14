variable "resource_group_name" {
  description = "The name of the resource group in which the resources will be created"
  default     = "vmssrg"
}

variable "location" {
  description = "The location where the resources will be created"
  default     = "West US"
}

variable "vm_size" {
  default     = "Standard_A0"
  description = "Size of the Virtual Machine based on Azure sizing"
}

variable "admin_username" {
  description = "The admin username of the VMSS that will be deployed"
  default     = "azureuser"
}

variable "admin_password" {
  description = "The admin password to be used on the VMSS that will be deployed. The password must meet the complexity requirements of Azure"
  default = ""
}

variable "ssh_key" {
  description = "Path to the public key to be used for ssh access to the VM"
  default     = "~/.ssh/id_rsa.pub"
}

variable "nb_instance" {
  description = "Specify the number of vm instances"
  default     = "1"
}

variable "vm_os_publisher" {
  description = "The name of the publisher of the image that you want to deploy"
  default = "Canonical"
}

variable "vm_os_offer" {
  description = "The name of the offer of the image that you want to deploy"
  default = "UbuntuServer"
}

variable "vm_os_sku" {
  description = "The sku of the image that you want to deploy"
  default = "14.04.2-LTS"
}

variable "vm_os_id" {
  description = "The ID of the image that you want to deploy if you are using a custom image."
  default = ""
}

variable "lb_port" {
    description = "Protocols to be used for lb health probes and rules. [frontend_port, protocol, backend_port]"
    default = {
    }
}

variable "tags" {
  type = "map"
  description = "A map of the tags to use on the resources that are deployed with this module."
  default = {
    tag1 = ""
    tag2 = ""
  }
}
