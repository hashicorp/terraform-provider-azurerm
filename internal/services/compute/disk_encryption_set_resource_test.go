// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DiskEncryptionSetResource struct{}

func TestAccDiskEncryptionSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")
	r := DiskEncryptionSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption_type").HasValue("EncryptionAtRestWithCustomerKey"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDiskEncryptionSet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")
	r := DiskEncryptionSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDiskEncryptionSet_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")
	r := DiskEncryptionSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDiskEncryptionSet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")
	r := DiskEncryptionSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDiskEncryptionSet_keyReplace(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")
	r := DiskEncryptionSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyReplace(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDiskEncryptionSet_keyRotateValid(t *testing.T) {
	// Regression test case for issue #22864
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")
	r := DiskEncryptionSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyRotateValid(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDiskEncryptionSet_keyRotateInvalid(t *testing.T) {
	// Regression test case for issue #22864
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")
	r := DiskEncryptionSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.keyRotateInvalid(data),
			ExpectError: regexp.MustCompile("'auto_key_rotation_enabled' field is set to 'true' expected a key vault key with a versionless ID but version information was found"),
		},
	})
}

func TestAccDiskEncryptionSet_keyRotateFalseInvalid(t *testing.T) {
	// Regression test case for issue #22864
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")
	r := DiskEncryptionSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.keyRotateFalseInvalid(data),
			ExpectError: regexp.MustCompile("'auto_key_rotation_enabled' field is set to 'false' expected a key vault key with a versioned ID but no version information was found"),
		},
	})
}

func TestAccDiskEncryptionSet_withEncryptionType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")
	r := DiskEncryptionSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withPlatformAndCustomerKeys(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDiskEncryptionSet_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")
	r := DiskEncryptionSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDiskEncryptionSet_systemAssignedUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")
	r := DiskEncryptionSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDiskEncryptionSet_updateIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")
	r := DiskEncryptionSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDiskEncryptionSet_withFederatedClientId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")
	r := DiskEncryptionSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withFederatedClientId(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withFederatedClientId(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDiskEncryptionSet_disablePurgeProtection(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_encryption_set", "test")
	r := DiskEncryptionSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disablePurgeProtection(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption_type").HasValue("EncryptionAtRestWithCustomerKey"),
			),
		},
		data.ImportStep(),
	})
}

func (DiskEncryptionSetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseDiskEncryptionSetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.DiskEncryptionSetsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute Disk Encryption Set %q", id.String())
	}

	model := resp.Model
	if model == nil {
		return nil, fmt.Errorf("retrieving Compute Disk Encryption Set %q: `model` was nil", id.String())
	}

	return utils.Bool(model.Id != nil), nil
}

func (DiskEncryptionSetResource) dependencies(data acceptance.TestData, purgeProtectionEnabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      recover_soft_deleted_key_vaults    = false
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                        = "acctestkv-%s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "standard"
  purge_protection_enabled    = %t
  enabled_for_disk_encryption = true
}

resource "azurerm_key_vault_access_policy" "service-principal" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Update",
    "GetRotationPolicy",
  ]

  secret_permissions = [
    "Get",
    "Delete",
    "Purge",
    "Set",
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "examplekey"
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

  depends_on = [azurerm_key_vault_access_policy.service-principal]
}

`, data.RandomInteger, data.Locations.Primary, data.RandomString, purgeProtectionEnabled)
}

func (r DiskEncryptionSetResource) systemAssignedDependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "disk-encryption" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "Get",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy",
  ]

  tenant_id = azurerm_disk_encryption_set.test.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.test.identity.0.principal_id
}
`, r.dependencies(data, true))
}

func (r DiskEncryptionSetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestDES-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id

  identity {
    type = "SystemAssigned"
  }
}
`, r.systemAssignedDependencies(data), data.RandomInteger)
}

func (r DiskEncryptionSetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_disk_encryption_set" "import" {
  name                = azurerm_disk_encryption_set.test.name
  resource_group_name = azurerm_disk_encryption_set.test.resource_group_name
  location            = azurerm_disk_encryption_set.test.location
  key_vault_key_id    = azurerm_disk_encryption_set.test.key_vault_key_id

  identity {
    type = "SystemAssigned"
  }
}
`, r.basic(data))
}

