# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_key_vault" "example" {
  name                = "${var.prefix}-key-vault"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"
}

resource "azurerm_key_vault_access_policy" "current_user" {
  key_vault_id = azurerm_key_vault.example.id

  tenant_id = azurerm_key_vault.example.tenant_id
  object_id = data.azurerm_client_config.current.object_id

  certificate_permissions = [
    "Get",
    "Import"
  ]
}

data "azuread_service_principal" "web_app_resource_provider" {
  client_id = data.azurerm_client_config.current.client_id
}

resource "azurerm_key_vault_access_policy" "web_app_resource_provider" {
  key_vault_id = azurerm_key_vault.example.id

  tenant_id = azurerm_key_vault.example.tenant_id
  object_id = data.azuread_service_principal.web_app_resource_provider.id

  secret_permissions = [
    "Get"
  ]

  certificate_permissions = [
    "Get"
  ]
}

resource "azurerm_key_vault_certificate" "example" {
  name         = "${var.prefix}-cert"
  key_vault_id = azurerm_key_vault.example.id

  certificate {
    contents = filebase64("certificate.pfx")
    password = "terraform"
  }

  certificate_policy {
    issuer_parameters {
      name = "Unknown"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }
  }

  depends_on = [
    azurerm_key_vault_access_policy.current_user
  ]
}

resource "azurerm_app_service_certificate" "example" {
  name                = "${var.prefix}-cert"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  key_vault_secret_id = azurerm_key_vault_certificate.example.secret_id
}
