# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "namespace_connection_string" {
  value = azurerm_servicebus_namespace.example.default_primary_connection_string
}

output "shared_access_policy_primarykey" {
  value = azurerm_servicebus_namespace.example.default_primary_key
}
