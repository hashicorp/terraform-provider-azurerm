resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "${var.prefix}-namespace"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "standard"

  tags {
    source = "terraform"
  }
}

resource "azurerm_servicebus_topic" "example" {
  name                = "${var.prefix}-topic"
  resource_group_name = "${azurerm_resource_group.example.name}"
  namespace_name      = "${azurerm_servicebus_namespace.example.name}"
  enable_partitioning = true
}

resource "azurerm_servicebus_subscription" "example" {
  name                = "${var.prefix}-subscription"
  resource_group_name = "${azurerm_resource_group.example.name}"
  namespace_name      = "${azurerm_servicebus_namespace.example.name}"
  topic_name          = "${azurerm_servicebus_topic.example.name}"
  max_delivery_count  = 1
}

resource "azurerm_servicebus_subscription_rule" "example" {
  name                = "${var.prefix}-rule"
  resource_group_name = "${azurerm_resource_group.example.name}"
  namespace_name      = "${azurerm_servicebus_namespace.example.name}"
  topic_name          = "${azurerm_servicebus_topic.example.name}"
  subscription_name   = "${azurerm_servicebus_subscription.example.name}"
  filter_type         = "CorrelationFilter"

  correlation_filter {
    correlation_id = "high"
    label          = "red"
  }
}
