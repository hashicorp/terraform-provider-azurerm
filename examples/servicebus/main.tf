resource "azurerm_resource_group" "example" {
  name     = "${var.resource_group}"
  location = "${var.location}"
}

resource "random_integer" "ri" {
  min = 10000
  max = 99999
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "tfex-servicebus${random_integer.ri.result}"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "standard"
}

resource "azurerm_servicebus_namespace_authorization_rule" "example" {
  name                = "tfex_servicebus${random_integer.ri.result}_authrule"
  namespace_name      = "${azurerm_servicebus_namespace.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  send   = true
  listen = true
  manage = true
}

resource "azurerm_servicebus_topic" "source" {
  name                = "tfex_servicebus${random_integer.ri.result}~topic"
  resource_group_name = "${azurerm_resource_group.example.name}"
  namespace_name      = "${azurerm_servicebus_namespace.example.name}"

  enable_partitioning = true
}

resource "azurerm_servicebus_topic_authorization_rule" "example" {
  name                = "tfex_servicebus${random_integer.ri.result}_authrule"
  namespace_name      = "${azurerm_servicebus_namespace.example.name}"
  topic_name          = "${azurerm_servicebus_topic.source.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  send   = true
  listen = true
  manage = true
}

resource "azurerm_servicebus_topic" "forward_to" {
  name                = "tfex_servicebus${random_integer.ri.result}_forwardto"
  resource_group_name = "${azurerm_resource_group.example.name}"
  namespace_name      = "${azurerm_servicebus_namespace.example.name}"

  enable_partitioning = true
}

resource "azurerm_servicebus_subscription" "example" {
  name                = "tfex_servicebus${random_integer.ri.result}_subscription"
  resource_group_name = "${azurerm_resource_group.example.name}"
  namespace_name      = "${azurerm_servicebus_namespace.example.name}"
  topic_name          = "${azurerm_servicebus_topic.source.name}"
  forward_to          = "${azurerm_servicebus_topic.forward_to.name}"
  max_delivery_count  = 1
}

resource "azurerm_servicebus_queue" "example" {
  name                = "tfex_servicebus_queue"
  resource_group_name = "${azurerm_resource_group.example.name}"
  namespace_name      = "${azurerm_servicebus_namespace.example.name}"

  enable_partitioning = true
}

resource "azurerm_servicebus_queue_authorization_rule" "example" {
  name                = "tfex_servicebus${random_integer.ri.result}_authrule"
  namespace_name      = "${azurerm_servicebus_namespace.example.name}"
  queue_name          = "${azurerm_servicebus_queue.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  send   = true
  listen = true
}
