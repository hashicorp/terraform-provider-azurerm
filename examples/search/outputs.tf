# Copyright IBM Corp. 2023, 2025
# SPDX-License-Identifier: MPL-2.0

output "azure_search_service" {
  value = azurerm_search_service.example.name
}
