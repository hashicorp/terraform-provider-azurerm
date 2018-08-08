resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_storage_account" "example" {
  name                     = "${var.prefix}storacc"
  resource_group_name      = "${azurerm_resource_group.example.name}"
  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "random_id" "example" {
  keepers = {
    azi_id = 1
  }

  byte_length = 8
}

resource "azurerm_app_service_plan" "example" {
  # App Service Plan name's need to be globally unique - so we suffix a random ID
  name                = "azure-functions-test-service-plan"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  kind                = "FunctionApp"

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}

resource "azurerm_function_app" "example" {
  # Function App name's need to be globally unique - so we suffix a random ID
  name                      = "${var.prefix}-function"
  location                  = "${azurerm_resource_group.example.location}"
  resource_group_name       = "${azurerm_resource_group.example.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.example.id}"
  storage_connection_string = "${azurerm_storage_account.example.primary_connection_string}"
}
