# Copyright IBM Corp. 2023, 2025
# SPDX-License-Identifier: MPL-2.0

output "fqdn" {
  value = azurerm_traffic_manager_profile.example.fqdn
}
