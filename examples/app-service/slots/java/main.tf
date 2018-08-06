resource "random_id" "example" {
  keepers = {
    azi_id = 1
  }

  byte_length = 4
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_app_service_plan" "example" {
  name                = "${var.prefix}-appserviceplan"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "example" {
  name                = "${var.prefix}-appservice${random_id.example.hex}"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  app_service_plan_id = "${azurerm_app_service_plan.example.id}"

  site_config {
    java_version           = "1.8"
    java_container         = "JETTY"
    java_container_version = "9.3"
  }
}

resource "azurerm_app_service_slot" "example" {
  name                = "first"
  app_service_name    = "${azurerm_app_service.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  app_service_plan_id = "${azurerm_app_service_plan.example.id}"

  site_config {
    java_version           = "1.8"
    java_container         = "JETTY"
    java_container_version = "9.3"
  }
}
