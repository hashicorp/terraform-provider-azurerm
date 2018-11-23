variable "resource_group_name" {
  type        = "string"
  description = "Name of the azure resource group."
  default     = "tfex-container-registry"
}

variable "resource_group_location" {
  type        = "string"
  description = "Location of the azure resource group."
  default     = "westus"
}

variable "sku" {
  type        = "string"
  description = "SKU for the Azure Container Registry"
  default     = "Premium"
}

variable "admin_enabled" {
  description = "Flag that indicates wether an admin account/password should be created or not"
  default     = false
}

variable "georeplication_locations" {
  type        = "list"
  description = "Azure locations where the Container Registry should be geo-replicated - Only with Premium SKU"
  default     = ["eastus"]
}