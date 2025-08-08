# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "prefix" {
  description = "The prefix used for all resources in this example"
  type        = string
  default     = "anf-ct-cmk"
}

variable "location" {
  description = "The Azure region where all resources in this example should be created"
  type        = string
  default     = "West Europe"
}

# Cross-tenant key vault details (existing resources in remote tenant)
variable "remote_subscription_id" {
  description = "The subscription ID of the remote tenant where the key vault exists"
  type        = string
}

variable "cross_tenant_key_vault_name" {
  description = "The name of the existing key vault in the remote tenant"
  type        = string
}

variable "cross_tenant_resource_group_name" {
  description = "The resource group name where the cross-tenant key vault exists"
  type        = string
}

variable "cross_tenant_key_name" {
  description = "The name of the existing encryption key in the cross-tenant key vault"
  type        = string
}

# Federated application details
variable "federated_client_id" {
  description = "The client ID of the multi-tenant Entra ID application used for cross-tenant access"
  type        = string

  validation {
    condition     = length(var.federated_client_id) > 0
    error_message = "The federated_client_id cannot be empty."
  }
}

# Private endpoint configuration
variable "private_endpoint_manual_approval" {
  description = "Whether the private endpoint connection requires manual approval"
  type        = bool
  default     = true
}

# Existing managed identity (in your tenant) to access the cross-tenant key vault
variable "user_assigned_identity_id" {
  description = "The full resource ID of the pre-created user-assigned managed identity used by the NetApp account (e.g., /subscriptions/<sub>/resourceGroups/<rg>/providers/Microsoft.ManagedIdentity/userAssignedIdentities/<name>)"
  type        = string
}

# Cross-tenant key vault resource ID (optional but mandatory for cross-tenant scenarios)
variable "cross_tenant_key_vault_resource_id" {
  description = "The full resource ID of the cross-tenant key vault (e.g., /subscriptions/<remote-sub>/resourceGroups/<remote-rg>/providers/Microsoft.KeyVault/vaults/<vault-name>)"
  type        = string
  default     = ""
}
