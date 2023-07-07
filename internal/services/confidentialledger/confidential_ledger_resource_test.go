// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package confidentialledger_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/confidentialledger/2022-05-13/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
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
		{
			Config: r.publicUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
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
		{
			Config: r.privateUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
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

func TestAccConfidentialLedger_certBased(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confidential_ledger", "test")
	r := ConfidentialLedgerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// add it
			Config: r.certBased(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// update it
			Config: r.certBasedUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// remove it
			Config: r.private(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ConfidentialLedgerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := confidentialledger.ParseLedgerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ConfidentialLedger.ConfidentialLedgerClient.LedgerGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r ConfidentialLedgerResource) public(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_confidential_ledger" "test" {
  name                = "acctest-tfci-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ledger_type         = "Public"

  azuread_based_service_principal {
    ledger_role_name = "Administrator"
    principal_id     = data.azurerm_client_config.current.object_id
    tenant_id        = data.azurerm_client_config.current.tenant_id
  }
}
`, template, data.RandomInteger)
}

func (r ConfidentialLedgerResource) publicUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_confidential_ledger" "test" {
  name                = "acctest-tfci-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ledger_type         = "Public"

  azuread_based_service_principal {
    ledger_role_name = "Administrator"
    principal_id     = data.azurerm_client_config.current.object_id
    tenant_id        = data.azurerm_client_config.current.tenant_id
  }

  azuread_based_service_principal {
    ledger_role_name = "Reader"
    principal_id     = azurerm_user_assigned_identity.first.principal_id
    tenant_id        = azurerm_user_assigned_identity.first.tenant_id
  }

  azuread_based_service_principal {
    ledger_role_name = "Reader"
    principal_id     = azurerm_user_assigned_identity.second.principal_id
    tenant_id        = azurerm_user_assigned_identity.second.tenant_id
  }

  tags = {
    Environment = "Testing"
  }
}
`, template, data.RandomInteger)
}

func (r ConfidentialLedgerResource) private(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_confidential_ledger" "test" {
  name                = "acctest-tfci-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ledger_type         = "Private"

  azuread_based_service_principal {
    ledger_role_name = "Administrator"
    principal_id     = data.azurerm_client_config.current.object_id
    tenant_id        = data.azurerm_client_config.current.tenant_id
  }
}
`, template, data.RandomInteger)
}

func (r ConfidentialLedgerResource) privateUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_confidential_ledger" "test" {
  name                = "acctest-tfci-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  ledger_type         = "Private"

  azuread_based_service_principal {
    ledger_role_name = "Administrator"
    principal_id     = data.azurerm_client_config.current.object_id
    tenant_id        = data.azurerm_client_config.current.tenant_id
  }

  azuread_based_service_principal {
    ledger_role_name = "Reader"
    principal_id     = azurerm_user_assigned_identity.first.principal_id
    tenant_id        = azurerm_user_assigned_identity.first.tenant_id
  }

  azuread_based_service_principal {
    ledger_role_name = "Reader"
    principal_id     = azurerm_user_assigned_identity.second.principal_id
    tenant_id        = azurerm_user_assigned_identity.second.tenant_id
  }

  tags = {
    Environment = "Testing"
  }
}
`, template, data.RandomInteger)
}

