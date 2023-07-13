# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

provider "azuread" {}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-HDInsight-ESP"
  location = var.location

  tags = {
    StorageType = "Standard_LRS"
    type        = "test"
  }
}

resource "azurerm_storage_account" "example" {
  name                     = "${var.prefix}sa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "${var.prefix}-sc"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.10.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "${var.prefix}-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = [cidrsubnet(azurerm_virtual_network.example.address_space.0, 8, 0)]
}

resource "azurerm_network_security_group" "example" {
  name                = "${var.prefix}-nsg"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  security_rule {
    name                       = "AllowSyncWithAzureAD"
    priority                   = 101
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "443"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "AllowRD"
    priority                   = 201
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "3389"
    source_address_prefix      = "CorpNetSaw"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "AllowPSRemoting"
    priority                   = 301
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "5986"
    source_address_prefix      = "AzureActiveDirectoryDomainServices"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "AllowLDAPS"
    priority                   = 401
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "636"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

resource "azurerm_subnet_network_security_group_association" "example" {
  subnet_id                 = azurerm_subnet.example.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azuread_group" "example" {
  display_name = "${var.prefix} AAD DC Administrators"
  description  = "Test for delegating group to administer Azure AD Domain Services"
  mail_enabled = true
}

data "azuread_domains" "example" {
  only_initial = true
}

resource "azuread_user" "example" {
  user_principal_name = "${var.prefix}AADDSAdminUser@${data.azuread_domains.example.domains.0.domain_name}"
  display_name        = "${var.prefix}AADDSAdminUser"
  password            = "TerrAform321!"
}

resource "azuread_group_member" "example" {
  group_object_id  = azuread_group.example.object_id
  member_object_id = azuread_user.example.object_id
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "${var.prefix}-uai"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

data "azurerm_subscription" "primary" {}

resource "azurerm_role_assignment" "example" {
  scope                = data.azurerm_subscription.primary.id
  role_definition_name = "HDInsight Domain Services Contributor"
  principal_id         = azurerm_user_assigned_identity.example.principal_id
}

resource "azurerm_active_directory_domain_service" "example" {
  name                = "${var.prefix}-aadds"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  domain_name           = "never.gonna.shut.you.down"
  sku                   = "Standard"
  filtered_sync_enabled = false

  initial_replica_set {
    subnet_id = azurerm_subnet.example.id
  }

  notifications {
    additional_recipients = ["notifyA@example.net", "notifyB@example.org"]
    notify_dc_admins      = true
    notify_global_admins  = true
  }

  secure_ldap {
    enabled                  = true
    external_access_enabled  = true
    pfx_certificate          = "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"
    pfx_certificate_password = "qwer5678"
  }

  security {
    ntlm_v1_enabled         = true
    sync_kerberos_passwords = true
    sync_ntlm_passwords     = true
    sync_on_prem_passwords  = true
    tls_v1_enabled          = true
  }

  tags = {
    Environment = "test"
  }

  depends_on = [
    azuread_group_member.example,
    azurerm_role_assignment.example,
    azurerm_subnet_network_security_group_association.example,
  ]
}

resource "azurerm_virtual_network_dns_servers" "example" {
  virtual_network_id = azurerm_virtual_network.example.id
  dns_servers        = azurerm_active_directory_domain_service.example.initial_replica_set.0.domain_controller_ip_addresses
}

resource "azurerm_hdinsight_hadoop_cluster" "example" {
  name                = "${var.prefix}-hdicluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_version     = "4.0"
  tier                = "Premium"

  component_version {
    hadoop = "3.1"
  }

  gateway {
    username = "sshuser"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_id = azurerm_storage_container.example.id
    storage_account_key  = azurerm_storage_account.example.primary_access_key
    is_default           = true
  }

  roles {
    head_node {
      vm_size            = "Standard_E4_V3"
      username           = "sshuser"
      password           = "TerrAform123!"
      subnet_id          = azurerm_subnet.example.id
      virtual_network_id = azurerm_virtual_network.example.id
    }

    worker_node {
      vm_size               = "Standard_E8_V3"
      username              = "sshuser"
      password              = "TerrAform123!"
      target_instance_count = 3
      subnet_id             = azurerm_subnet.example.id
      virtual_network_id    = azurerm_virtual_network.example.id
    }

    zookeeper_node {
      vm_size            = "Standard_A2_V2"
      username           = "sshuser"
      password           = "TerrAform123!"
      subnet_id          = azurerm_subnet.example.id
      virtual_network_id = azurerm_virtual_network.example.id
    }
  }

  security_profile {
    aadds_resource_id       = azurerm_active_directory_domain_service.example.resource_id
    domain_name             = azurerm_active_directory_domain_service.example.domain_name
    domain_username         = azuread_user.example.user_principal_name
    domain_user_password    = azuread_user.example.password
    ldaps_urls              = ["ldaps://${azurerm_active_directory_domain_service.example.domain_name}:636"]
    msi_resource_id         = azurerm_user_assigned_identity.example.id
    cluster_users_group_dns = [azuread_group.example.display_name]
  }

  depends_on = [
    azurerm_virtual_network_dns_servers.example,
  ]
}
