// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package ai_foundry_project_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type AIFoundry struct{}

func (r AIFoundry) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

%s

resource "azurerm_ai_foundry" "test" {
  name                = "acctestaihub-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  storage_account_id  = azurerm_storage_account.test.id
  key_vault_id        = azurerm_key_vault.test.id

  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r AIFoundry) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    application_insights {
      disable_generated_rule = true
    }
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

%s

resource "azurerm_application_insights" "test" {
  name                = "acctestai-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  // Premium sku is required when creating a hub with AllowInternetOutbound or AllowOnlyApprovedOutbound isolation mode
  sku = "Premium"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_ai_foundry" "test" {
  name                = "acctestaihub-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  storage_account_id  = azurerm_storage_account.test.id
  key_vault_id        = azurerm_key_vault.test.id

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  application_insights_id        = azurerm_application_insights.test.id
  container_registry_id          = azurerm_container_registry.test.id
  primary_user_assigned_identity = azurerm_user_assigned_identity.test.id
  public_network_access          = "Disabled"
  description                    = "AI Hub created by Terraform"
  friendly_name                  = "AI Hub"
  high_business_impact_enabled   = true

  managed_network {
    isolation_mode = "AllowInternetOutbound"
  }

  tags = {
    env = "test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (AIFoundry) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aiservices-%[1]d"
  location = "%[2]s"
}

resource "azurerm_key_vault" "test" {
  name                       = "acctestvault%[3]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  rbac_authorization_enabled = false
  tenant_id                  = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  purge_protection_enabled   = true
  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Get",
    "Delete",
    "Purge",
    "GetRotationPolicy",
  ]
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
