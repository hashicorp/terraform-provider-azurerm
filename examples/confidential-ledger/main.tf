provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "eastus"
}

resource "azurerm_confidential_ledger" "example" {
	name                = "example-ledger"
	ledger_type         = "Public"
	location            = azurerm_resource_group.example.location
	resource_group_name = azurerm_resource_group.example.name

	tags = {
		IsExample = "True"
	}
  }