terraform {
  required_version = ">= 1.6"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
  }
}

# ---------- Providers ----------
provider "azurerm" {
  features {}
}

# variables
variable "resource_group_name" {
  type        = string
  default     = "example-iotoperations"
  description = "Resource Group name to house Azure resources."
}

variable "location" {
  type        = string
  default     = "West Europe"
  description = "Azure location/region."
}

# Optional: tags
variable "tags" {
  type = map(string)
  default = {
    Environment = "Dev"
    Owner       = "team"
  }
}

# -Data (useful for subscription, tenant, object ids)
data "azurerm_client_config" "current" {}

# Resources
resource "azurerm_resource_group" "example" {
  name     = var.resource_group_name
  location = var.location
  tags     = var.tags
}

# ---------- Outputs ----------
output "resource_group_id" {
  value       = azurerm_resource_group.example.id
  description = "ID of the created resource group."
}

