package keyvault_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type KeyVaultSecretResource struct {
}

func TestAccKeyVaultSecret_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")
	r := KeyVaultSecretResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("value").HasValue("rick-and-morty"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVaultSecret_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")
	r := KeyVaultSecretResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("value").HasValue("rick-and-morty"),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_key_vault_secret"),
		},
	})
}

func TestAccKeyVaultSecret_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")
	r := KeyVaultSecretResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccKeyVaultSecret_disappearsWhenParentKeyVaultDeleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")
	r := KeyVaultSecretResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.destroyParentKeyVault, "azurerm_key_vault.test"),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccKeyVaultSecret_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")
	r := KeyVaultSecretResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("not_before_date").HasValue("2019-01-01T01:02:03Z"),
				check.That(data.ResourceName).Key("expiration_date").HasValue("2020-01-01T01:02:03Z"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.hello").HasValue("world"),
				check.That(data.ResourceName).Key("versionless_id").HasValue(fmt.Sprintf("https://acctestkv-%s.vault.azure.net/secrets/secret-%s", data.RandomString, data.RandomString)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVaultSecret_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")
	r := KeyVaultSecretResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("value").HasValue("rick-and-morty"),
			),
		},
		{
			Config: r.basicUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("value").HasValue("szechuan"),
			),
		},
	})
}

func TestAccKeyVaultSecret_updatingValueChangedExternally(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")
	r := KeyVaultSecretResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("value").HasValue("rick-and-morty"),
				data.CheckWithClient(r.updateSecretValue("mad-scientist")),
			),
			ExpectNonEmptyPlan: true,
		},
		{
			Config: r.updateTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:   r.updateTags(data),
			PlanOnly: true,
		},
		data.ImportStep(),
	})
}

func TestAccKeyVaultSecret_recovery(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")
	r := KeyVaultSecretResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.softDeleteRecovery(data, false, "first"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("value").HasValue("first"),
			),
		},
		{
			Config:  r.softDeleteRecovery(data, false, "first"),
			Destroy: true,
		},
		{
			// purge true here to make sure when we end the test there's no soft-deleted items left behind
			Config: r.softDeleteRecovery(data, true, "second"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("value").HasValue("second"),
			),
		},
	})
}

func TestAccKeyVaultSecret_withExternalAccessPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")
	r := KeyVaultSecretResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withExternalAccessPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withExternalAccessPolicyUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVaultSecret_purge(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_secret", "test")
	r := KeyVaultSecretResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("value").HasValue("rick-and-morty"),
			),
		},
		{
			Config:  r.basic(data),
			Destroy: true,
		},
	})
}

