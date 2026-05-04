# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {
    netapp {
      prevent_volume_destruction = true
    }
  }
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location

  tags = {
    "SkipNRMSNSG" = "true"
  }
}

resource "azurerm_netapp_account" "example" {
  name                = "${var.prefix}-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  identity {
    type = "SystemAssigned"
  }

  tags = {
    "CreatedOnDate" = "2022-07-08T23:50:21Z"
  }
}

resource "azurerm_key_vault" "example" {
  name                            = "${var.prefix}anfakv"
  location                        = azurerm_resource_group.example.location
  resource_group_name             = azurerm_resource_group.example.name
  enabled_for_disk_encryption     = true
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  purge_protection_enabled        = true
  tenant_id                       = var.tenant_id
  sku_name                        = "standard"

  access_policy {
    tenant_id = azurerm_netapp_account.example.identity.0.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = []
    secret_permissions      = []
    storage_permissions     = []
    key_permissions = [
      "Get",
      "Create",
      "Delete",
      "WrapKey",
      "UnwrapKey",
      "GetRotationPolicy",
      "SetRotationPolicy",
    ]
  }

  access_policy {
    tenant_id = azurerm_netapp_account.example.identity.0.tenant_id
    object_id = azurerm_netapp_account.example.identity.0.principal_id

    certificate_permissions = []
    secret_permissions      = []
    storage_permissions     = []
    key_permissions = [
      "Get",
      "Encrypt",
      "Decrypt"
    ]
  }

  tags = {
    "CreatedOnDate" = "2022-07-08T23:50:21Z"
  }
}

resource "azurerm_key_vault_key" "example" {
  name         = "${var.prefix}anfenckey"
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

resource "azurerm_netapp_account_encryption" "example" {
  netapp_account_id                     = azurerm_netapp_account.example.id
  system_assigned_identity_principal_id = azurerm_netapp_account.example.identity.0.principal_id
  encryption_key                        = azurerm_key_vault_key.example.versionless_id
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-virtual-network"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.6.0.0/16"]

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
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

resource "azurerm_private_endpoint" "example" {
  name                = "${var.prefix}-pe-akv"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  subnet_id           = azurerm_subnet.example-non-delegated.id

  private_service_connection {
    name                           = "${var.prefix}-pe-sc-akv"
    private_connection_resource_id = azurerm_key_vault.example.id
    is_manual_connection           = false
    subresource_names              = ["Vault"]
  }

  tags = {
    CreatedOnDate = "2023-10-03T19:58:43.6509795Z"
  }
}

resource "azurerm_netapp_pool" "example" {
  name                = "${var.prefix}-capacity-pool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Standard"
  size_in_tb          = 4
  qos_type            = "Manual"

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }

  depends_on = [
    azurerm_netapp_account_encryption.example
  ]
}

resource "azurerm_netapp_volume_group_oracle" "example" {
  name                   = "${var.prefix}-volume-group-oracle"
  location               = azurerm_resource_group.example.location
  resource_group_name    = azurerm_resource_group.example.name
  account_name           = azurerm_netapp_account.example.name
  group_description      = "Example volume group for Oracle"
  application_identifier = "TST"

  volume {
    name                          = "${var.prefix}-volume-ora1"
    volume_path                   = "${var.prefix}-my-unique-file-ora-path-1"
    service_level                 = "Standard"
    capacity_pool_id              = azurerm_netapp_pool.example.id
    subnet_id                     = azurerm_subnet.example-delegated.id
    zone                          = "1"
    volume_spec_name              = "ora-data1"
    storage_quota_in_gb           = 1024
    throughput_in_mibps           = 24
    protocols                     = ["NFSv4.1"]
    security_style                = "unix"
    snapshot_directory_visible    = false
    encryption_key_source         = "Microsoft.KeyVault"
    key_vault_private_endpoint_id = azurerm_private_endpoint.example.id
    network_features              = "Standard"

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }
  }

  volume {
    name                          = "${var.prefix}-volume-oraLog"
    volume_path                   = "${var.prefix}-my-unique-file-oralog-path"
    service_level                 = "Standard"
    capacity_pool_id              = azurerm_netapp_pool.example.id
    subnet_id                     = azurerm_subnet.example-delegated.id
    zone                          = "1"
    volume_spec_name              = "ora-log"
    storage_quota_in_gb           = 1024
    throughput_in_mibps           = 24
    protocols                     = ["NFSv4.1"]
    security_style                = "unix"
    snapshot_directory_visible    = false
    encryption_key_source         = "Microsoft.KeyVault"
    key_vault_private_endpoint_id = azurerm_private_endpoint.example.id
    network_features              = "Standard"

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }
  }
}
