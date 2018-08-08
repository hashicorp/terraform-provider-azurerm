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

resource "azurerm_app_service_plan" "test" {
  # App Service Plan name's need to be globally unique - so we suffix a random ID
  name                = "${var.prefix}-asp${random_id.example.id}"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  # Function App name's need to be globally unique - so we suffix a random ID
  name                      = "${var.prefix}-function${random_id.example.id}"
  location                  = "${azurerm_resource_group.example.location}"
  resource_group_name       = "${azurerm_resource_group.example.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.example.id}"
  storage_connection_string = "${azurerm_storage_account.example.primary_connection_string}"
}
