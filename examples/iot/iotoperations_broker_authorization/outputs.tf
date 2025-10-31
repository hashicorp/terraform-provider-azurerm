# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "iotoperations_broker_authorization_id" {
  description = "The ARM resource ID of the IoT Operations broker authorization"
  value       = azurerm_iotoperations_broker_authorization.example.id
}