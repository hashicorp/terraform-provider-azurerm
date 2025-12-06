# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

output "cluster_egress_ip" {
  value = data.azurerm_public_ip.example.ip_address
}