func (r ConfidentialLedgerResource) certBased(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_confidential_ledger" "test" {
  name                = "acctest-tfci-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ledger_type         = "Private"

  azuread_based_service_principal {
    ledger_role_name = "Administrator"
    principal_id     = data.azurerm_client_config.current.object_id
    tenant_id        = data.azurerm_client_config.current.tenant_id
  }

  certificate_based_security_principal {
    ledger_role_name = "Administrator"
    pem_public_key   = "-----BEGIN CERTIFICATE-----MIIBsjCCATigAwIBAgIUZWIbyG79TniQLd2UxJuU74tqrKcwCgYIKoZIzj0EAwMwEDEOMAwGA1UEAwwFdXNlcjAwHhcNMjEwMzE2MTgwNjExWhcNMjIwMzE2MTgwNjExWjAQMQ4wDAYDVQQDDAV1c2VyMDB2MBAGByqGSM49AgEGBSuBBAAiA2IABBiWSo/j8EFit7aUMm5lF+lUmCu+IgfnpFD+7QMgLKtxRJ3aGSqgS/GpqcYVGddnODtSarNE/HyGKUFUolLPQ5ybHcouUk0kyfA7XMeSoUA4lBz63Wha8wmXo+NdBRo39qNTMFEwHQYDVR0OBBYEFPtuhrwgGjDFHeUUT4nGsXaZn69KMB8GA1UdIwQYMBaAFPtuhrwgGjDFHeUUT4nGsXaZn69KMA8GA1UdEwEB/wQFMAMBAf8wCgYIKoZIzj0EAwMDaAAwZQIxAOnozm2CyqRwSSQLls5r+mUHRGRyXHXwYtM4Dcst/VEZdmS9fqvHRCHbjUlO/+HNfgIwMWZ4FmsjD3wnPxONOm9YdVn/PRD7SsPRPbOjwBiE4EBGaHDsLjYAGDSGi7NJnSkA-----END CERTIFICATE-----"
  }
}
`, template, data.RandomInteger)
}

func (r ConfidentialLedgerResource) certBasedUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_confidential_ledger" "test" {
  name                = "acctest-tfci-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ledger_type         = "Public"

  azuread_based_service_principal {
    ledger_role_name = "Administrator"
    principal_id     = data.azurerm_client_config.current.object_id
    tenant_id        = data.azurerm_client_config.current.tenant_id
  }

  certificate_based_security_principal {
    ledger_role_name = "Reader"
    pem_public_key   = "-----BEGIN CERTIFICATE-----MIIBsjCCATigAwIBAgIUZWIbyG79TniQLd2UxJuU74tqrKcwCgYIKoZIzj0EAwMwEDEOMAwGA1UEAwwFdXNlcjAwHhcNMjEwMzE2MTgwNjExWhcNMjIwMzE2MTgwNjExWjAQMQ4wDAYDVQQDDAV1c2VyMDB2MBAGByqGSM49AgEGBSuBBAAiA2IABBiWSo/j8EFit7aUMm5lF+lUmCu+IgfnpFD+7QMgLKtxRJ3aGSqgS/GpqcYVGddnODtSarNE/HyGKUFUolLPQ5ybHcouUk0kyfA7XMeSoUA4lBz63Wha8wmXo+NdBRo39qNTMFEwHQYDVR0OBBYEFPtuhrwgGjDFHeUUT4nGsXaZn69KMB8GA1UdIwQYMBaAFPtuhrwgGjDFHeUUT4nGsXaZn69KMA8GA1UdEwEB/wQFMAMBAf8wCgYIKoZIzj0EAwMDaAAwZQIxAOnozm2CyqRwSSQLls5r+mUHRGRyXHXwYtM4Dcst/VEZdmS9fqvHRCHbjUlO/+HNfgIwMWZ4FmsjD3wnPxONOm9YdVn/PRD7SsPRPbOjwBiE4EBGaHDsLjYAGDSGi7NJnSkA-----END CERTIFICATE-----"
  }

  tags = {
    Environment = "Testing"
  }
}
`, template, data.RandomInteger)
}

func (r ConfidentialLedgerResource) requiresImport(data acceptance.TestData) string {
	template := r.public(data)
	return fmt.Sprintf(`
%s

resource "azurerm_confidential_ledger" "import" {
  name                = azurerm_confidential_ledger.test.name
  resource_group_name = azurerm_confidential_ledger.test.resource_group_name
  location            = azurerm_confidential_ledger.test.location
  ledger_type         = azurerm_confidential_ledger.test.ledger_type

  azuread_based_service_principal {
    ledger_role_name = azurerm_confidential_ledger.test.azuread_based_service_principal.0.ledger_role_name
    principal_id     = azurerm_confidential_ledger.test.azuread_based_service_principal.0.principal_id
    tenant_id        = azurerm_confidential_ledger.test.azuread_based_service_principal.0.tenant_id
  }
}
`, template)
}

func (ConfidentialLedgerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-ledger-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "first" {
  name                = "acctest-uai1-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_user_assigned_identity" "second" {
  name                = "acctest-uai2-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary)
}
