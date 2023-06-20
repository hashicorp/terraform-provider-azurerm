provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-functions-python-rg"
  location = var.location
}

resource "azurerm_storage_account" "example" {
  name                     = "${var.prefix}storageacct"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "example" {
  name                = "${var.prefix}-sp"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  os_type             = "Linux"
  sku_name            = "S1"
}

resource "azurerm_linux_function_app" "example" {
  name                = "${var.prefix}-python-example-app"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  service_plan_id     = azurerm_service_plan.example.id

  storage_account_name       = azurerm_storage_account.example.name
  storage_account_access_key = azurerm_storage_account.example.primary_access_key

  site_config {
    application_stack {
      python_version = "3.9"
    }
  }
}

resource "azurerm_function_app_function" "example" {
  name            = "example-python-function"
  function_app_id = azurerm_linux_function_app.example.id
  language        = "Python"
  file {
    name    = "__init__.py"
    content = file("./SampleApp/PythonSampleApp/__init__.py")
  }
  test_data = file("./SampleApp/PythonSampleApp/sample.dat")
  #  test_data = jsonencode({
  #    "name" = "Azure"
  #  })
  config_json = file("./SampleApp/PythonSampleApp/function.json")
  #  config_json = jsonencode({
  #    "scriptFile" = "__init__.py"
  #    "bindings" = [
  #      {
  #        "authLevel" = "anonymous"
  #        "direction" = "in"
  #        "methods" = [
  #          "get",
  #          "post",
  #        ]
  #        "name" = "req"
  #        "type" = "httpTrigger"
  #      },
  #      {
  #        "direction" = "out"
  #        "name"      = "$return"
  #        "type"      = "http"
  #      },
  #    ]
  #  })
}