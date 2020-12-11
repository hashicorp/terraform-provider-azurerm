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
  kind                = "xenon"
  is_xenon            = true

  sku {
    tier = "PremiumContainer"
    size = "PC2"
  }
}

resource "azurerm_app_service" "example" {
  name                = "${var.prefix}-appservice"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  app_service_plan_id = "${azurerm_app_service_plan.example.id}"

  site_config {
    windows_fx_version = "DOCKER|mcr.microsoft.com/azure-app-service/samples/aspnethelloworld:latest"
  }

  app_settings = {
    "DOCKER_REGISTRY_SERVER_URL" = "https://mcr.microsoft.com",
    "DOCKER_REGISTRY_SERVER_USERNAME" = "",
    "DOCKER_REGISTRY_SERVER_PASSWORD" = "",
  }
}
