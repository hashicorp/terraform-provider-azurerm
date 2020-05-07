variable "prefix" {
  description = "The prefix which should be used for all resources in this example"
}

variable "location" {
  description = "The Azure Region in which all resources in this example should be created."
}

variable "resource_group_name" {
  description = "The name of the resource group"
}

variable "subnet_id" {
  description = "Subnet ID for the Domain Controllers"
}

variable "active_directory_domain_name" {
  description = "the domain name for Active Directory, for example `consoto.local`"
}

variable "admin_username" {
  description = "Username for the Domain Administrator user"
}

variable "admin_password" {
  description = "Password for the Adminstrator user"
}

variable "active_directory_username" {
  description = "The username of an account with permissions to bind machines to the Active Directory Domain"
}

variable "active_directory_password" {
  description = "The password of the account with permissions to bind machines to the Active Directory Domain"
}

locals {
  virtual_machine_name = join("-", [var.prefix, "client"])
  wait_for_domain_command = "while (!(Test-Connection -TargetName ${var.active_directory_domain_name} -Count 1 -Quiet) -and ($retryCount++ -le 360)) { Start-Sleep 10 }"
}