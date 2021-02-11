package keyvault_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type KeyVaultKeyResource struct {
}

func TestAccKeyVaultKey_basicEC(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	r := KeyVaultKeyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicEC(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("key_size"),
	})
}

func TestAccKeyVaultKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	r := KeyVaultKeyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicEC(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_key_vault_key"),
		},
	})
}

func TestAccKeyVaultKey_basicECHSM(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	r := KeyVaultKeyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicECHSM(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccKeyVaultKey_curveEC(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	r := KeyVaultKeyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.curveEC(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVaultKey_basicRSA(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	r := KeyVaultKeyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicRSA(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("key_size"),
	})
}

func TestAccKeyVaultKey_basicRSAHSM(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	r := KeyVaultKeyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicRSAHSM(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("key_size"),
	})
}

func TestAccKeyVaultKey_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	r := KeyVaultKeyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("not_before_date").HasValue("2020-01-01T01:02:03Z"),
				check.That(data.ResourceName).Key("expiration_date").HasValue("2021-01-01T01:02:03Z"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.hello").HasValue("world"),
				check.That(data.ResourceName).Key("versionless_id").HasValue(fmt.Sprintf("https://acctestkv-%s.vault.azure.net/keys/key-%s", data.RandomString, data.RandomString)),
			),
		},
		data.ImportStep("key_size"),
	})
}

func TestAccKeyVaultKey_softDeleteRecovery(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	r := KeyVaultKeyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.softDeleteRecovery(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("not_before_date").HasValue("2020-01-01T01:02:03Z"),
				check.That(data.ResourceName).Key("expiration_date").HasValue("2021-01-01T01:02:03Z"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.hello").HasValue("world"),
			),
		},
		data.ImportStep("key_size"),
		{
			Config:  r.softDeleteRecovery(data, false),
			Destroy: true,
		},
		{
			Config: r.softDeleteRecovery(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("not_before_date").HasValue("2020-01-01T01:02:03Z"),
				check.That(data.ResourceName).Key("expiration_date").HasValue("2021-01-01T01:02:03Z"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.hello").HasValue("world"),
			),
		},
	})
}

func TestAccKeyVaultKey_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	r := KeyVaultKeyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicRSA(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_opts.#").HasValue("6"),
				check.That(data.ResourceName).Key("key_opts.0").HasValue("decrypt"),
			),
		},
		{
			Config: r.basicUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_opts.#").HasValue("5"),
				check.That(data.ResourceName).Key("key_opts.0").HasValue("encrypt"),
			),
		},
	})
}

func TestAccKeyVaultKey_updatedExternally(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	r := KeyVaultKeyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicEC(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.updateExpiryDate("2029-02-02T12:59:00Z")),
			),
			ExpectNonEmptyPlan: true,
		},
		{
			Config: r.basicECUpdatedExternally(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:   r.basicECUpdatedExternally(data),
			PlanOnly: true,
		},
		data.ImportStep("key_size"),
	})
}

func TestAccKeyVaultKey_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	r := KeyVaultKeyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basicEC,
			TestResource: r,
		}),
	})
}

func TestAccKeyVaultKey_disappearsWhenParentKeyVaultDeleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	r := KeyVaultKeyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicEC(data),
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.destroyParentKeyVault, "azurerm_key_vault.test"),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccKeyVaultKey_withExternalAccessPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key", "test")
	r := KeyVaultKeyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withExternalAccessPolicy(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("key_size"),
		{
			Config: r.withExternalAccessPolicyUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("key_size"),
	})
}

func (r KeyVaultKeyResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
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

	resp, err := client.GetKey(ctx, id.KeyVaultBaseUrl, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving Key Vault Key %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.Key != nil), nil
}

func (KeyVaultKeyResource) destroyParentKeyVault(ctx context.Context, client *clients.Client, state *terraform.InstanceState) error {
	ok, err := KeyVaultResource{}.Destroy(ctx, client, state)
	if err != nil {
		return err
	}

	if ok == nil || !*ok {
		return fmt.Errorf("deleting parent key vault failed")
	}

	return nil
}

func (KeyVaultKeyResource) updateExpiryDate(expiryDate string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
		name := state.Attributes["name"]
		keyVaultId, err := parse.VaultID(state.Attributes["key_vault_id"])
		if err != nil {
			return err
		}

		vaultBaseUrl, err := clients.KeyVault.BaseUriForKeyVault(ctx, *keyVaultId)
		if err != nil {
			return fmt.Errorf("looking up base uri for Key %q from %q: %+v", name, keyVaultId, err)
		}

		expirationDate, err := time.Parse(time.RFC3339, expiryDate)
		if err != nil {
			return err
		}
		expirationUnixTime := date.UnixTime(expirationDate)
		update := keyvault.KeyUpdateParameters{
			KeyAttributes: &keyvault.KeyAttributes{
				Expires: &expirationUnixTime,
			},
		}
		if _, err = clients.KeyVault.ManagementClient.UpdateKey(ctx, *vaultBaseUrl, name, "", update); err != nil {
			return fmt.Errorf("updating secret: %+v", err)
		}

		return nil
	}
}

