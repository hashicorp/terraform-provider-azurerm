resource "azurerm_resource_group" "rg" {
  name     = "${var.resource_group}"
  location = "${var.location}"
}

resource "random_integer" "ri" {
  min = 10000
  max = 99999
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "tfex-servicebus${random_integer.ri.result}"
  location            = "${var.location}"
  resource_group_name = "${var.resource_group}"
  sku                 = "standard"
  capacity            = 1
}

/*
resource "azurerm_servicebus_topic" "test" {
  name                = "tfex_servicebus${random_integer.ri.result}_topic"
  resource_group_name = "${var.resource_group}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"

  enable_partitioning = true
}

resource "azurerm_servicebus_subscription" "test" {
  name                = "tfex_servicebus${random_integer.ri.result}_subscription"
  resource_group_name = "${var.resource_group}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name          = "${azurerm_servicebus_topic.test.name}"
  forward_to          = "${azurerm_servicebus_topic.forward_to.name}"
  max_delivery_count  = 1
}


resource "azurerm_servicebus_topic" "forward_to" {
  name                = "tfex_servicebus${random_integer.ri.result}_forwardto"
  resource_group_name = "${var.resource_group}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"

  enable_partitioning = true
}
*/

