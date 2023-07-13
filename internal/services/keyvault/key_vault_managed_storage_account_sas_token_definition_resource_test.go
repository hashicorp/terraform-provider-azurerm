// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KeyVaultManagedStorageAccountSasTokenDefinitionResource struct{}

func TestAccKeyVaultManagedStorageAccountSasTokenDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_storage_account_sas_token_definition", "test")
	r := KeyVaultManagedStorageAccountSasTokenDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("managed_storage_account_id"),
	})
}

func TestAccKeyVaultManagedStorageAccountSasTokenDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_storage_account_sas_token_definition", "test")
	r := KeyVaultManagedStorageAccountSasTokenDefinitionResource{}

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

func TestAccKeyVaultManagedStorageAccountSasTokenDefinition_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_storage_account_sas_token_definition", "test")
	r := KeyVaultManagedStorageAccountSasTokenDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "P1D"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sas_type").HasValue("account"),
				check.That(data.ResourceName).Key("validity_period").HasValue("P1D"),
				check.That(data.ResourceName).Key("secret_id").HasValue(fmt.Sprintf("https://acctestkv-%s.vault.azure.net/secrets/acctestKVstorage-acctestKVsasdefinition", data.RandomString)),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.hello").HasValue("world"),
			),
		},
		data.ImportStep("managed_storage_account_id"),
	})
}

func TestAccKeyVaultManagedStorageAccountSasTokenDefinition_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_storage_account_sas_token_definition", "test")
	r := KeyVaultManagedStorageAccountSasTokenDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "P1D"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sas_type").HasValue("account"),
				check.That(data.ResourceName).Key("validity_period").HasValue("P1D"),
				check.That(data.ResourceName).Key("secret_id").HasValue(fmt.Sprintf("https://acctestkv-%s.vault.azure.net/secrets/acctestKVstorage-acctestKVsasdefinition", data.RandomString)),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.hello").HasValue("world"),
			),
		},
		data.ImportStep("managed_storage_account_id"),
		{
			Config: r.complete(data, "P2D"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sas_type").HasValue("account"),
				check.That(data.ResourceName).Key("validity_period").HasValue("P2D"),
				check.That(data.ResourceName).Key("secret_id").HasValue(fmt.Sprintf("https://acctestkv-%s.vault.azure.net/secrets/acctestKVstorage-acctestKVsasdefinition", data.RandomString)),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.hello").HasValue("world"),
			),
		},
		data.ImportStep("managed_storage_account_id"),
	})
}

func TestAccKeyVaultManagedStorageAccountSasTokenDefinition_recovery(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_storage_account_sas_token_definition", "test")
	r := KeyVaultManagedStorageAccountSasTokenDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.softDeleteRecovery(data, false, "1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("managed_storage_account_id"),
		{
			Config:  r.softDeleteRecovery(data, false, "1"),
			Destroy: true,
		},
		{
			// purge true here to make sure when we end the test there's no soft-deleted items left behind
			Config: r.softDeleteRecovery(data, true, "2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("managed_storage_account_id"),
	})
}

func (r KeyVaultManagedStorageAccountSasTokenDefinitionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_managed_storage_account" "test" {
  name                         = "acctestKVstorage"
  key_vault_id                 = azurerm_key_vault.test.id
  storage_account_id           = azurerm_storage_account.test.id
  storage_account_key          = "key1"
  regenerate_key_automatically = false
  regeneration_period          = "P1D"
}

resource "azurerm_key_vault_managed_storage_account_sas_token_definition" "test" {
  name                       = "acctestKVsasdefinition"
  managed_storage_account_id = azurerm_key_vault_managed_storage_account.test.id
  sas_type                   = "account"
  sas_template_uri           = data.azurerm_storage_account_sas.test.sas
  validity_period            = "P1D"
}
`, r.template(data))
}

func (r KeyVaultManagedStorageAccountSasTokenDefinitionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_managed_storage_account_sas_token_definition" "import" {
  name                       = azurerm_key_vault_managed_storage_account_sas_token_definition.test.name
  managed_storage_account_id = azurerm_key_vault_managed_storage_account.test.id
  sas_type                   = "account"
  sas_template_uri           = data.azurerm_storage_account_sas.test.sas
  validity_period            = "P1D"
}
`, r.basic(data))
}

