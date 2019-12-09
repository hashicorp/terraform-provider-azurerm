resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_spring_cloud" "example" {
  name                     = "${var.prefix}-sc"
  resource_group           = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location

  tags = {
    env = "test"
  }
}

resource "azurerm_spring_cloud_config_server" "example" {
  spring_cloud_id = "${azurerm_spring_cloud.example.id}"

  uri = "https://github.com/Azure-Samples/piggymetrics"
  label = "config"
}

