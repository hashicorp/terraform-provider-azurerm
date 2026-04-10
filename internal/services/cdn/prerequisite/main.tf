provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "dns" {
  name     = var.resource_group_name
  location = var.location
}

resource "azurerm_dns_zone" "parent" {
  name                = var.dns_zone_name
  resource_group_name = azurerm_resource_group.dns.name

  tags = var.tags
}
