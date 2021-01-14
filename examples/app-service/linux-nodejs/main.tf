provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "main" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_app_service_plan" "main" {
  name                = "${var.prefix}-asp"
  location            = "${azurerm_resource_group.main.location}"
  resource_group_name = "${azurerm_resource_group.main.name}"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "main" {
  name                = "${var.prefix}-appservice"
  location            = "${azurerm_resource_group.main.location}"
  resource_group_name = "${azurerm_resource_group.main.name}"
  app_service_plan_id = "${azurerm_app_service_plan.main.id}"

  site_config {
    linux_fx_version = "NODE|10.14"
  }
}
