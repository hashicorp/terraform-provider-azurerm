provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_eventhub_cluster" "example" {
  name                = "${var.prefix}-ehcluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Dedicated_1"
}

resource "azurerm_eventhub_namespace" "example" {
  name                 = "${var.prefix}-ehnamespace"
  location             = azurerm_resource_group.example.location
  resource_group_name  = azurerm_resource_group.example.name
  sku                  = "Standard"
  dedicated_cluster_id = azurerm_eventhub_cluster.example.id
}

resource "azurerm_eventhub" "example" {
  name                = "${var.prefix}-eventhub"
  resource_group_name = azurerm_resource_group.example.name
  namespace_name      = azurerm_eventhub_namespace.example.name
  partition_count     = 40
  message_retention   = 1
}
