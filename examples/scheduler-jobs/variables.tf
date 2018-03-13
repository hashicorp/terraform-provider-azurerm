variable "resource_group_name" {
  type        = "string"
  description = "Name of the azure resource group."
  default     = "tfex-job_collection"
}

variable "resource_group_location" {
  type        = "string"
  description = "Location of the azure resource group."
  default     = "westus"
}
