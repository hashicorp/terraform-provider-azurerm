# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "csr" {
  value = azurerm_app_service_certificate_order.test.csr
}
