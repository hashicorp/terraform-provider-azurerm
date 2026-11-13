// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type NetAppAccountEncryptionResource struct{}

func (r NetAppAccountEncryptionResource) cmkSystemAssigned(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_client_config" "current" {
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }

  tags = {
    "CreatedOnDate" = "2022-07-08T23-50-21Z"
  }
}

resource "azurerm_key_vault" "test" {
  name                            = "anfakv%[2]d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  rbac_authorization_enabled      = false
  enabled_for_disk_encryption     = true
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  purge_protection_enabled        = true
  tenant_id                       = "%[3]s"
  sku_name                        = "standard"

  access_policy {
    tenant_id = azurerm_netapp_account.test.identity.0.tenant_id
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
    tenant_id = azurerm_netapp_account.test.identity.0.tenant_id
    object_id = azurerm_netapp_account.test.identity.0.principal_id

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
    "CreatedOnDate" = "2022-07-08T23-50-21Z"
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "anfenckey%[2]d"
  key_vault_id = azurerm_key_vault.test.id
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

resource "azurerm_netapp_account_encryption" "test" {
  netapp_account_id                     = azurerm_netapp_account.test.id
  system_assigned_identity_principal_id = azurerm_netapp_account.test.identity.0.principal_id
  encryption_key                        = azurerm_key_vault_key.test.versionless_id
}
`, r.template(data), data.RandomInteger, tenantID)
}

func (r NetAppAccountEncryptionResource) cmkUserAssigned(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "user-assigned-identity-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "CreatedOnDate" = "2022-07-08T23-50-21Z"
  }
}

data "azurerm_client_config" "current" {
}

resource "azurerm_key_vault" "test" {
  name                            = "anfkv%[4]d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  rbac_authorization_enabled      = false
  enabled_for_disk_encryption     = true
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  purge_protection_enabled        = true
  soft_delete_retention_days      = 7
  tenant_id                       = "%[3]s"
  sku_name                        = "standard"

  access_policy {
    tenant_id = "%[3]s"
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
    tenant_id = "%[3]s"
    object_id = azurerm_user_assigned_identity.test.principal_id

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
    "CreatedOnDate" = "2022-07-08T23-50-21Z"
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "anfenckey%[2]d"
  key_vault_id = azurerm_key_vault.test.id
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

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_account_encryption" "test" {
  netapp_account_id         = azurerm_netapp_account.test.id
  user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  encryption_key            = azurerm_key_vault_key.test.versionless_id
}
`, r.template(data), data.RandomInteger, tenantID, data.RandomIntOfLength(17))
}

func (NetAppAccountEncryptionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }

    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }

    netapp {
      prevent_volume_destruction             = false
      delete_backups_on_backup_vault_destroy = true
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%[1]d"
  location = "%[2]s"

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true",
    "SkipNRMSNSG"      = "true"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
