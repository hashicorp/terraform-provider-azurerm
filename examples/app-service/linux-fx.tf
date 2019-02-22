resource "azurerm_app_service" "test" {
  # ...
  site_config = {
    # ...
    linux_fx_version = "COMPOSE|<base64encoded_docker_compose_file>"
  }

  lifecycle {
    ignore_changes = [
      "site_config.0.linux_fx_version", # will be updated with every deployment
    ]
  }
}