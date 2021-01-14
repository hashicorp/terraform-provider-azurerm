provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

# NOTE: the Name used for Redis needs to be globally unique
resource "azurerm_redis_cache" "example" {
  name                = "${var.prefix}-redis"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  capacity            = 0
  family              = "C"
  sku_name            = "Basic"
  enable_non_ssl_port = false
}
