// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-05-01/storageaccounts"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StorageAccountCustomerManagedKeyResource struct{}

func TestAccStorageAccountCustomerManagedKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	r := StorageAccountCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_storage_account_customer_managed_key.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Delete the encryption settings resource and verify it is gone
			Config: r.template(data),
			Check: acceptance.ComposeTestCheckFunc(
				// Then ensure the encryption settings on the storage account
				// have been reverted to their default state
				data.CheckWithClient(r.accountHasDefaultSettings),
			),
		},
	})
}

func TestAccStorageAccountCustomerManagedKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_customer_managed_key", "test")
	r := StorageAccountCustomerManagedKeyResource{}

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

func TestAccStorageAccountCustomerManagedKey_updateKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_customer_managed_key", "test")
	r := StorageAccountCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageAccountCustomerManagedKey_testKeyVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_customer_managed_key", "test")
	r := StorageAccountCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.autoKeyRotation(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageAccountCustomerManagedKey_managedHSM(t *testing.T) {
	if os.Getenv("ARM_TEST_HSM_KEY") == "" {
		t.Skip("Skipping as ARM_TEST_HSM_KEY is not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_account_customer_managed_key", "test")
	r := StorageAccountCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedHSM(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.managedHSMUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageAccountCustomerManagedKey_remoteKeyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_customer_managed_key", "test")
	r := StorageAccountCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.remoteKeyVault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccStorageAccountCustomerManagedKey_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_customer_managed_key", "test")
	r := StorageAccountCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("user_assigned_identity_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageAccountCustomerManagedKey_userAssignedIdentityWithFederatedIdentity(t *testing.T) {
	// Multiple tenants are needed for this test
	altTenantId := os.Getenv("ARM_TENANT_ID_ALT")
	subscriptionIdAltTenant := os.Getenv("ARM_SUBSCRIPTION_ID_ALT_TENANT")

	if altTenantId == "" || subscriptionIdAltTenant == "" {
		t.Skip("One of ARM_TENANT_ID_ALT, ARM_SUBSCRIPTION_ID_ALT_TENANT are not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_account_customer_managed_key", "test")
	r := StorageAccountCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.federatedIdentity(data, altTenantId, subscriptionIdAltTenant),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("federated_identity_client_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageAccountCustomerManagedKeyResource) accountHasDefaultSettings(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	accountId, err := commonids.ParseStorageAccountID(state.Attributes["id"])
	if err != nil {
		return err
	}

	resp, err := client.Storage.ResourceManager.StorageAccounts.GetProperties(ctx, *accountId, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("Bad: Get on storageServiceClient: %+v", err)
	}

	if response.WasNotFound(resp.HttpResponse) {
		return fmt.Errorf("Bad: %s does not exist", accountId)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if encryption := props.Encryption; encryption != nil {
				if services := encryption.Services; services != nil {
					if !*services.Blob.Enabled {
						return fmt.Errorf("enable_blob_encryption not set to default")
					}
					if !*services.File.Enabled {
						return fmt.Errorf("enable_file_encryption not set to default")
					}
				}

				if encryption.KeySource != nil && *encryption.KeySource != storageaccounts.KeySourceMicrosoftPointStorage {
					return fmt.Errorf("%q should be %q", *encryption.KeySource, string(storageaccounts.KeySourceMicrosoftPointStorage))
				}
			} else {
				return fmt.Errorf("storage account encryption properties not found")
			}
		}
	}

	return nil
}

func (r StorageAccountCustomerManagedKeyResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	accountId, err := commonids.ParseStorageAccountID(state.Attributes["storage_account_id"])
	if err != nil {
		return nil, err
	}

	resp, err := client.Storage.ResourceManager.StorageAccounts.GetProperties(ctx, *accountId, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}

		return nil, fmt.Errorf("Bad: Get on storageServiceClient: %+v", err)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if encryption := props.Encryption; encryption != nil {
				if encryption.KeySource != nil && *encryption.KeySource == storageaccounts.KeySourceMicrosoftPointKeyvault {
					return pointer.To(true), nil
				}

				return nil, fmt.Errorf("%q should be %q", *encryption.KeySource, string(storageaccounts.KeySourceMicrosoftPointKeyvault))
			}
		}
	}

	return pointer.To(false), nil
}

func (r StorageAccountCustomerManagedKeyResource) basic(data acceptance.TestData) string {
	if !features.FivePointOh() {
		return fmt.Sprintf(`
%s

resource "azurerm_storage_account_customer_managed_key" "test" {
  storage_account_id = azurerm_storage_account.test.id
  key_vault_id       = azurerm_key_vault.test.id
  key_name           = azurerm_key_vault_key.first.name
  key_version        = azurerm_key_vault_key.first.version
}
`, r.template(data))
	}

	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_customer_managed_key" "test" {
  storage_account_id = azurerm_storage_account.test.id
  key_vault_key_id   = azurerm_key_vault_key.first.id
}
`, r.template(data))
}

func (r StorageAccountCustomerManagedKeyResource) requiresImport(data acceptance.TestData) string {
	if !features.FivePointOh() {
		return fmt.Sprintf(`
%s

resource "azurerm_storage_account_customer_managed_key" "import" {
  storage_account_id = azurerm_storage_account_customer_managed_key.test.storage_account_id
  key_vault_id       = azurerm_storage_account_customer_managed_key.test.key_vault_id
  key_name           = azurerm_storage_account_customer_managed_key.test.key_name
  key_version        = azurerm_storage_account_customer_managed_key.test.key_version
}
`, r.basic(data))
	}

	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_customer_managed_key" "import" {
  storage_account_id = azurerm_storage_account_customer_managed_key.test.storage_account_id
  key_vault_key_id   = azurerm_storage_account_customer_managed_key.test.key_vault_key_id
}
`, r.basic(data))
}

func (r StorageAccountCustomerManagedKeyResource) updated(data acceptance.TestData) string {
	if !features.FivePointOh() {
		return fmt.Sprintf(`
%s

resource "azurerm_key_vault_key" "second" {
  name         = "second"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.storage,
  ]
}

resource "azurerm_storage_account_customer_managed_key" "test" {
  storage_account_id = azurerm_storage_account.test.id
  key_vault_id       = azurerm_key_vault.test.id
  key_name           = azurerm_key_vault_key.second.name
  key_version        = azurerm_key_vault_key.second.version
}
`, r.template(data))
	}

	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_key" "second" {
  name         = "second"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.storage,
  ]
}