func (KeyVaultKeyResource) Destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	name := state.Attributes["name"]
	keyVaultId, err := parse.VaultID(state.Attributes["key_vault_id"])
	if err != nil {
		return nil, err
	}

	vaultBaseUrl, err := client.KeyVault.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return nil, fmt.Errorf("looking up Secret %q vault url from id %q: %+v", name, keyVaultId, err)
	}

	if _, err := client.KeyVault.ManagementClient.DeleteKey(ctx, *vaultBaseUrl, name); err != nil {
		return nil, fmt.Errorf("deleting keyVaultManagementClient: %+v", err)
	}

	return utils.Bool(true), nil
}

func (r KeyVaultKeyResource) basicEC(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "EC"
  key_size     = 2048

  key_opts = [
    "sign",
    "verify",
  ]
}
`, r.templateStandard(data), data.RandomString)
}

func (r KeyVaultKeyResource) basicECUpdatedExternally(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_key" "test" {
  name            = "key-%s"
  key_vault_id    = azurerm_key_vault.test.id
  key_type        = "EC"
  key_size        = 2048
  expiration_date = "2029-02-02T12:59:00Z"

  key_opts = [
    "sign",
    "verify",
  ]

  tags = {
    Rick = "Morty"
  }
}
`, r.templateStandard(data), data.RandomString)
}

func (r KeyVaultKeyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_key" "import" {
  name         = azurerm_key_vault_key.test.name
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "EC"
  key_size     = 2048

  key_opts = [
    "sign",
    "verify",
  ]
}
`, r.basicEC(data))
}

func (r KeyVaultKeyResource) basicRSA(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
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
`, r.templateStandard(data), data.RandomString)
}

func (r KeyVaultKeyResource) basicRSAHSM(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA-HSM"
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
`, r.templatePremium(data), data.RandomString)
}

func (r KeyVaultKeyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_key" "test" {
  name            = "key-%s"
  key_vault_id    = azurerm_key_vault.test.id
  key_type        = "RSA"
  key_size        = 2048
  not_before_date = "2020-01-01T01:02:03Z"
  expiration_date = "2021-01-01T01:02:03Z"

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  tags = {
    "hello" = "world"
  }
}
`, r.templateStandard(data), data.RandomString)
}

func (r KeyVaultKeyResource) basicUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}
`, r.templateStandard(data), data.RandomString)
}

func (r KeyVaultKeyResource) curveEC(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "EC"
  curve        = "P-521"

  key_opts = [
    "sign",
    "verify",
  ]
}
`, r.templateStandard(data), data.RandomString)
}

func (r KeyVaultKeyResource) basicECHSM(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "EC-HSM"
  curve        = "P-521"

  key_opts = [
    "sign",
    "verify",
  ]
}
`, r.templatePremium(data), data.RandomString)
}

func (r KeyVaultKeyResource) softDeleteRecovery(data acceptance.TestData, purge bool) string {
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

resource "azurerm_key_vault_key" "test" {
  name            = "key-%s"
  key_vault_id    = azurerm_key_vault.test.id
  key_type        = "RSA"
  key_size        = 2048
  not_before_date = "2020-01-01T01:02:03Z"
  expiration_date = "2021-01-01T01:02:03Z"

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  tags = {
    "hello" = "world"
  }
}
`, purge, r.templateStandard(data), data.RandomString)
}

func (KeyVaultKeyResource) withExternalAccessPolicy(data acceptance.TestData) string {
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
    environment = "accTest"
  }
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "create",
    "delete",
    "get",
    "purge",
    "recover",
    "update",
  ]

  secret_permissions = [
    "get",
    "delete",
    "set",
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "EC"
  key_size     = 2048

  key_opts = [
    "sign",
    "verify",
  ]

  depends_on = [azurerm_key_vault_access_policy.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (KeyVaultKeyResource) withExternalAccessPolicyUpdate(data acceptance.TestData) string {
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
    environment = "accTest"
  }
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "create",
    "delete",
    "encrypt",
    "get",
    "purge",
    "recover",
    "update",
  ]

  secret_permissions = [
    "get",
    "delete",
    "set",
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "EC"
  key_size     = 2048

  key_opts = [
    "sign",
    "verify",
  ]

  depends_on = [azurerm_key_vault_access_policy.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (r KeyVaultKeyResource) templateStandard(data acceptance.TestData) string {
	return r.template(data, "standard")
}

func (r KeyVaultKeyResource) templatePremium(data acceptance.TestData) string {
	return r.template(data, "premium")
}

func (KeyVaultKeyResource) template(data acceptance.TestData, sku string) string {
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
  sku_name                   = "%s"
  soft_delete_enabled        = true
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "create",
      "delete",
      "get",
      "purge",
      "recover",
      "update",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
    ]
  }

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, sku)
}
