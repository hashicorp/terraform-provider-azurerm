# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "azure_search_service" {
  value = azurerm_search_service.example.name
}
