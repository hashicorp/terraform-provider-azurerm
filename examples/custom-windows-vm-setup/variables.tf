variable "prefix" {
  description = "The prefix to use for all resources"
  default     = "winvm2022"
}

variable "location" {
  description = "The Azure region to deploy resources in"
  default     = "East US"
}

variable "subscription_id" {
  description = "Azure Subscription ID"
}

variable "tenant_id" {
  description = "Azure Tenant ID"
}

variable "client_id" {
  description = "Azure Client ID"
}

variable "client_secret" {
  description = "Azure Client Secret"
}

variable "admin_password" {
  description = "Admin password for the Windows VM"
  type        = string
}





