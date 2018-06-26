variable "resource_group_name" {
  type        = "string"
  description = "Name of the azure resource group."
  default     = "tfex-appserviceplan"
}

variable "location" {
  type        = "string"
  description = "Location of the azure resource group."
  default     = "westus"
}

variable "app_service_plan_sku_tier" {
  type        = "string"
  description = "SKU tier of the App Service Plan"
  default     = "Basic"                            # Basic | Standard | ...
}

variable "app_service_plan_sku_size" {
  type        = "string"
  description = "SKU size of the App Service Plan"
  default     = "B1"                               # B1 | S1 | ...
}
