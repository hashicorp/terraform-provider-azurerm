# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "iotoperations_instance_id" {
  description = "The ARM resource ID of the IoT Operations instance"
  value       = azurerm_iotoperations_instance.example.id
}