func (KeyVaultSecretResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.KeyVault.ManagementClient
	keyVaultsClient := clients.KeyVault

	id, err := parse.ParseNestedItemID(state.ID)
	if err != nil {
		return nil, err
	}

	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, clients.Resource, id.KeyVaultBaseUrl)
	if err != nil || keyVaultIdRaw == nil {
		return nil, fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	keyVaultId, err := parse.VaultID(*keyVaultIdRaw)
	if err != nil {
		return nil, err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil || !ok {
		return nil, fmt.Errorf("checking if key vault %q for Certificate %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}

	// we always want to get the latest version
	resp, err := client.GetSecret(ctx, id.KeyVaultBaseUrl, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("making Read request on Azure KeyVault Secret %s: %+v", id.Name, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (KeyVaultSecretResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	dataPlaneClient := client.KeyVault.ManagementClient

	name := state.Attributes["name"]
	keyVaultId, err := parse.VaultID(state.Attributes["key_vault_id"])
	if err != nil {
		return nil, err
	}
	vaultBaseUrl, err := client.KeyVault.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return nil, fmt.Errorf("looking up Secret %q vault url from id %q: %+v", name, keyVaultId, err)
	}

	if _, err := dataPlaneClient.DeleteSecret(ctx, *vaultBaseUrl, name); err != nil {
		return nil, fmt.Errorf("Bad: Delete on keyVaultManagementClient: %+v", err)
	}

	return utils.Bool(true), nil
}

func (KeyVaultSecretResource) destroyParentKeyVault(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	ok, err := KeyVaultResource{}.Destroy(ctx, client, state)
	if err != nil {
		return err
	}

	if ok == nil || !*ok {
		return fmt.Errorf("deleting parent key vault failed")
	}

	return nil
}

func (r KeyVaultSecretResource) updateSecretValue(value string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		dataPlaneClient := clients.KeyVault.ManagementClient

		name := state.Attributes["name"]
		keyVaultId, err := parse.VaultID(state.Attributes["key_vault_id"])
		if err != nil {
			return err
		}

		vaultBaseUrl, err := clients.KeyVault.BaseUriForKeyVault(ctx, *keyVaultId)
		if err != nil {
			return fmt.Errorf("looking up Secret %q vault url from id %q: %+v", name, keyVaultId, err)
		}

		updated := keyvault.SecretSetParameters{
			Value: utils.String(value),
		}
		if _, err = dataPlaneClient.SetSecret(ctx, *vaultBaseUrl, name, updated); err != nil {
			return fmt.Errorf("updating secret: %+v", err)
		}
		return nil
	}
}

func (r KeyVaultSecretResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%s"
  value        = "rick-and-morty"
  key_vault_id = azurerm_key_vault.test.id
}
`, r.template(data), data.RandomString)
}

func (r KeyVaultSecretResource) updateTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%s"
  value        = "mad-scientist"
  key_vault_id = azurerm_key_vault.test.id

  tags = {
    Rick = "Morty"
  }
}
`, r.template(data), data.RandomString)
}

func (r KeyVaultSecretResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_secret" "import" {
  name         = azurerm_key_vault_secret.test.name
  value        = azurerm_key_vault_secret.test.value
  key_vault_id = azurerm_key_vault_secret.test.key_vault_id
}
`, r.basic(data))
}

func (r KeyVaultSecretResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_secret" "test" {
  name            = "secret-%s"
  value           = "<rick><morty /></rick>"
  key_vault_id    = azurerm_key_vault.test.id
  content_type    = "application/xml"
  not_before_date = "2019-01-01T01:02:03Z"
  expiration_date = "2020-01-01T01:02:03Z"

  tags = {
    "hello" = "world"
  }
}
`, r.template(data), data.RandomString)
}

func (r KeyVaultSecretResource) basicUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%s"
  value        = "szechuan"
  key_vault_id = azurerm_key_vault.test.id
}
`, r.template(data), data.RandomString)
}

func (r KeyVaultSecretResource) softDeleteRecovery(data acceptance.TestData, purge bool, value string) string {
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

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%s"
  value        = "%s"
  key_vault_id = azurerm_key_vault.test.id
}
`, purge, r.template(data), data.RandomString, value)
}

func (KeyVaultSecretResource) withExternalAccessPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv-%s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_enabled        = true
  soft_delete_retention_days = 7

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id
  key_permissions = [
    "Create",
    "Get",
  ]
  secret_permissions = [
    "Set",
    "Get",
    "Delete",
    "Purge",
    "Recover"
  ]
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%s"
  value        = "rick-and-morty"
  key_vault_id = azurerm_key_vault.test.id
  depends_on   = [azurerm_key_vault_access_policy.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (KeyVaultSecretResource) withExternalAccessPolicyUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv-%s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id
  key_permissions = [
    "Create",
    "Get",
  ]
  secret_permissions = [
    "Set",
    "Get",
    "Delete",
    "Purge",
    "Recover"
  ]
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%s"
  value        = "rick-and-morty"
  key_vault_id = azurerm_key_vault.test.id
  depends_on   = [azurerm_key_vault_access_policy.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (KeyVaultSecretResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv-%s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Get",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "List",
      "Purge",
      "Recover",
      "Set",
    ]
  }

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
