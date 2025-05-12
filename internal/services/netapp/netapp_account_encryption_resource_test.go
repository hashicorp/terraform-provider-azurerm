// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/netappaccounts"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetAppAccountEncryptionResource struct{}

func TestAccNetAppAccountEncryption_cmkSystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account_encryption", "test")
	r := NetAppAccountEncryptionResource{}

	tenantID := os.Getenv("ARM_TENANT_ID")

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cmkSystemAssigned(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption_key").IsSet(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppAccountEncryption_cmkUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account_encryption", "test")
	r := NetAppAccountEncryptionResource{}

	tenantID := os.Getenv("ARM_TENANT_ID")

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cmkUserAssigned(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption_key").IsSet(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppAccountEncryption_updateKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account_encryption", "test")
	r := NetAppAccountEncryptionResource{}

	tenantID := os.Getenv("ARM_TENANT_ID")

	regexInitialKey := regexp.MustCompile(`anfenckey\d+$`)
	regexNewKey := regexp.MustCompile(`.*anfenckey-new.*`)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyUpdate1(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption_key").MatchesRegex(regexInitialKey),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyUpdate2(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption_key").MatchesRegex(regexNewKey),
			),
		},
		data.ImportStep(),
	})
}

func (t NetAppAccountEncryptionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := netappaccounts.ParseNetAppAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NetApp.AccountClient.AccountsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Netapp Account (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

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
    "CreatedOnDate" = "2022-07-08T23:50:21Z"
  }
}

resource "azurerm_key_vault" "test" {
  name                            = "anfakv%[2]d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
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
    "CreatedOnDate" = "2022-07-08T23:50:21Z"
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
    "CreatedOnDate" = "2022-07-08T23:50:21Z"
  }
}

data "azurerm_client_config" "current" {
}

resource "azurerm_key_vault" "test" {
  name                            = "anfakv%[2]d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  enabled_for_disk_encryption     = true
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  purge_protection_enabled        = true
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
    "CreatedOnDate" = "2022-07-08T23:50:21Z"
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
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_account_encryption" "test" {
  netapp_account_id         = azurerm_netapp_account.test.id
  user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  encryption_key            = azurerm_key_vault_key.test.versionless_id
}
`, r.template(data), data.RandomInteger, tenantID)
}

func (r NetAppAccountEncryptionResource) keyUpdate1(data acceptance.TestData, tenantID string) string {
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
    "SkipNRMSNSG"   = "true",
    "CreatedOnDate" = "2022-07-08T23:50:21Z"
  }
}

resource "azurerm_key_vault" "test" {
  name                            = "anfakv%[2]d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
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
    "CreatedOnDate" = "2022-07-08T23:50:21Z"
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

resource "azurerm_key_vault_key" "test-new-key" {
  name         = "anfenckey-new%[2]d"
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

func (r NetAppAccountEncryptionResource) keyUpdate2(data acceptance.TestData, tenantID string) string {
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
    "SkipNRMSNSG"   = "true",
    "CreatedOnDate" = "2022-07-08T23:50:21Z"
  }
}

resource "azurerm_key_vault" "test" {
  name                            = "anfakv%[2]d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
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
    "CreatedOnDate" = "2022-07-08T23:50:21Z"
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

resource "azurerm_key_vault_key" "test-new-key" {
  name         = "anfenckey-new%[2]d"
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

  depends_on = [
    azurerm_key_vault_key.test
  ]
}

resource "azurerm_netapp_account_encryption" "test" {
  netapp_account_id                     = azurerm_netapp_account.test.id
  system_assigned_identity_principal_id = azurerm_netapp_account.test.identity.0.principal_id
  encryption_key                        = azurerm_key_vault_key.test-new-key.versionless_id
}
`, r.template(data), data.RandomInteger, tenantID)
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
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true",
    "SkipNRMSNSG"      = "true"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
