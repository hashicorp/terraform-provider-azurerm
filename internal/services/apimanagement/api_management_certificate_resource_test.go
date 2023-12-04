// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementCertificateResource struct{}

func TestAccApiManagementCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_certificate", "test")
	r := ApiManagementCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("expiration").Exists(),
				check.That(data.ResourceName).Key("subject").Exists(),
				check.That(data.ResourceName).Key("thumbprint").Exists(),
			),
		},
		{
			ResourceName:      data.ResourceName,
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				// not returned from the API
				"data",
				"password",
			},
		},
	})
}

func TestAccApiManagementCertificate_basicKeyVaultSystemIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_certificate", "test")
	r := ApiManagementCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicKeyVaultSystemIdentity(data, "cert1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_vault_secret_id").Exists(),
				check.That(data.ResourceName).Key("expiration").Exists(),
				check.That(data.ResourceName).Key("subject").Exists(),
				check.That(data.ResourceName).Key("thumbprint").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementCertificate_basicKeyVaultUserIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_certificate", "test")
	r := ApiManagementCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicKeyVaultUserIdentity(data, "cert1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_vault_secret_id").Exists(),
				check.That(data.ResourceName).Key("key_vault_identity_client_id").Exists(),
				check.That(data.ResourceName).Key("expiration").Exists(),
				check.That(data.ResourceName).Key("subject").Exists(),
				check.That(data.ResourceName).Key("thumbprint").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementCertificate_basicKeyVaultUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_certificate", "test")
	r := ApiManagementCertificateResource{}

	certUpdatedRegex := regexp.MustCompile(fmt.Sprintf(`https://acct%d\.vault\.azure\.net/secrets/cert2/[a-z0-9]{32}`, data.RandomInteger))

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicKeyVaultSystemIdentity(data, "cert1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_vault_secret_id").Exists(),
				check.That(data.ResourceName).Key("expiration").Exists(),
				check.That(data.ResourceName).Key("subject").Exists(),
				check.That(data.ResourceName).Key("thumbprint").Exists(),
			),
		},
		{
			Config: r.basicKeyVaultSystemIdentity(data, "cert2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_vault_secret_id").MatchesRegex(certUpdatedRegex),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementCertificate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_certificate", "test")
	r := ApiManagementCertificateResource{}

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

func (ApiManagementCertificateResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := certificate.ParseCertificateID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.CertificatesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (ApiManagementCertificateResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}

resource "azurerm_api_management_certificate" "test" {
  name                = "example-cert"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  data                = filebase64("testdata/keyvaultcert.pfx")
  password            = ""
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ApiManagementCertificateResource) basicKeyVaultSystemIdentity(data acceptance.TestData, certificate string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_api_management.test.identity.0.tenant_id
  object_id    = azurerm_api_management.test.identity.0.principal_id

  secret_permissions = [
    "Get",
  ]

  certificate_permissions = [
    "Get",
  ]
}

resource "azurerm_api_management_certificate" "test" {
  name                = "example-cert"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name

  key_vault_secret_id = azurerm_key_vault_certificate.%s.secret_id
}
`, r.templateKeyVault(data), data.RandomInteger, certificate)
}

func (r ApiManagementCertificateResource) basicKeyVaultUserIdentity(data acceptance.TestData, certificate string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.test.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id

  secret_permissions = [
    "Get",
  ]

  certificate_permissions = [
    "Get",
  ]
}

resource "azurerm_api_management_certificate" "test" {
  name                = "example-cert"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name

  key_vault_secret_id          = azurerm_key_vault_certificate.%s.secret_id
  key_vault_identity_client_id = azurerm_user_assigned_identity.test.client_id
}
`, r.templateKeyVault(data), data.RandomInteger, data.RandomInteger, certificate)
}

func (ApiManagementCertificateResource) templateKeyVault(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "test" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acct%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tenant_id = data.azurerm_client_config.test.tenant_id

  sku_name = "standard"
}

resource "azurerm_key_vault_access_policy" "sptest" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.test.tenant_id
  object_id    = data.azurerm_client_config.test.object_id

  secret_permissions = [
    "Delete",
    "Get",
    "Purge",
    "Set",
  ]

  certificate_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Import",
  ]
}

resource "azurerm_key_vault_certificate" "cert1" {
  name         = "cert1"
  key_vault_id = azurerm_key_vault.test.id

  depends_on = [azurerm_key_vault_access_policy.sptest]

  certificate {
    contents = filebase64("testdata/api_management_api_test.pfx")
    password = "terraform"
  }

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }
  }
}

resource "azurerm_key_vault_certificate" "cert2" {
  name         = "cert2"
  key_vault_id = azurerm_key_vault.test.id

  depends_on = [azurerm_key_vault_access_policy.sptest]

  certificate {
    contents = filebase64("testdata/api_management_api2_test.pfx")
    password = "terraform"
  }

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ApiManagementCertificateResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_certificate" "import" {
  name                = azurerm_api_management_certificate.test.name
  api_management_name = azurerm_api_management_certificate.test.api_management_name
  resource_group_name = azurerm_api_management_certificate.test.resource_group_name
  data                = azurerm_api_management_certificate.test.data
  password            = azurerm_api_management_certificate.test.password
}
`, r.basic(data))
}
