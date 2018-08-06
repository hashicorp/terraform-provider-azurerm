resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_app_service_plan" "example" {
  name                = "${var.prefix}-appserviceplan"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  kind                = "Linux"

  sku {
    tier = "Standard"
    size = "S1"
  }

  properties {
    reserved = true
  }
}
