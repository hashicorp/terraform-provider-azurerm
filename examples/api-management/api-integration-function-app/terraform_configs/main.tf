data "azurerm_resource_group" "rg" {
  name = var.resource_group_name
}

resource "random_string" "random_name" {
  length  = 13
  lower   = true
  numeric = false
  special = false
  upper   = false
}

resource "azurerm_api_management" "example" {
  name                = "myapi${random_string.random_name.result}"
  location            = data.azurerm_resource_group.rg.location
  resource_group_name = data.azurerm_resource_group.rg.name
  publisher_email     = "test@contoso.com"
  publisher_name      = "example publisher"
  sku_name            = "Developer_1"

}

resource "azurerm_api_management_api" "example" {
  name                  = "example-api"
  resource_group_name   = data.azurerm_resource_group.rg.name
  api_management_name   = azurerm_api_management.example.name
  revision              = "1"
  display_name          = "Example API"
  api_type              = "http"
  protocols             = ["https"]
  subscription_required = false
}

resource "azurerm_api_management_api_operation" "example" {
  operation_id        = "op1"
  api_name            = azurerm_api_management_api.example.name
  api_management_name = azurerm_api_management.example.name
  resource_group_name = data.azurerm_resource_group.rg.name
  display_name        = "GET Resource"
  method              = "GET"
  url_template        = "funcfromcli" # Ensure this value matches the Azure Function's name exactly as it appears in the invoke URL.
  response {
    status_code = 200
    description = "Successful GET request"
  }
}

resource "azurerm_api_management_api_policy" "example" {
  api_name            = azurerm_api_management_api.example.name
  resource_group_name = data.azurerm_resource_group.rg.name
  api_management_name = azurerm_api_management.example.name
    xml_content = templatefile("policy.xml", {
    base-url        = "https://${azurerm_linux_function_app.example.default_hostname}/api"
  })
}


resource "azurerm_api_management_backend" "example" {
  name                = "example-backend"
  resource_group_name = data.azurerm_resource_group.rg.name
  api_management_name = azurerm_api_management.example.name
  protocol            = "http"
  url                 = "https://${azurerm_linux_function_app.example.default_hostname}/api" #make sure this ends exactly like this. just "api" without slash
}

resource "azurerm_storage_account" "example" {
  name                     = random_string.random_name.result
  resource_group_name      = data.azurerm_resource_group.rg.name
  location                 = data.azurerm_resource_group.rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

}

resource "azurerm_service_plan" "example" {
  name                = "example-app-service-plan"
  resource_group_name = data.azurerm_resource_group.rg.name
  location            = data.azurerm_resource_group.rg.location
  os_type             = "Linux"
  sku_name            = "B1"

}

resource "azurerm_linux_function_app" "example" {
  name                = "myfuncapp${random_string.random_name.result}"
  resource_group_name = data.azurerm_resource_group.rg.name
  location            = data.azurerm_resource_group.rg.location
  service_plan_id     = azurerm_service_plan.example.id

  storage_account_name       = azurerm_storage_account.example.name
  storage_account_access_key = azurerm_storage_account.example.primary_access_key
  site_config {

    application_stack {
      node_version = "20"

    }
  }
}