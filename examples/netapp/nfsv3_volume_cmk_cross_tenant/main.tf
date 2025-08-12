# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Example of NetApp Account Encryption with Cross-Tenant Customer-Managed Keys
# This example demonstrates how to configure NetApp encryption using a key vault
# that already exists in a different Entra ID tenant (cross-tenant scenario)
#
# IMPORTANT: In a real cross-tenant scenario, you cannot create resources in the remote tenant.
# The key vault and keys must already exist and be managed by the remote tenant administrators.
# This example references existing cross-tenant resources by ID/URL and requires a pre-created
# user-assigned identity ID to be provided via variables.

terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~>4.0"
    }
    time = {
      source  = "hashicorp/time"
      version = "~>0.9"
    }
  }
}

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.6.0.0/16"]
}

resource "azurerm_subnet" "example-delegated" {
  name                 = "${var.prefix}-delegated-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.6.1.0/24"]

  delegation {
    name = "exampledelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_subnet" "example-non-delegated" {
  name                 = "${var.prefix}-non-delegated-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.6.0.0/24"]
}

# User-assigned identity that can access cross-tenant key vault (pre-existing)
# Provided via variable: var.user_assigned_identity_id

# Private DNS zone for Key Vault (required for private endpoint resolution)
resource "azurerm_private_dns_zone" "keyvault" {
  name                = "privatelink.vaultcore.azure.net"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "keyvault" {
  name                  = "${var.prefix}-dns-link"
  resource_group_name   = azurerm_resource_group.example.name
  private_dns_zone_name = azurerm_private_dns_zone.keyvault.name
  virtual_network_id    = azurerm_virtual_network.example.id
  registration_enabled  = false
}

# Private endpoint to the cross-tenant key vault
# In cross-tenant scenarios, this often requires manual approval
resource "azurerm_private_endpoint" "cross_tenant" {
  name                = "${var.prefix}-pe-ct-akv"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  subnet_id           = azurerm_subnet.example-non-delegated.id

  private_service_connection {
    name = "${var.prefix}-pe-sc-ct-akv"
    # Using the resource ID format for cross-tenant key vault
    private_connection_resource_id = "/subscriptions/${var.remote_subscription_id}/resourceGroups/${var.cross_tenant_resource_group_name}/providers/Microsoft.KeyVault/vaults/${var.cross_tenant_key_vault_name}"
    is_manual_connection           = var.private_endpoint_manual_approval
    subresource_names              = ["Vault"]
    request_message                = var.private_endpoint_manual_approval ? "Cross-tenant access for NetApp encryption" : null
  }

  private_dns_zone_group {
    name                 = "${var.prefix}-dns-zone-group"
    private_dns_zone_ids = [azurerm_private_dns_zone.keyvault.id]
  }

  depends_on = [azurerm_private_dns_zone_virtual_network_link.keyvault]
}

# Time-based wait for private endpoint approval
resource "time_sleep" "private_endpoint_approval_time_wait" {
  count = var.private_endpoint_manual_approval ? 1 : 0

  create_duration = "${var.private_endpoint_approval_wait_time}m"

  triggers = {
    private_endpoint_id = azurerm_private_endpoint.cross_tenant.id
  }

  depends_on = [azurerm_private_endpoint.cross_tenant]
}

resource "azurerm_netapp_account" "example" {
  name                = "${var.prefix}-netappaccount"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  identity {
    type = "UserAssigned"
    identity_ids = [
      var.user_assigned_identity_id
    ]
  }
}

# NetApp Account Encryption with Cross-Tenant CMK
# This will use the actual remote key vault URL
resource "azurerm_netapp_account_encryption" "example" {
  netapp_account_id = azurerm_netapp_account.example.id

  user_assigned_identity_id = var.user_assigned_identity_id

  encryption_key = "https://${var.cross_tenant_key_vault_name}.vault.azure.net/keys/${var.cross_tenant_key_name}"

  # This is the Client ID of the multi-tenant Entra ID application
  # that has been granted access to the cross-tenant key vault
  federated_client_id = var.federated_client_id

  # The full resource ID of the cross-tenant key vault (mandatory for cross-tenant scenarios)
  cross_tenant_key_vault_resource_id = var.cross_tenant_key_vault_resource_id

  depends_on = [
   azurerm_netapp_account.example,
   time_sleep.private_endpoint_approval_time_wait
  ]
}

resource "azurerm_netapp_pool" "example" {
  name                = "${var.prefix}-pool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Standard"
  size_in_tb          = 4

  depends_on = [
    azurerm_netapp_account_encryption.example
  ]
}

resource "azurerm_netapp_volume" "example" {
  name                          = "${var.prefix}-vol"
  location                      = azurerm_resource_group.example.location
  resource_group_name           = azurerm_resource_group.example.name
  account_name                  = azurerm_netapp_account.example.name
  pool_name                     = azurerm_netapp_pool.example.name
  volume_path                   = "${var.prefix}-my-unique-file-path-vol"
  service_level                 = "Standard"
  subnet_id                     = azurerm_subnet.example-delegated.id
  storage_quota_in_gb           = 100
  network_features              = "Standard"
  encryption_key_source         = "Microsoft.KeyVault"
  key_vault_private_endpoint_id = azurerm_private_endpoint.cross_tenant.id

  export_policy_rule {
    rule_index          = 1
    allowed_clients     = ["0.0.0.0/0"]
    protocols_enabled   = ["NFSv3"]
    unix_read_only      = false
    unix_read_write     = true
    root_access_enabled = true
  }

  depends_on = [
    azurerm_netapp_account_encryption.example,
    azurerm_private_endpoint.cross_tenant
  ]
}
