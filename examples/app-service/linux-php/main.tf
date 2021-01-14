provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_app_service_plan" "example" {
  name                = "${var.prefix}-asp"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "example" {
  name                = "${var.prefix}-appservice"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  app_service_plan_id = "${azurerm_app_service_plan.example.id}"

  site_config {
    linux_fx_version = "PHP|7.0"
    scm_type         = "LocalGit"
  }
}
