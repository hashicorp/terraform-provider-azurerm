resource "azurerm_app_service_source_control_token" "example" {
  type  = "GitHub"
  token = "${var.github_token}"
}
