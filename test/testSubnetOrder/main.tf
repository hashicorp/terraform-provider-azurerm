provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "xiaxintestRG-subnet"
  location = "east us"
}

resource "azurerm_virtual_network" "test" {
  name                = "xiaxintestvirtnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.AzureCosmosDB", "Microsoft.EventHub", "Microsoft.KeyVault", "Microsoft.ServiceBus","Microsoft.Sql", "Microsoft.Storage"]
}
