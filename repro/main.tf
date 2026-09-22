terraform {
  required_providers {
    azurerm = {
      version = ">=5"
    }
  }
}

provider "azurerm" {
  features {}
}

data "azurerm_resource_group" "test" {
  name = "test"

  timeouts {
  #   read = "10s"
  }
}
