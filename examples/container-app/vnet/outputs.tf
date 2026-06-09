# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

output "container_app_url" {
  value       = "http://${azurerm_container_app.example.ingress[0].fqdn}"
  description = "This is the url that can be used to access the app from other resources in the same virtual network (e.g. APIM)"
}
