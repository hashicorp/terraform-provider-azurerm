# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "public_ip_address" {
  value = "${azurerm_public_ip.static.ip_address}"
}
