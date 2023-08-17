terraform {
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "azurerm" {
  # Configuration options
  features {

  }
}

# resource "azurerm_resource_group" "example" {
#   name     = "azurerm-resources"
#   location = "West Europe"
# }

resource "azurerm_storage_account" "example" {
  name                     = "storageaccountnamex"
  resource_group_name      = "azurerm-resources"
  location                 = "westeurope"
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}
