# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "graph_services_account_id" {
  description = "The ID of the Graph Services Account"
  value       = azurerm_graph_services_account.example.id
}

output "billing_plan_id" {
  description = "The billing plan ID of the Graph Services Account"
  value       = azurerm_graph_services_account.example.billing_plan_id
}
