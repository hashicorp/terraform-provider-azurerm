package keyvault_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type KeyVaultCertificateIssuerResource struct {
}

func TestAccKeyVaultCertificateIssuer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_issuer", "test")
	r := KeyVaultCertificateIssuerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccKeyVaultCertificateIssuer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_issuer", "test")
	r := KeyVaultCertificateIssuerResource{}

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

func TestAccKeyVaultCertificateIssuer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_issuer", "test")
	r := KeyVaultCertificateIssuerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccKeyVaultCertificateIssuer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_issuer", "test")
	r := KeyVaultCertificateIssuerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccKeyVaultCertificateIssuer_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_issuer", "test")
	r := KeyVaultCertificateIssuerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func (r KeyVaultCertificateIssuerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.KeyVault.ManagementClient
	keyVaultsClient := clients.KeyVault

	id, err := parse.IssuerID(state.ID)
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

	resp, err := client.GetCertificateIssuer(ctx, id.KeyVaultBaseUrl, id.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to make Read request on Azure KeyVault Certificate Issuer %s: %+v", id.Name, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r KeyVaultCertificateIssuerResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	dataPlaneClient := client.KeyVault.ManagementClient
	keyVaultsClient := client.KeyVault

	name := state.Attributes["name"]
	keyVaultId, err := parse.VaultID(state.Attributes["key_vault_id"])
	if err != nil {
		return nil, err
	}

	vaultBaseUrl, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return utils.Bool(false), fmt.Errorf("failed to look up base URI from id %q: %+v", keyVaultId, err)
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return nil, fmt.Errorf("failed to check if key vault %q for Certificate Issuer %q in Vault at url %q exists: %v", keyVaultId.ID(), name, *vaultBaseUrl, err)
	}
	if !ok {
		return utils.Bool(false), fmt.Errorf("Certificate Issuer %q Key Vault %q was not found in Key Vault at URI %q", name, keyVaultId.ID(), *vaultBaseUrl)
	}

	resp, err := dataPlaneClient.DeleteCertificateIssuer(ctx, *vaultBaseUrl, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(true), nil
		}

		return nil, fmt.Errorf("Bad: Delete on keyVaultManagementClient: %+v", err)
	}

	return utils.Bool(true), nil
}

func (KeyVaultCertificateIssuerResource) basic(data acceptance.TestData) string {
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
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "delete",
      "import",
      "get",
      "manageissuers",
      "setissuers",
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]
  }
}

resource "azurerm_key_vault_certificate_issuer" "test" {
  name          = "acctestKVCI-%d"
  key_vault_id  = azurerm_key_vault.test.id
  provider_name = "OneCertV2-PrivateCA"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r KeyVaultCertificateIssuerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_certificate_issuer" "import" {
  name          = azurerm_key_vault_certificate_issuer.test.name
  key_vault_id  = azurerm_key_vault_certificate_issuer.test.key_vault_id
  org_id        = azurerm_key_vault_certificate_issuer.test.org_id
  account_id    = azurerm_key_vault_certificate_issuer.test.account_id
  password      = "test"
  provider_name = azurerm_key_vault_certificate_issuer.test.provider_name
}

`, r.basic(data))
}

func (KeyVaultCertificateIssuerResource) complete(data acceptance.TestData) string {
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
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "delete",
      "import",
      "get",
      "manageissuers",
      "setissuers",
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]
  }
}

resource "azurerm_key_vault_certificate_issuer" "test" {
  name          = "acctestKVCI-%d"
  key_vault_id  = azurerm_key_vault.test.id
  account_id    = "test-account"
  password      = "test"
  provider_name = "DigiCert"

  org_id = "accTestOrg"
  admin {
    email_address = "admin@contoso.com"
    first_name    = "First"
    last_name     = "Last"
    phone         = "01234567890"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}
