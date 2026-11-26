# Copyright IBM Corp. 2014, 2026
# SPDX-License-Identifier: MPL-2.0

output "azure_search_service" {
  value = azurerm_search_service.example.name
}
