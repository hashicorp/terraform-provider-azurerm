variable "resource_group_name" {
  description = "The name of the resource group"
  type        = string
  default     = "example-netapp-protocol-conversion-rg"
}

variable "location" {
  description = "The Azure region for the resources"
  type        = string
  default     = "westus3"
}

variable "prefix" {
  description = "The prefix for all resources"
  type        = string
  default     = "example"
}

variable "protocol_type" {
  description = "The NFS protocol type (NFSv3 or NFSv4.1)"
  type        = string
  default     = "NFSv3"

  validation {
    condition     = contains(["NFSv3", "NFSv4.1"], var.protocol_type)
    error_message = "Protocol type must be either NFSv3 or NFSv4.1."
  }
}