func (r DiskEncryptionSetResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_disk_encryption_set" "test" {
  name                      = "acctestDES-%d"
  resource_group_name       = azurerm_resource_group.test.name
  location                  = azurerm_resource_group.test.location
  key_vault_key_id          = azurerm_key_vault_key.test.versionless_id
  auto_key_rotation_enabled = true

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Hello = "woRld"
  }
}
`, r.systemAssignedDependencies(data), data.RandomInteger)
}

func (r DiskEncryptionSetResource) keyRotateInvalid(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_disk_encryption_set" "test" {
  name                      = "acctestDES-%d"
  resource_group_name       = azurerm_resource_group.test.name
  location                  = azurerm_resource_group.test.location
  key_vault_key_id          = azurerm_key_vault_key.test.id
  auto_key_rotation_enabled = true

  identity {
    type = "SystemAssigned"
  }
}
`, r.systemAssignedDependencies(data), data.RandomInteger)
}

func (r DiskEncryptionSetResource) keyRotateFalseInvalid(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestDES-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.versionless_id

  identity {
    type = "SystemAssigned"
  }
}
`, r.systemAssignedDependencies(data), data.RandomInteger)
}

func (r DiskEncryptionSetResource) keyReplace(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_key" "new" {
  name         = "newKey"
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

  depends_on = ["azurerm_key_vault_access_policy.service-principal"]
}

resource "azurerm_disk_encryption_set" "test" {
  name                      = "acctestDES-%d"
  resource_group_name       = azurerm_resource_group.test.name
  location                  = azurerm_resource_group.test.location
  key_vault_key_id          = azurerm_key_vault_key.new.versionless_id
  auto_key_rotation_enabled = true

  identity {
    type = "SystemAssigned"
  }
}
`, r.systemAssignedDependencies(data), data.RandomInteger)
}

func (r DiskEncryptionSetResource) keyRotateValid(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_disk_encryption_set" "test" {
  name                      = "acctestDES-%d"
  resource_group_name       = azurerm_resource_group.test.name
  location                  = azurerm_resource_group.test.location
  key_vault_key_id          = azurerm_key_vault_key.test.versionless_id
  auto_key_rotation_enabled = true

  identity {
    type = "SystemAssigned"
  }
}
`, r.systemAssignedDependencies(data), data.RandomInteger)
}

func (r DiskEncryptionSetResource) withPlatformAndCustomerKeys(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestDES-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id
  encryption_type     = "EncryptionAtRestWithPlatformAndCustomerKeys"

  identity {
    type = "SystemAssigned"
  }
}
`, r.systemAssignedDependencies(data), data.RandomInteger)
}

func (r DiskEncryptionSetResource) systemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault_access_policy" "user-assigned" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_user_assigned_identity.test.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id

  key_permissions = [
    "Get",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy",
  ]

  depends_on = [azurerm_key_vault_access_policy.service-principal]
}

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestDES-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id

  identity {
    type = "SystemAssigned"
  }

  depends_on = [azurerm_key_vault_access_policy.user-assigned]
}
	`, r.dependencies(data, true), data.RandomInteger)
}

func (r DiskEncryptionSetResource) userAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault_access_policy" "user-assigned" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_user_assigned_identity.test.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id

  key_permissions = [
    "Get",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy",
  ]

  depends_on = [azurerm_key_vault_access_policy.service-principal]
}

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestDES-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  depends_on = [azurerm_key_vault_access_policy.user-assigned]
}
	`, r.dependencies(data, true), data.RandomInteger)
}

func (r DiskEncryptionSetResource) systemAssignedUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault_access_policy" "user-assigned" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_user_assigned_identity.test.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id

  key_permissions = [
    "Get",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy",
  ]

  depends_on = [azurerm_key_vault_access_policy.service-principal]
}

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestDES-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  depends_on = [azurerm_key_vault_access_policy.user-assigned]
}
	`, r.systemAssignedDependencies(data), data.RandomInteger)
}

func (r DiskEncryptionSetResource) withFederatedClientId(data acceptance.TestData, clean bool) string {
	federatedClientId := `federated_client_id = azurerm_user_assigned_identity.test.client_id`
	if clean {
		federatedClientId = ""
	}

	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault_access_policy" "user-assigned" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_user_assigned_identity.test.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id

  key_permissions = [
    "Get",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy",
  ]

  depends_on = [azurerm_key_vault_access_policy.service-principal]
}

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestDES-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id
%[3]s

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  depends_on = [azurerm_key_vault_access_policy.user-assigned]
}
	`, r.dependencies(data, true), data.RandomInteger, federatedClientId)
}

func (r DiskEncryptionSetResource) disablePurgeProtection(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "disk-encryption" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "Get",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy",
  ]

  tenant_id = azurerm_disk_encryption_set.test.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.test.identity.0.principal_id
}

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestDES-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id

  identity {
    type = "SystemAssigned"
  }
}
`, r.dependencies(data, false), data.RandomInteger)
}