func (r KeyVaultManagedStorageAccountSasTokenDefinitionResource) complete(data acceptance.TestData, validyPeriod string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_managed_storage_account" "test" {
  name                         = "acctestKVstorage"
  key_vault_id                 = azurerm_key_vault.test.id
  storage_account_id           = azurerm_storage_account.test.id
  storage_account_key          = "key1"
  regenerate_key_automatically = false
  regeneration_period          = "P1D"
}

resource "azurerm_key_vault_managed_storage_account_sas_token_definition" "test" {
  name                       = "acctestKVsasdefinition"
  managed_storage_account_id = azurerm_key_vault_managed_storage_account.test.id
  sas_type                   = "account"
  sas_template_uri           = data.azurerm_storage_account_sas.test.sas
  validity_period            = "%s"

  tags = {
    "hello" = "world"
  }
}
`, r.template(data), validyPeriod)
}

func (r KeyVaultManagedStorageAccountSasTokenDefinitionResource) softDeleteRecovery(data acceptance.TestData, purge bool, name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy    = "%t"
      recover_soft_deleted_key_vaults = true
    }
  }
}

%s

resource "azurerm_key_vault_managed_storage_account" "test" {
  name                         = "acctestKVstorage%s"
  key_vault_id                 = azurerm_key_vault.test.id
  storage_account_id           = azurerm_storage_account.test.id
  storage_account_key          = "key1"
  regenerate_key_automatically = false
  regeneration_period          = "P1D"
}

resource "azurerm_key_vault_managed_storage_account_sas_token_definition" "test" {
  name                       = "acctestKVsasdefinition%s"
  managed_storage_account_id = azurerm_key_vault_managed_storage_account.test.id
  sas_type                   = "account"
  sas_template_uri           = data.azurerm_storage_account_sas.test.sas
  validity_period            = "P1D"
}
`, purge, r.template(data), name, name)
}

func (KeyVaultManagedStorageAccountSasTokenDefinitionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctest-kv-RG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

data "azurerm_storage_account_sas" "test" {
  connection_string = azurerm_storage_account.test.primary_connection_string
  https_only        = true

  resource_types {
    service   = true
    container = false
    object    = false
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = "2021-04-30T00:00:00Z"
  expiry = "2023-04-30T00:00:00Z"

  permissions {
    read    = true
    write   = true
    delete  = false
    list    = false
    add     = true
    create  = true
    update  = false
    process = false
    tag     = false
    filter  = false
  }
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Get",
      "Delete"
    ]

    storage_permissions = [
      "Get",
      "List",
      "Set",
      "SetSAS",
      "GetSAS",
      "DeleteSAS",
      "Update",
      "RegenerateKey"
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (KeyVaultManagedStorageAccountSasTokenDefinitionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	dataPlaneClient := client.KeyVault.ManagementClient
	keyVaultsClient := client.KeyVault

	id, err := parse.SasDefinitionID(state.ID)
	if err != nil {
		return nil, err
	}

	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, client.Resource, id.KeyVaultBaseUrl)
	if err != nil || keyVaultIdRaw == nil {
		return nil, fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return nil, err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil || !ok {
		return nil, fmt.Errorf("checking if key vault %q for Managed Storage Account Sas Definition %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}

	resp, err := dataPlaneClient.GetSasDefinition(ctx, id.KeyVaultBaseUrl, id.StorageAccountName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Key Vault Managed Storage Account Sas Definition %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.ID != nil), nil
}
