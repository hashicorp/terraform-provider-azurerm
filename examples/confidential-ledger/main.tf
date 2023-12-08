# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "eastus"
}

data "azurerm_client_config" "current" {
}

resource "azurerm_confidential_ledger" "example" {
  name                = "example-ledger"
  ledger_type         = "Public"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  aad_based_security_principals {
    principal_id     = data.azurerm_client_config.current.object_id
    tenant_id        = data.azurerm_client_config.current.tenant_id
    ledger_role_name = "Administrator"
  }

  tags = {
    IsExample = "True"
  }
}