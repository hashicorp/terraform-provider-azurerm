# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "iotoperations_broker_authentication_id" {
  description = "The ARM resource ID of the IoT Operations broker authentication"
  value       = azurerm_iotoperations_broker_authentication.example.id
}