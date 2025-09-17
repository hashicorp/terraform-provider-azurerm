provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-iotoperations"
  location = "West Europe"
}

# Example IoT Operations Instance
resource "azurerm_iotoperations_instance" "example" {
  name                = "example-iotinstance"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  # Optional: example extended location usage
  # extended_location_name = "your-extended-location"
  # extended_location_type = "CustomLocation"

  description = "Example IoT Operations instance managed by Terraform"
  version     = "1.0.0"

  tags = {
    Environment = "Dev"
    Owner       = "team"
  }
}