resource "azurerm_storage_account_customer_managed_key" "test" {
  storage_account_id = azurerm_storage_account.test.id
  key_vault_key_id   = azurerm_key_vault_key.second.id
}
`, r.template(data))
}

func (r StorageAccountCustomerManagedKeyResource) autoKeyRotation(data acceptance.TestData) string {
	if !features.FivePointOh() {
		return fmt.Sprintf(`
%s

resource "azurerm_storage_account_customer_managed_key" "test" {
  storage_account_id = azurerm_storage_account.test.id
  key_vault_id       = azurerm_key_vault.test.id
  key_name           = azurerm_key_vault_key.first.name
}
`, r.template(data))
	}

	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_customer_managed_key" "test" {
  storage_account_id = azurerm_storage_account.test.id
  key_vault_key_id   = azurerm_key_vault_key.first.versionless_id
}
`, r.template(data))
}

func (r StorageAccountCustomerManagedKeyResource) managedHSM(data acceptance.TestData) string {
	if !features.FivePointOh() {
		return fmt.Sprintf(`
%s

resource "azurerm_storage_account_customer_managed_key" "test" {
  storage_account_id = azurerm_storage_account.test.id
  managed_hsm_key_id = azurerm_key_vault_managed_hardware_security_module_key.test.versioned_id
}
`, r.templateHSM(data))
	}

	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_customer_managed_key" "test" {
  storage_account_id = azurerm_storage_account.test.id
  key_vault_key_id   = azurerm_key_vault_managed_hardware_security_module_key.test.versioned_id
}
`, r.templateHSM(data))
}

func (r StorageAccountCustomerManagedKeyResource) managedHSMUpdate(data acceptance.TestData) string {
	if !features.FivePointOh() {
		return fmt.Sprintf(`
%s

resource "azurerm_storage_account_customer_managed_key" "test" {
  storage_account_id = azurerm_storage_account.test.id
  managed_hsm_key_id = azurerm_key_vault_managed_hardware_security_module_key.test.id
}
`, r.templateHSM(data))
	}

	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_customer_managed_key" "test" {
  storage_account_id = azurerm_storage_account.test.id
  key_vault_key_id   = azurerm_key_vault_managed_hardware_security_module_key.test.id
}
`, r.templateHSM(data))
}

