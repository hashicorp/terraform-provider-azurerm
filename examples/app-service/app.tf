# Configure the Microsoft Azure Provider
provider "azurerm" {
  # if you're using a Service Principal (shared account) then either set the environment variables, or fill these in:  # subscription_id = "..."  # client_id       = "..."  # client_secret   = "..."  # tenant_id       = "..."
}

resource "azurerm_resource_group" "default" {
  name     = "${var.resource_group_name}"
  location = "${var.location}"
}

resource "random_integer" "ri" {
  min = 10000
  max = 99999
}

resource "azurerm_app_service_plan" "default" {
  name                = "tfex-appservice-${random_integer.ri.result}-plan"
  location            = "${azurerm_resource_group.default.location}"
  resource_group_name = "${azurerm_resource_group.default.name}"

  sku {
    tier = "${var.app_service_plan_sku_tier}"
    size = "${var.app_service_plan_sku_size}"
  }
}

resource "azurerm_app_service" "default" {
  name                = "tfex-appservice-${random_integer.ri.result}"
  location            = "${azurerm_resource_group.default.location}"
  resource_group_name = "${azurerm_resource_group.default.name}"
  app_service_plan_id = "${azurerm_app_service_plan.default.id}"

  site_config {
    dotnet_framework_version = "v4.0"
    remote_debugging_enabled = true
    remote_debugging_version = "VS2015"
  }

  # app_settings {
  #   "SOME_KEY" = "some-value"
  # }
  # connection_string {
  #   name  = "Database"
  #   type  = "SQLServer"
  #   value = "Server=some-server.mydomain.com;Integrated Security=SSPI"
  # }
}

output "app_service_name" {
  value = "${azurerm_app_service.default.name}"
}

output "app_service_default_hostname" {
  value = "https://${azurerm_app_service.default.default_site_hostname}"
}
