provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_service_plan" "example" {
  name                = "${var.prefix}-sp-zipdeploy"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  os_type             = "Linux"
  sku_name            = "S1"
}

resource "azurerm_linux_web_app" "example" {
  name                = "${var.prefix}-example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  service_plan_id     = azurerm_service_plan.example.id

  app_settings = {
    "WEBSITES_ENABLE_APP_SERVICE_STORAGE" = "false"
  }

  site_config {
    application_stack {
      docker_image     = "jackofallops/azure-containerapps-python-acctest"
      docker_image_tag = "v0.0.1"
    }
  }

  auth_settings {
    enabled                       = true
    issuer                        = "https://sts.windows.net/${data.azurerm_client_config.current.tenant_id}"
    default_provider              = "AzureActiveDirectory"
    token_refresh_extension_hours = 24
    token_store_enabled           = true
    unauthenticated_client_action = "RedirectToLoginPage"

    additional_login_parameters = {
      test_key = "test_value"
    }

    allowed_external_redirect_urls = ["https://example.com"]

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }

    microsoft {
      client_id     = "microsoftclientid"
      client_secret = "microsoftclientsecret"

      oauth_scopes = [
        "microsoftscope",
      ]
    }
  }
}