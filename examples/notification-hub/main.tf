provider "azurerm" {
  features {}
}

resource "azurerm_notification_hub" "example" {
  name                = "${var.prefix}-nh"
  resource_group_name = azurerm_resource_group.example.name
  namespace_name      = azurerm_notification_hub_namespace.example.name
  location            = azurerm_notification_hub_namespace.example.location
}

resource "azurerm_notification_hub_namespace" "example" {
  name                = "${var.prefix}-nh-ns"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Standard"
  namespace_type      = "NotificationHub"
}
