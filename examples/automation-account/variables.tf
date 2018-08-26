variable "resource_group_name" {
  type        = "string"
  description = "Name of the azure resource group."
  default     = "tfex-automation_account"
}

variable "resource_group_location" {
  type        = "string"
  description = "Location of the azure resource group."
  default     = "west europe"
}
