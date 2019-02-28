resource "azurerm_resource_group" "example" {
  name     = "${var.resource_group}"
  location = "${var.location}"
}

resource "random_integer" "ri" {
  min = 10000
  max = 99999
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "tfex-eventhub${random_integer.ri.result}-namespace"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  sku      = "Standard"
  capacity = 2

  tags = {
    environment = "Examples"
  }
}

resource "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = "tfex-eventhub-namespace-authrule"
  namespace_name      = "${azurerm_eventhub_namespace.example.name}"
  eventhub_name       = "${azurerm_eventhub.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  listen = true
  send   = true
  manage = false
}

resource "azurerm_eventhub" "example" {
  name                = "tfex-eventhub${random_integer.ri.result}"
  namespace_name      = "${azurerm_eventhub_namespace.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  partition_count   = 2
  message_retention = 1
}

resource "azurerm_eventhub_authorization_rule" "test" {
  name                = "tfex-eventhub-authrule"
  namespace_name      = "${azurerm_eventhub_namespace.example.name}"
  eventhub_name       = "${azurerm_eventhub.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  listen = true
  send   = true
  manage = true
}

resource "azurerm_eventhub_consumer_group" "example" {
  name                = "tfex-eventhub${random_integer.ri.result}-consumer"
  namespace_name      = "${azurerm_eventhub_namespace.example.name}"
  eventhub_name       = "${azurerm_eventhub.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  user_metadata       = "some-meta-data"
}
