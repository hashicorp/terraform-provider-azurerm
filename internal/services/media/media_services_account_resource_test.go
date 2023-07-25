// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package media_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/accounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MediaServicesAccountResource struct{}

func TestAccMediaServicesAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")
	r := MediaServicesAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("storage_account.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaServicesAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")
	r := MediaServicesAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("storage_account.#").HasValue("1"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMediaServicesAccount_multipleAccounts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")
	r := MediaServicesAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleAccounts(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_account.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config:   r.multipleAccountsUpdated(data),
			PlanOnly: true,
		},
		data.ImportStep(),
	})
}

func TestAccMediaServicesAccount_storageAuthWithManagedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")
	r := MediaServicesAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAuthWithManagedIdentity(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaServicesAccount_multiplePrimaries(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")
	r := MediaServicesAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.multiplePrimaries(data),
			ExpectError: regexp.MustCompile("Only one Storage Account can be set as Primary"),
		},
	})
}

func TestAccMediaServicesAccount_publicNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")
	r := MediaServicesAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.publicNetworkAccess(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.publicNetworkAccess(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.publicNetworkAccess(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaServicesAccount_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")
	r := MediaServicesAccountResource{}

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
		{
			Config: r.systemAndUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaServicesAccount_encryptionCustomerKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")
	r := MediaServicesAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.encryptionSystemKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaServicesAccount_encryptionUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")
	r := MediaServicesAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.encryptionSystemKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.encryptionCustomerKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.encryptionSystemKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaServicesAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")
	r := MediaServicesAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
			),
		},
		data.ImportStep(),
	})
}

func (MediaServicesAccountResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := accounts.ParseMediaServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.V20211101Client.Accounts.MediaservicesGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r MediaServicesAccountResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }

  tags = {
    environment = "staging"
  }
}
`, template, data.RandomString)
}

func (r MediaServicesAccountResource) publicNetworkAccess(data acceptance.TestData, enabled bool) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_services_account" "test" {
  name                          = "acctestmsa%s"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  public_network_access_enabled = %t

  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }

  tags = {
    environment = "staging"
  }
}
`, template, data.RandomString, enabled)
}

func (r MediaServicesAccountResource) encryptionSystemKey(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
  encryption {
    type = "SystemKey"
  }

  tags = {
    environment = "staging"
  }
}
`, template, data.RandomString)
}

func (r MediaServicesAccountResource) encryptionCustomerKey(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault" "test" {
  name                       = "mediakv-%s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  purge_protection_enabled   = true
}

data "azurerm_client_config" "current" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_key_vault_access_policy" "server" {
  key_vault_id       = azurerm_key_vault.test.id
  tenant_id          = data.azurerm_client_config.current.tenant_id
  object_id          = azurerm_user_assigned_identity.test.principal_id
  key_permissions    = ["Get", "UnwrapKey", "WrapKey", "List", "Encrypt"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id       = azurerm_key_vault.test.id
  tenant_id          = data.azurerm_client_config.current.tenant_id
  object_id          = data.azurerm_client_config.current.object_id
  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_key" "test" {
  name         = "test-kvk-%[2]s"
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
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.server,
  ]
}

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }
  encryption {
    type                     = "CustomerKey"
    key_vault_key_identifier = azurerm_key_vault_key.test.id
    managed_identity {
      user_assigned_identity_id = azurerm_user_assigned_identity.test.id
    }
  }

  tags = {
    environment = "staging"
  }
}
`, template, data.RandomString)
}

func (r MediaServicesAccountResource) systemAssignedIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }
  identity {
    type = "SystemAssigned"
  }
  tags = {
    environment = "staging"
  }
}
`, template, data.RandomString)
}

func (r MediaServicesAccountResource) userAssignedIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
  tags = {
    environment = "staging"
  }
}
`, template, data.RandomString)
}

func (r MediaServicesAccountResource) systemAndUserAssignedIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
  tags = {
    environment = "staging"
  }
}
`, template, data.RandomString)
}

func (r MediaServicesAccountResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_services_account" "import" {
  name                = azurerm_media_services_account.test.name
  location            = azurerm_media_services_account.test.location
  resource_group_name = azurerm_media_services_account.test.resource_group_name

  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }

  tags = {
    environment = "staging"
  }
}
`, template)
}

func (r MediaServicesAccountResource) storageAuthWithManagedIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_storage_account" "third" {
  name                     = "acctestsa3%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}

resource "azurerm_role_assignment" "storageAccountRoleAssignment" {
  scope                = azurerm_storage_account.third.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_media_services_account" "test" {
  name                        = "acctestmsa%[2]s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  storage_authentication_type = "ManagedIdentity"

  storage_account {
    id         = azurerm_storage_account.third.id
    is_primary = true
    managed_identity {
      user_assigned_identity_id    = azurerm_user_assigned_identity.test.id
      use_system_assigned_identity = false
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    environment = "staging"
  }

  depends_on = [
    azurerm_role_assignment.storageAccountRoleAssignment
  ]
}
`, template, data.RandomString)
}

func (r MediaServicesAccountResource) multipleAccounts(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "second" {
  name                     = "acctestsa2%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }

  storage_account {
    id         = azurerm_storage_account.second.id
    is_primary = false
  }
}
`, template, data.RandomString, data.RandomString)
}

func (r MediaServicesAccountResource) multipleAccountsUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "second" {
  name                     = "acctestsa2%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  storage_account {
    id         = azurerm_storage_account.second.id
    is_primary = false
  }

  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }
}
`, template, data.RandomString, data.RandomString)
}

func (r MediaServicesAccountResource) multiplePrimaries(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "second" {
  name                     = "acctestsa2%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }

  storage_account {
    id         = azurerm_storage_account.second.id
    is_primary = true
  }
}
`, template, data.RandomString, data.RandomString)
}

func (r MediaServicesAccountResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }

  tags = {
    environment = "staging"
  }

  identity {
    type = "SystemAssigned"
  }

  key_delivery_access_control {
    default_action = "Deny"
    ip_allow_list  = ["0.0.0.0/0"]
  }
}
`, template, data.RandomString)
}

func (MediaServicesAccountResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-media-%d"
  location = "%s"
}

resource "azurerm_storage_account" "first" {
  name                     = "acctestsa1%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
