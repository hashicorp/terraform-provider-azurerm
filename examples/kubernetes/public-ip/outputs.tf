# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "cluster_egress_ip" {
  value = data.azurerm_public_ip.example.ip_address
}
