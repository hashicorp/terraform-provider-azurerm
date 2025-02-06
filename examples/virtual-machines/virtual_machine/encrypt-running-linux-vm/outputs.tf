# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "hostname" {
  value = var.hostname
}

output "BitLockerKey" {
  value     = azurerm_template_deployment.linux_vm.outputs["BitLockerKey"]
  sensitive = true
}
