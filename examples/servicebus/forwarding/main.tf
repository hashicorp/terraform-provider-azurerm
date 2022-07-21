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

resource "azurerm_servicebus_topic" "source" {
  name                = "${var.prefix}-sbt-source"
  namespace_id        = azurerm_servicebus_namespace.example.id
  enable_partitioning = true
}

resource "azurerm_servicebus_topic_authorization_rule" "example" {
  name     = "${var.prefix}-sbtauth-source"
  topic_id = azurerm_servicebus_topic.source.id
  send     = true
  listen   = true
  manage   = true
}

resource "azurerm_servicebus_topic" "destination" {
  name                = "${var.prefix}-sbt-destination"
  namespace_id        = azurerm_servicebus_namespace.example.id
  enable_partitioning = true
}

resource "azurerm_servicebus_subscription" "example" {
  name               = "${var.prefix}-sbsubscription"
  topic_id           = azurerm_servicebus_topic.source.id
  forward_to         = azurerm_servicebus_topic.destination.name
  max_delivery_count = 1
}
