provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "main" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_storage_account" "main" {
  name                     = "${var.prefix}storageacct"
  resource_group_name      = azurerm_resource_group.main.name
  location                 = azurerm_resource_group.main.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_application_insights" "main" {
  name                = "${var.prefix}-appinsights"
  resource_group_name = azurerm_resource_group.main.name
  location            = azurerm_resource_group.main.location
  application_type    = "web"
}

resource "azurerm_app_service_plan" "main" {
  name                = "${var.prefix}-asp"
  resource_group_name = azurerm_resource_group.main.name
  location            = azurerm_resource_group.main.location
  kind                = "FunctionApp"

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}

resource "azurerm_function_app" "main" {
  name                      = "${var.prefix}-function"
  resource_group_name       = azurerm_resource_group.main.name
  location                  = azurerm_resource_group.main.location
  app_service_plan_id       = azurerm_app_service_plan.main.id
  storage_connection_string = azurerm_storage_account.main.primary_connection_string

  app_settings = {
    AppInsights_InstrumentationKey = azurerm_application_insights.main.instrumentation_key
  }
  
  identity {
    type = "SystemAssigned"   
  }
}

data "azurerm_subscription" "primary" {
}

resource "azurerm_role_assignment" "main" {
  scope                = data.azurerm_subscription.primary.id
  role_definition_name = var.role_definition_name
  principal_id         = azurerm_function_app.main.identity.0.principal_id
}