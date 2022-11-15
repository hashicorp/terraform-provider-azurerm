provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "${var.prefix}-sbnamespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_namespace_authorization_rule" "example" {
  name         = "${var.prefix}-sbnauth"
  namespace_id = azurerm_servicebus_namespace.example.id
  send         = true
  listen       = true
  manage       = true
}

resource "azurerm_servicebus_topic" "example" {
  name                = "${var.prefix}-sbtopic"
  namespace_id        = azurerm_servicebus_namespace.example.id
  enable_partitioning = true
}

resource "azurerm_servicebus_subscription" "example" {
  name               = "${var.prefix}-sbsubscription"
  topic_id           = azurerm_servicebus_topic.example.id
  max_delivery_count = 1
}

resource "azurerm_servicebus_queue" "example" {
  name                = "${var.prefix}-sbqueue"
  namespace_id        = azurerm_servicebus_namespace.example.id
  enable_partitioning = true
}
