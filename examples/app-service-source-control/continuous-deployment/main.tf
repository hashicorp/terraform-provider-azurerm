resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

# NOTE: Source Control Tokens are configured at the subscription level, 
# not on each App Service - as such this can only be configured Subscription-wide.
resource "azurerm_app_service_source_control_token" "example" {
  type  = "GitHub"
  token = "${var.github_token}"
}

resource "azurerm_app_service_plan" "example" {
  name                = "${var.prefix}-asp"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "example" {
  name                = "${var.prefix}-appservice"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  app_service_plan_id = "${azurerm_app_service_plan.example.id}"

  lifecycle {
    ignore_changes = [site_config.0.scm_type]
  }
}

resource "azurerm_app_service_source_control" "example" {
  app_service_id = "${azurerm_app_service.example.id}"
  repo_url       = "${var.repo_url}"
  branch         = "master"
}
