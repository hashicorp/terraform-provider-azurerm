provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "${var.prefix}-ehnamespace"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "Standard"
  capacity            = 2

  tags = {
    environment = "Examples"
  }
}

resource "azurerm_eventhub_namespace_authorization_rule" "example" {
  name                = "${var.prefix}-nsauth-rule"
  namespace_name      = "${azurerm_eventhub_namespace.example.name}"
  eventhub_name       = "${azurerm_eventhub.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  listen = true
  send   = true
  manage = false
}

resource "azurerm_eventhub" "example" {
  name                = "${var.prefix}-eh1"
  namespace_name      = "${azurerm_eventhub_namespace.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  partition_count   = 2
  message_retention = 1
}

resource "azurerm_eventhub_authorization_rule" "test" {
  name                = "${var.prefix}-enauth-rule"
  namespace_name      = "${azurerm_eventhub_namespace.example.name}"
  eventhub_name       = "${azurerm_eventhub.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  listen = true
  send   = true
  manage = true
}

resource "azurerm_eventhub_consumer_group" "example" {
  name                = "${var.prefix}-ehcg"
  namespace_name      = "${azurerm_eventhub_namespace.example.name}"
  eventhub_name       = "${azurerm_eventhub.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  user_metadata       = "some-meta-data"
}
