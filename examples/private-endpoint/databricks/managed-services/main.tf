# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_databricks_workspace_private_endpoint_connection" "example" {
  workspace_id        = azurerm_databricks_workspace.example.id
  private_endpoint_id = azurerm_private_endpoint.databricks.id
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-databricks-private-endpoint-ms-dbfscmk"
  location = "eastus2"
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-vnet-databricks"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "public" {
  name                 = "${var.prefix}-sn-public"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "databricks-del-pub-${var.prefix}"

    service_delegation {
      actions = [
          "Microsoft.Network/virtualNetworks/subnets/join/action",
          "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
          "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
        ]
      name = "Microsoft.Databricks/workspaces"
    }
  }
}

resource "azurerm_subnet" "private" {
  name                 = "${var.prefix}-sn-private"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "databricks-del-pri-${var.prefix}"

    service_delegation {
      actions = [
          "Microsoft.Network/virtualNetworks/subnets/join/action",
          "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
          "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
        ]
      name = "Microsoft.Databricks/workspaces"
    }
  }
}

resource "azurerm_subnet" "endpoint" {
  name                 = "${var.prefix}-sn-private-endpoint"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.3.0/24"]

  private_endpoint_network_policies_enabled = false
}

resource "azurerm_subnet_network_security_group_association" "private" {
  subnet_id                 = azurerm_subnet.private.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_subnet_network_security_group_association" "public" {
  subnet_id                 = azurerm_subnet.public.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_network_security_group" "example" {
  name                = "${var.prefix}-nsg-databricks"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_databricks_workspace" "example" {
  name                        = "${var.prefix}-DBW"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  sku                         = "premium"
  managed_resource_group_name = "${var.prefix}-DBW-managed-private-endpoint-ms-dbfscmk"

  customer_managed_key_enabled              = true
  managed_services_cmk_key_vault_key_id     = azurerm_key_vault_key.example.id
  public_network_access_enabled             = false
  network_security_group_rules_required     = "NoAzureDatabricksRules"

  custom_parameters {
    no_public_ip        = true
    public_subnet_name  = azurerm_subnet.public.name
    private_subnet_name = azurerm_subnet.private.name
    virtual_network_id  = azurerm_virtual_network.example.id

    public_subnet_network_security_group_association_id  = azurerm_subnet_network_security_group_association.public.id
    private_subnet_network_security_group_association_id = azurerm_subnet_network_security_group_association.private.id
  }

  tags = {
    Environment = "Production"
    Pricing     = "Premium"
  }
}

resource "azurerm_databricks_workspace_root_dbfs_customer_managed_key" "example" {
  depends_on = [azurerm_key_vault_access_policy.databricks]

  workspace_id     = azurerm_databricks_workspace.example.id
  key_vault_key_id = azurerm_key_vault_key.example.id
}

resource "azurerm_private_endpoint" "databricks" {
  depends_on = [azurerm_databricks_workspace_root_dbfs_customer_managed_key.example]

  name                = "${var.prefix}-pe-databricks"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = "${var.prefix}-psc"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_databricks_workspace.example.id
    subresource_names              = ["databricks_ui_api"]
  }
}

resource "azurerm_private_dns_zone" "example" {
  depends_on = [azurerm_private_endpoint.databricks]

  name                = "privatelink.azuredatabricks.net"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_private_dns_cname_record" "example" {
  name                = azurerm_databricks_workspace.example.workspace_url
  zone_name           = azurerm_private_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300
  record              = "eastus2-c2.azuredatabricks.net"
}

resource "azurerm_key_vault" "example" {
  name                = "${var.prefix}-keyvault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_key" "example" {
  depends_on = [azurerm_key_vault_access_policy.terraform]

  name         = "${var.prefix}-certificate"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_key_vault_access_policy" "terraform" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = azurerm_key_vault.example.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "get",
    "list",
    "create",
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
    "delete",
    "restore",
    "recover",
    "update",
    "purge",
  ]
}

resource "azurerm_key_vault_access_policy" "databricks" {
  depends_on = [azurerm_databricks_workspace.example]

  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = azurerm_databricks_workspace.example.storage_account_identity.0.tenant_id
  object_id    = azurerm_databricks_workspace.example.storage_account_identity.0.principal_id

  key_permissions = [
    "get",
    "unwrapKey",
    "wrapKey",
  ]
}

resource "azurerm_key_vault_access_policy" "managed" {
  key_vault_id   = azurerm_key_vault.example.id
  tenant_id      = azurerm_key_vault.example.tenant_id
  object_id      = "See the README.md file for instructions on how to lookup the correct value to enter here"

  key_permissions = [
    "get",
    "unwrapKey",
    "wrapKey",
  ]
}

output "databricks_workspace_private_endpoint_connection_workspace_id" {
  value = data.azurerm_databricks_workspace_private_endpoint_connection.example.workspace_id
}

output "databricks_workspace_private_endpoint_connection_private_endpoint_id" {
  value = data.azurerm_databricks_workspace_private_endpoint_connection.example.private_endpoint_id
}

output "databricks_workspace_private_endpoint_connection_name" {
  value = data.azurerm_databricks_workspace_private_endpoint_connection.example.connections.0.name
}

output "databricks_workspace_private_endpoint_connection_workspace_private_endpoint_id" {
  value = data.azurerm_databricks_workspace_private_endpoint_connection.example.connections.0.workspace_private_endpoint_id
}

output "databricks_workspace_private_endpoint_connection_status" {
  value = data.azurerm_databricks_workspace_private_endpoint_connection.example.connections.0.status
}

output "databricks_workspace_private_endpoint_connection_description" {
  value = data.azurerm_databricks_workspace_private_endpoint_connection.example.connections.0.description
}

output "databricks_workspace_private_endpoint_connection_action_required" {
  value = data.azurerm_databricks_workspace_private_endpoint_connection.example.connections.0.action_required
}