# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "iotoperations_broker_id" {
  description = "The ARM resource ID of the IoT Operations broker"
  value       = azurerm_iotoperations_broker.example.id
}