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
  description = "The admin username of the VMSS being deployed"
  default     = "azureuser"
}

variable "admin_password" {
  default = "CpmplexP@ssw0rd"
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

variable "remote_port" { 
    description = "A map of the port that you want to use for remote access [protocol, backend_port]. Frontend port will be automatically generated starting at 50000 and in the output."
    default = {
        ssh = ["Tcp", "22"]
    }
}

variable "lb_port" {
    description = "Protocols to be used for lb health probes and rules. [frontend_port, protocol, backend_port]"
    default = {
        http = ["80", "Tcp", "80"]
        https = ["443", "Tcp", "443"]
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
