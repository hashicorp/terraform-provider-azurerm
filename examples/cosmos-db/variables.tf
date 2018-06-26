variable "resource_group_name" {
  type        = "string"
  description = "Name of the azure resource group."
  default     = "tfex-cosmos-db"
}

variable "resource_group_location" {
  type        = "string"
  description = "Location of the azure resource group."
  default     = "westus"
}

variable "failover_location" {
  type        = "string"
  description = "Location of the failover instance."
  default     = "eastus"
}