func (r StorageAccountCustomerManagedKeyResource) userAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_customer_managed_key" "test" {
  storage_account_id        = azurerm_storage_account.test.id
  key_vault_key_id          = azurerm_key_vault_key.first.id
  user_assigned_identity_id = azurerm_user_assigned_identity.test.id
}
`, r.userAssignedIdentityTemplate(data))
}

func (r StorageAccountCustomerManagedKeyResource) remoteKeyVault(data acceptance.TestData) string {
	clientData := data.Client()
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azurerm-alt" {
  subscription_id = "%s"
  tenant_id       = "%s"
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "remotetest" {
  provider = azurerm-alt

  name     = "acctestRG-alt-%d"
  location = "%s"
}

resource "azurerm_key_vault" "remotetest" {
  provider = azurerm-alt

  name                       = "acctestkv%s"
  location                   = azurerm_resource_group.remotetest.location
  resource_group_name        = azurerm_resource_group.remotetest.name
  tenant_id                  = "%s"
  sku_name                   = "standard"
  purge_protection_enabled   = true
  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_access_policy" "storage" {
  provider = azurerm-alt

  key_vault_id = azurerm_key_vault.remotetest.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_storage_account.test.identity.0.principal_id

  key_permissions    = ["Get", "Create", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_access_policy" "client" {
  provider = azurerm-alt

  key_vault_id = azurerm_key_vault.remotetest.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_key" "remote" {
  provider = azurerm-alt

  name         = "remote"
  key_vault_id = azurerm_key_vault.remotetest.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.storage,
  ]
}

resource "azurerm_resource_group" "test" {
  provider = azurerm

  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  provider = azurerm

  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  identity {
    type = "SystemAssigned"
  }

  lifecycle {
    ignore_changes = ["customer_managed_key"]
  }
}

resource "azurerm_storage_account_customer_managed_key" "test" {
  provider = azurerm

  storage_account_id = azurerm_storage_account.test.id
  key_vault_key_id   = azurerm_key_vault_key.remote.id
}
`, clientData.SubscriptionIDAlt, clientData.TenantID, data.RandomInteger, data.Locations.Primary, data.RandomString, clientData.TenantID, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageAccountCustomerManagedKeyResource) userAssignedIdentityTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
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

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestmi%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv%s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  purge_protection_enabled   = true
  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_access_policy" "storage" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id

  key_permissions    = ["Get", "Create", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_key" "first" {
  name         = "first"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.storage,
  ]
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  lifecycle {
    ignore_changes = ["customer_managed_key"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString)
}

func (r StorageAccountCustomerManagedKeyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
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
  name                       = "acctestkv%s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  purge_protection_enabled   = true
  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_access_policy" "storage" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_storage_account.test.identity.0.principal_id

  key_permissions    = ["Get", "Create", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_key" "first" {
  name         = "first"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.storage,
  ]
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  identity {
    type = "SystemAssigned"
  }

  lifecycle {
    ignore_changes = ["customer_managed_key"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (r StorageAccountCustomerManagedKeyResource) federatedIdentity(data acceptance.TestData, altTenantId, subscriptionIdAltTenant string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azurerm-alt" {
  tenant_id       = "%[1]s"
  subscription_id = "%[2]s"

  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

provider "azuread" {}

provider "azuread" {
  alias     = "alt"
  tenant_id = "%[1]s"
}

data "azurerm_client_config" "current" {}

data "azurerm_client_config" "remote" {
  provider = azurerm-alt
}

data "azuread_client_config" "current" {}

data "azuread_client_config" "remote" {
  provider = azuread.alt
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[3]d"
  location = "%[4]s"
}

resource "azuread_application" "test" {
  display_name     = "acctestapp-%[5]s"
  sign_in_audience = "AzureADMultipleOrgs"
  owners           = [data.azuread_client_config.current.object_id]
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestmi-%[5]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azuread_application_federated_identity_credential" "test" {
  application_id = azuread_application.test.id
  display_name   = "acctestcred-%[5]s"
  description    = "Federated Identity Credential for CMK"
  audiences      = ["api://AzureADTokenExchange"]
  issuer         = "https://login.microsoftonline.com/${data.azurerm_client_config.current.tenant_id}/v2.0"
  subject        = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_resource_group" "remotetest" {
  provider = azurerm-alt
  name     = "acctestRG-alt-%[3]d"
  location = "%[4]s"
}

resource "azuread_service_principal" "remotetest" {
  provider  = azuread.alt
  owners    = [data.azuread_client_config.remote.object_id]
  client_id = azuread_application.test.client_id
}

resource "azurerm_key_vault" "remotetest" {
  provider = azurerm-alt

  name                       = "acctestkv%[5]s"
  location                   = azurerm_resource_group.remotetest.location
  resource_group_name        = azurerm_resource_group.remotetest.name
  tenant_id                  = data.azurerm_client_config.remote.tenant_id
  sku_name                   = "standard"
  purge_protection_enabled   = true
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.remote.tenant_id
    object_id = data.azurerm_client_config.remote.object_id

    key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
    secret_permissions = ["Get"]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.remote.tenant_id
    object_id = azuread_service_principal.remotetest.object_id

    key_permissions = [
      "Get", "List", "UnwrapKey", "WrapKey",
    ]
  }

}

resource "azurerm_key_vault_key" "remotetest" {
  provider = azurerm-alt

  name         = "remote"
  key_vault_id = azurerm_key_vault.remotetest.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[5]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  lifecycle {
    ignore_changes = [customer_managed_key]
  }
}

resource "azurerm_storage_account_customer_managed_key" "test" {
  storage_account_id = azurerm_storage_account.test.id
  key_vault_key_id   = azurerm_key_vault_key.remotetest.id

  user_assigned_identity_id    = azurerm_user_assigned_identity.test.id
  federated_identity_client_id = azuread_application.test.client_id
}
`, altTenantId, subscriptionIdAltTenant, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageAccountCustomerManagedKeyResource) templateHSM(data acceptance.TestData) string {
	uuid1, _ := uuid.GenerateUUID()
	uuid2, _ := uuid.GenerateUUID()
	uuid3, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-KV-%[1]s"
  location = "%[2]s"
}

resource "azurerm_key_vault" "test" {
  name                       = "acc%[3]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Recover",
      "Update",
      "GetRotationPolicy",
    ]

    secret_permissions = [
      "Delete",
      "Get",
      "Set",
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "DeleteIssuers",
      "Get",
      "Purge",
      "Update"
    ]
  }
}

resource "azurerm_key_vault_certificate" "cert" {
  count        = 3
  name         = "acchsmcert${count.index}"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      extended_key_usage = []

      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=hello-world"
      validity_in_months = 12
    }
  }
}

resource "azurerm_key_vault_managed_hardware_security_module" "test" {
  name                       = "kvHsm%[3]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  sku_name                   = "Standard_B1"
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  admin_object_ids           = [data.azurerm_client_config.current.object_id]
  purge_protection_enabled   = false
  soft_delete_retention_days = 7

  security_domain_key_vault_certificate_ids = [for cert in azurerm_key_vault_certificate.cert : cert.id]
  security_domain_quorum                    = 3
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestmi%[1]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "test" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[4]s"
  scope              = "/keys"
  role_definition_id = "/Microsoft.KeyVault/providers/Microsoft.Authorization/roleDefinitions/21dbd100-6940-42c2-9190-5d6cb909625b"
  principal_id       = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "test1" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[5]s"
  scope              = "/keys"
  role_definition_id = "/Microsoft.KeyVault/providers/Microsoft.Authorization/roleDefinitions/515eb02d-2335-4d2d-92f2-b1cbdf9c3778"
  principal_id       = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "user" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[6]s"
  scope              = "/keys"
  role_definition_id = "/Microsoft.KeyVault/providers/Microsoft.Authorization/roleDefinitions/21dbd100-6940-42c2-9190-5d6cb909625b"
  principal_id       = azurerm_storage_account.test.identity.0.principal_id
}

resource "azurerm_key_vault_managed_hardware_security_module_key" "test" {
  name           = "acctestHSMK-%[1]s"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  key_type       = "RSA-HSM"
  key_size       = 2048
  key_opts       = ["unwrapKey", "wrapKey"]

  depends_on = [
    azurerm_key_vault_managed_hardware_security_module_role_assignment.test,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.test1,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.user
  ]
}

resource "azurerm_key_vault_managed_hardware_security_module_key" "update" {
  name           = "acctestHSMK2-%[1]s"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  key_type       = "RSA-HSM"
  key_size       = 2048
  key_opts       = ["unwrapKey", "wrapKey"]

  depends_on = [
    azurerm_key_vault_managed_hardware_security_module_role_assignment.test,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.test1,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.user
  ]
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[1]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  identity {
    type = "SystemAssigned"
  }

  lifecycle {
    ignore_changes = ["customer_managed_key"]
  }
}
`, data.RandomString, data.Locations.Primary, data.RandomInteger, uuid1, uuid2, uuid3)
}
