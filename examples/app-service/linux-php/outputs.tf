# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

locals {
  scm_username = azurerm_linux_web_app.example.site_credential.0.name
  scm_password = azurerm_linux_web_app.example.site_credential.0.password
  repo_uri     = replace(azurerm_app_service_source_control.example.repo_url, "https://", "")
}

output "repository_url" {
  value = "https://${local.scm_username}:${local.scm_password}@${local.repo_uri}/${azurerm_linux_web_app.example.name}.git"
}

output "app_name" {
  value = azurerm_linux_web_app.example.default_hostname
}
