package confidentialledger_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/sdk/2021-05-13-preview/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ConfidentialLedgerResource struct{}

func TestAccConfidentialLedger_public(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confidential_ledger", "test")
	r := ConfidentialLedgerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.public(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConfidentialLedger_private(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confidential_ledger", "test")
	r := ConfidentialLedgerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.private(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConfidentialLedger_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confidential_ledger", "test")
	r := ConfidentialLedgerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.public(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccConfidentialLedger_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confidential_ledger", "test")
	r := ConfidentialLedgerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConfidentialLedger_aadBasedServicePrincipals(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confidential_ledger", "test")
	r := ConfidentialLedgerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.aadBasedServicePrincipals(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("aad_based_security_principals.#").HasValue("3"),
				check.That(data.ResourceName).Key("aad_based_security_principals.0.ledger_role_name").HasValue("Administrator"),
				check.That(data.ResourceName).Key("aad_based_security_principals.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("aad_based_security_principals.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("aad_based_security_principals.1.ledger_role_name").HasValue("Contributor"),
				check.That(data.ResourceName).Key("aad_based_security_principals.1.principal_id").Exists(),
				check.That(data.ResourceName).Key("aad_based_security_principals.1.tenant_id").Exists(),
				check.That(data.ResourceName).Key("aad_based_security_principals.2.ledger_role_name").HasValue("Reader"),
				check.That(data.ResourceName).Key("aad_based_security_principals.2.principal_id").Exists(),
				check.That(data.ResourceName).Key("aad_based_security_principals.2.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConfidentialLedger_certBasedAdministrator(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confidential_ledger", "test")
	r := ConfidentialLedgerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.certBasedAdministrator(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cert_based_security_principals.#").HasValue("1"),
				check.That(data.ResourceName).Key("cert_based_security_principals.0.ledger_role_name").HasValue("Administrator"),
				check.That(data.ResourceName).Key("cert_based_security_principals.0.cert").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConfidentialLedger_certBasedContributor(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confidential_ledger", "test")
	r := ConfidentialLedgerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.certBasedContributor(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cert_based_security_principals.#").HasValue("1"),
				check.That(data.ResourceName).Key("cert_based_security_principals.0.ledger_role_name").HasValue("Contributor"),
				check.That(data.ResourceName).Key("cert_based_security_principals.0.cert").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConfidentialLedger_certBasedReader(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confidential_ledger", "test")
	r := ConfidentialLedgerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.certBasedReader(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cert_based_security_principals.#").HasValue("1"),
				check.That(data.ResourceName).Key("cert_based_security_principals.0.ledger_role_name").HasValue("Reader"),
				check.That(data.ResourceName).Key("cert_based_security_principals.0.cert").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConfidentialLedger_combinedServicePrincipals(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confidential_ledger", "test")
	r := ConfidentialLedgerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.combinedServicePrincipals(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("aad_based_security_principals.#").HasValue("3"),
				check.That(data.ResourceName).Key("aad_based_security_principals.0.ledger_role_name").HasValue("Administrator"),
				check.That(data.ResourceName).Key("aad_based_security_principals.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("aad_based_security_principals.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("aad_based_security_principals.1.ledger_role_name").HasValue("Contributor"),
				check.That(data.ResourceName).Key("aad_based_security_principals.1.principal_id").Exists(),
				check.That(data.ResourceName).Key("aad_based_security_principals.1.tenant_id").Exists(),
				check.That(data.ResourceName).Key("aad_based_security_principals.2.ledger_role_name").HasValue("Reader"),
				check.That(data.ResourceName).Key("aad_based_security_principals.2.principal_id").Exists(),
				check.That(data.ResourceName).Key("aad_based_security_principals.2.tenant_id").Exists(),
				check.That(data.ResourceName).Key("cert_based_security_principals.#").HasValue("1"),
				check.That(data.ResourceName).Key("cert_based_security_principals.0.ledger_role_name").HasValue("Administrator"),
				check.That(data.ResourceName).Key("cert_based_security_principals.0.cert").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConfidentialLedger_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confidential_ledger", "test")
	r := ConfidentialLedgerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.public(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.publicUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (t ConfidentialLedgerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := confidentialledger.ParseLedgerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ConfidentialLedger.ConfidentialLedgereClient.LedgerGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (ConfidentialLedgerResource) public(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acl-terraform-acc-%d"
  location = "%s"
}

resource "azurerm_confidential_ledger" "test" {
  name                = "terraform-acc-%d"
  ledger_type         = "Public"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ConfidentialLedgerResource) private(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acl-terraform-acc-%d"
  location = "%s"
}

resource "azurerm_confidential_ledger" "test" {
  name                = "terraform-acc-%d"
  ledger_type         = "Private"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ConfidentialLedgerResource) requiresImport(data acceptance.TestData) string {
	template := r.public(data)
	return fmt.Sprintf(`
%s

resource "azurerm_confidential_ledger" "import" {
  name                = azurerm_confidential_ledger.test.name
  ledger_type         = azurerm_confidential_ledger.test.ledger_type
  location            = azurerm_confidential_ledger.test.location
  resource_group_name = azurerm_confidential_ledger.test.resource_group_name
}
`, template)
}

func (ConfidentialLedgerResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acl-terraform-acc-%d"
  location = "%s"
}

resource "azurerm_confidential_ledger" "test" {
  name                = "terraform-acc-%d"
  ledger_type         = "Private"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "development"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ConfidentialLedgerResource) aadBasedServicePrincipals(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acl-terraform-acc-%d"
  location = "%s"
}

resource "azurerm_confidential_ledger" "test" {
  aad_based_security_principals {
    principal_id     = "34621747-6fc8-4771-a2eb-72f31c461f20"
    tenant_id        = "bce123b9-2b7b-4975-8360-5ca0b9b1cd00"
    ledger_role_name = "Administrator"
  }
  aad_based_security_principals {
    principal_id     = "34621747-6fc8-4771-a2eb-72f31c461f21"
    tenant_id        = "bce123b9-2b7b-4975-8360-5ca0b9b1cd01"
    ledger_role_name = "Contributor"
  }
  aad_based_security_principals {
    principal_id     = "34621747-6fc8-4771-a2eb-72f31c461f22"
    tenant_id        = "bce123b9-2b7b-4975-8360-5ca0b9b1cd02"
    ledger_role_name = "Reader"
  }

  name                = "terraform-acc-%d"
  ledger_type         = "Public"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "development"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ConfidentialLedgerResource) certBasedAdministrator(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acl-terraform-acc-%d"
  location = "%s"
}

resource "azurerm_confidential_ledger" "test" {
  cert_based_security_principals {
    cert             = "-----BEGIN CERTIFICATE-----MIIBsjCCATigAwIBAgIUZWIbyG79TniQLd2UxJuU74tqrKcwCgYIKoZIzj0EAwMwEDEOMAwGA1UEAwwFdXNlcjAwHhcNMjEwMzE2MTgwNjExWhcNMjIwMzE2MTgwNjExWjAQMQ4wDAYDVQQDDAV1c2VyMDB2MBAGByqGSM49AgEGBSuBBAAiA2IABBiWSo/j8EFit7aUMm5lF+lUmCu+IgfnpFD+7QMgLKtxRJ3aGSqgS/GpqcYVGddnODtSarNE/HyGKUFUolLPQ5ybHcouUk0kyfA7XMeSoUA4lBz63Wha8wmXo+NdBRo39qNTMFEwHQYDVR0OBBYEFPtuhrwgGjDFHeUUT4nGsXaZn69KMB8GA1UdIwQYMBaAFPtuhrwgGjDFHeUUT4nGsXaZn69KMA8GA1UdEwEB/wQFMAMBAf8wCgYIKoZIzj0EAwMDaAAwZQIxAOnozm2CyqRwSSQLls5r+mUHRGRyXHXwYtM4Dcst/VEZdmS9fqvHRCHbjUlO/+HNfgIwMWZ4FmsjD3wnPxONOm9YdVn/PRD7SsPRPbOjwBiE4EBGaHDsLjYAGDSGi7NJnSkA-----END CERTIFICATE-----"
    ledger_role_name = "Administrator"
  }

  name                = "terraform-acc-%d"
  ledger_type         = "Public"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "development"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ConfidentialLedgerResource) certBasedContributor(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acl-terraform-acc-%d"
  location = "%s"
}

resource "azurerm_confidential_ledger" "test" {
  cert_based_security_principals {
    cert             = "-----BEGIN CERTIFICATE-----MIIBsjCCATigAwIBAgIUZWIbyG79TniQLd2UxJuU74tqrKcwCgYIKoZIzj0EAwMwEDEOMAwGA1UEAwwFdXNlcjAwHhcNMjEwMzE2MTgwNjExWhcNMjIwMzE2MTgwNjExWjAQMQ4wDAYDVQQDDAV1c2VyMDB2MBAGByqGSM49AgEGBSuBBAAiA2IABBiWSo/j8EFit7aUMm5lF+lUmCu+IgfnpFD+7QMgLKtxRJ3aGSqgS/GpqcYVGddnODtSarNE/HyGKUFUolLPQ5ybHcouUk0kyfA7XMeSoUA4lBz63Wha8wmXo+NdBRo39qNTMFEwHQYDVR0OBBYEFPtuhrwgGjDFHeUUT4nGsXaZn69KMB8GA1UdIwQYMBaAFPtuhrwgGjDFHeUUT4nGsXaZn69KMA8GA1UdEwEB/wQFMAMBAf8wCgYIKoZIzj0EAwMDaAAwZQIxAOnozm2CyqRwSSQLls5r+mUHRGRyXHXwYtM4Dcst/VEZdmS9fqvHRCHbjUlO/+HNfgIwMWZ4FmsjD3wnPxONOm9YdVn/PRD7SsPRPbOjwBiE4EBGaHDsLjYAGDSGi7NJnSkA-----END CERTIFICATE-----"
    ledger_role_name = "Contributor"
  }

  name                = "terraform-acc-%d"
  ledger_type         = "Public"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "development"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ConfidentialLedgerResource) certBasedReader(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acl-terraform-acc-%d"
  location = "%s"
}

resource "azurerm_confidential_ledger" "test" {
  cert_based_security_principals {
    cert             = "-----BEGIN CERTIFICATE-----MIIBsjCCATigAwIBAgIUZWIbyG79TniQLd2UxJuU74tqrKcwCgYIKoZIzj0EAwMwEDEOMAwGA1UEAwwFdXNlcjAwHhcNMjEwMzE2MTgwNjExWhcNMjIwMzE2MTgwNjExWjAQMQ4wDAYDVQQDDAV1c2VyMDB2MBAGByqGSM49AgEGBSuBBAAiA2IABBiWSo/j8EFit7aUMm5lF+lUmCu+IgfnpFD+7QMgLKtxRJ3aGSqgS/GpqcYVGddnODtSarNE/HyGKUFUolLPQ5ybHcouUk0kyfA7XMeSoUA4lBz63Wha8wmXo+NdBRo39qNTMFEwHQYDVR0OBBYEFPtuhrwgGjDFHeUUT4nGsXaZn69KMB8GA1UdIwQYMBaAFPtuhrwgGjDFHeUUT4nGsXaZn69KMA8GA1UdEwEB/wQFMAMBAf8wCgYIKoZIzj0EAwMDaAAwZQIxAOnozm2CyqRwSSQLls5r+mUHRGRyXHXwYtM4Dcst/VEZdmS9fqvHRCHbjUlO/+HNfgIwMWZ4FmsjD3wnPxONOm9YdVn/PRD7SsPRPbOjwBiE4EBGaHDsLjYAGDSGi7NJnSkA-----END CERTIFICATE-----"
    ledger_role_name = "Reader"
  }

  name                = "terraform-acc-%d"
  ledger_type         = "Public"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "development"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ConfidentialLedgerResource) combinedServicePrincipals(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acl-terraform-acc-%d"
  location = "%s"
}

resource "azurerm_confidential_ledger" "test" {
  aad_based_security_principals {
    principal_id     = "34621747-6fc8-4771-a2eb-72f31c461f20"
    tenant_id        = "bce123b9-2b7b-4975-8360-5ca0b9b1cd00"
    ledger_role_name = "Administrator"
  }
  aad_based_security_principals {
    principal_id     = "34621747-6fc8-4771-a2eb-72f31c461f21"
    tenant_id        = "bce123b9-2b7b-4975-8360-5ca0b9b1cd01"
    ledger_role_name = "Contributor"
  }
  aad_based_security_principals {
    principal_id     = "34621747-6fc8-4771-a2eb-72f31c461f22"
    tenant_id        = "bce123b9-2b7b-4975-8360-5ca0b9b1cd02"
    ledger_role_name = "Reader"
  }

  cert_based_security_principals {
    cert             = "-----BEGIN CERTIFICATE-----MIIBsjCCATigAwIBAgIUZWIbyG79TniQLd2UxJuU74tqrKcwCgYIKoZIzj0EAwMwEDEOMAwGA1UEAwwFdXNlcjAwHhcNMjEwMzE2MTgwNjExWhcNMjIwMzE2MTgwNjExWjAQMQ4wDAYDVQQDDAV1c2VyMDB2MBAGByqGSM49AgEGBSuBBAAiA2IABBiWSo/j8EFit7aUMm5lF+lUmCu+IgfnpFD+7QMgLKtxRJ3aGSqgS/GpqcYVGddnODtSarNE/HyGKUFUolLPQ5ybHcouUk0kyfA7XMeSoUA4lBz63Wha8wmXo+NdBRo39qNTMFEwHQYDVR0OBBYEFPtuhrwgGjDFHeUUT4nGsXaZn69KMB8GA1UdIwQYMBaAFPtuhrwgGjDFHeUUT4nGsXaZn69KMA8GA1UdEwEB/wQFMAMBAf8wCgYIKoZIzj0EAwMDaAAwZQIxAOnozm2CyqRwSSQLls5r+mUHRGRyXHXwYtM4Dcst/VEZdmS9fqvHRCHbjUlO/+HNfgIwMWZ4FmsjD3wnPxONOm9YdVn/PRD7SsPRPbOjwBiE4EBGaHDsLjYAGDSGi7NJnSkA-----END CERTIFICATE-----"
    ledger_role_name = "Administrator"
  }

  name                = "terraform-acc-%d"
  ledger_type         = "Public"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "development"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ConfidentialLedgerResource) publicUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acl-terraform-acc-%d"
  location = "%s"
}

resource "azurerm_confidential_ledger" "test" {
  name                = "terraform-acc-%d"
  ledger_type         = "Public"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    Updated = "Yes"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
