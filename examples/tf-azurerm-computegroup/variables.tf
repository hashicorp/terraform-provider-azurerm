variable "resource_group_name" {
    description = "The name of the resource group in which the resources will be created"
    default = "vmssrg"
}

variable "location" { 
    description = "The location where the resources will be created"
    default = "West US"
}

variable "vm_size" {
    default = "Standard_A0"
    description = "Size of the Virtual Machine based on Azure sizing"
}

variable "admin_username" {
    description = "The admin username of the VMSS being deployed"
    default = "azureuser"
 }

variable "admin_password" {
    default = ""
}

variable "ssh_key" {
    description = "Path to the public key to be used for ssh access to the VM"
    default = "~/.ssh/demo_key.pub"
}

variable "nb_instance" { 
    description = "Specify the number of vm instances"
    default = "1"
}

variable "protocol"{ 
    default = { 
        ssh = ["22", "Tcp", "50000", "50119"]
        http = ["80", "Tcp", "51000", "51119"]
    }
}

variable "os" { 
    type = "map"
    default = {  
        id = ""
        publisher = "Canonical"
        offer = "UbuntuServer"
        sku = "14.04.2-LTS"
        version = "latest"
    }
}

variable "tags"{
    type = "map"
    default = {
        tag1 = "dev"
        tag2 = "demo"
    }
}
