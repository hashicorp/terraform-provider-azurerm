# Configure the Microsoft Azure Provider
provider "azurerm" {
  # if you're using a Service Principal (shared account) then either set the environment variables, or fill these in:
  # subscription_id = "..."
  # client_id       = "..."
  # client_secret   = "..."
  # tenant_id       = "..."
}

resource "azurerm_resource_group" "g1" {
  name     = "${var.resource_group_name}"
  location = "${var.location}"
}

resource "azurerm_app_service_plan" "plan0" {
  name                = "service-plan0"
  location            = "${var.location}"
  resource_group_name = "${azurerm_resource_group.g1.name}"

  sku {
    tier = "Basic" # Basic | Standard | ...
    size = "B1"    # B1 | S1 | ...
  }
}

# underscores not supported as app_service name -> if not: you will receive error 400

resource "azurerm_app_service" "common_service" {
  name                = "${var.app_service_name}"
  location            = "${var.location}"
  resource_group_name = "${azurerm_resource_group.g1.name}"
  app_service_plan_id = "${azurerm_app_service_plan.plan0.id}"

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

  provisioner "local-exec" {
    command = "curl -k -u ${var.deploy_user}:${var.deploy_pass} -X PUT --data-binary @${var.deployZipFile} https://${azurerm_app_service.common_service.name}.scm.azurewebsites.net/api/zip/site/wwwroot/"

    # interpreter = ["cmd"]
  }
}

output "service" {
  value = "${azurerm_app_service.common_service.name}"
}

output "service_url" {
  value = "https://${azurerm_app_service.common_service.name}.azurewebsites.net"
}

output "admin_url" {
  value = "https://${azurerm_app_service.common_service.name}.scm.azurewebsites.net"
}
