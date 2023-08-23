# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "fqdn" {
  value = azurerm_traffic_manager_profile.example.fqdn
}
