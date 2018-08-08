variable "prefix" {
  description = "The Prefix used for all resources in this example."
}

variable "location" {
  description = "The Azure Region in which the resources in this example should exist."
}

variable "instance_count" {
  description = "Number of VM instances (100 or less)."
  default     = "3"
}

variable "admin_username" {
  description = "The Admin Username used for all VM's in this Scale Set."
}

variable "admin_password" {
  description = "The Admin Password used for all VM's in this Scale Set."
}
