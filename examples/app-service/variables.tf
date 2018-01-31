variable "resource_group_name" {
  type = "string"
}

variable "location" {
  type = "string"
}

variable "app_service_name" {
  type = "string"
}

variable "app_service_plan_sku_tier" {
  type    = "string"
  default = "Basic"  # Basic | Standard | ...
}

variable "app_service_plan_sku_size" {
  type    = "string"
  default = "B1"     # B1 | S1 | ...
}
