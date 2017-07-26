variable "resource_group_name" {
  type        = "string"
  description = "Name of the azure resource group."
}

variable "resource_group_location" {
  type        = "string"
  description = "Location of the azure resource group."
}

variable "admin_username" {
  type        = "string"
  description = "User name for authentication to the Virtual Machine"
}

variable "admin_ssh_publickey" {
  type        = "string"
  description = "SSH public key for authentication to the Virtual Machine."
}

variable "dns_name_for_public_ip" {
  type        = "string"
  description = "Unique DNS Name for the Public IP used to access the Virtual Machine."
}

variable "ubuntu_os_version" {
  type        = "string"
  description = "The Ubuntu version for deploying the Docker containers. This will pick a fully patched image of this given Ubuntu version. Allowed values: 16.04.0-LTS 16.10"
}
