// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementWorkspaceCertificateResource struct{}

func TestAccApiManagementWorkspaceCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_certificate", "test")
	r := ApiManagementWorkspaceCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("certificate_data_base64"),
	})
}

func TestAccApiManagementWorkspaceCertificate_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_certificate", "test")
	r := ApiManagementWorkspaceCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("certificate_data_base64", "password"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("certificate_data_base64", "password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("certificate_data_base64"),
		{
			Config: r.keyVaultWithIdentity(data),
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
		data.ImportStep("certificate_data_base64", "password"),
	})
}

func TestAccApiManagementWorkspaceCertificate_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_certificate", "test")
	r := ApiManagementWorkspaceCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("certificate_data_base64", "password"),
	})
}

func TestAccApiManagementWorkspaceCertificate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_certificate", "test")
	r := ApiManagementWorkspaceCertificateResource{}

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

func TestAccApiManagementWorkspaceCertificate_keyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_certificate", "test")
	r := ApiManagementWorkspaceCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementWorkspaceCertificate_keyVaultWithIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_certificate", "test")
	r := ApiManagementWorkspaceCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultWithIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultWithIdentityUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultWithIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ApiManagementWorkspaceCertificateResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := certificate.ParseWorkspaceCertificateID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.CertificateClient_v2024_05_01.WorkspaceCertificateGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ApiManagementWorkspaceCertificateResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_certificate" "test" {
  name                        = "acctest-cert-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  certificate_data_base64     = filebase64("testdata/testacc_nopassword.pfx")
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceCertificateResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_certificate" "test" {
  name                        = "acctest-cert-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  certificate_data_base64     = filebase64("testdata/testacc_nopassword.pfx")
  password                    = ""
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceCertificateResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_certificate" "test" {
  name                        = "acctest-cert-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  certificate_data_base64     = filebase64("testdata/testacc.pfx")
  password                    = "TestAcc"
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceCertificateResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_workspace_certificate" "import" {
  name                        = azurerm_api_management_workspace_certificate.test.name
  api_management_workspace_id = azurerm_api_management_workspace_certificate.test.api_management_workspace_id
  certificate_data_base64     = azurerm_api_management_workspace_certificate.test.certificate_data_base64
}
`, r.basic(data))
}

func (r ApiManagementWorkspaceCertificateResource) keyVault(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
}

%s

%s

resource "azurerm_api_management_workspace_certificate" "test" {
  name                        = "acctest-cert-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  key_vault_secret_id         = azurerm_key_vault_certificate.test.secret_id
}
`, r.templateWithIdentity(data), r.keyVaultTemplate(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceCertificateResource) keyVaultUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
}

%s

%s

resource "azurerm_api_management_workspace_certificate" "test" {
  name                        = "acctest-cert-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  key_vault_secret_id         = azurerm_key_vault_certificate.test2.secret_id
}
`, r.templateWithIdentity(data), r.keyVaultTemplate(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceCertificateResource) keyVaultWithIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
}

%s

%s

resource "azurerm_api_management_workspace_certificate" "test" {
  name                             = "acctest-cert-%d"
  api_management_workspace_id      = azurerm_api_management_workspace.test.id
  key_vault_secret_id              = azurerm_key_vault_certificate.test.secret_id
  user_assigned_identity_client_id = azurerm_user_assigned_identity.test.client_id
}
`, r.templateWithUserAssignedIdentity(data), r.keyVaultTemplateWithUserAssignedIdentity(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceCertificateResource) keyVaultWithIdentityUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
}

%s

%s

resource "azurerm_api_management_workspace_certificate" "test" {
  name                             = "acctest-cert-%d"
  api_management_workspace_id      = azurerm_api_management_workspace.test.id
  key_vault_secret_id              = azurerm_key_vault_certificate.test2.secret_id
  user_assigned_identity_client_id = azurerm_user_assigned_identity.test2.client_id
}
`, r.templateWithUserAssignedIdentity(data), r.keyVaultTemplateWithUserAssignedIdentity(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceCertificateResource) keyVaultTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {
}

resource "azurerm_key_vault" "test" {
  name                = "acct%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Get",
      "Set",
      "Delete",
      "Purge",
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Import",
    ]
  }

  access_policy {
    tenant_id = azurerm_api_management.test.identity.0.tenant_id
    object_id = azurerm_api_management.test.identity.0.principal_id

    secret_permissions = [
      "Get",
      "Set",
      "Delete",
      "Purge",
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Import",
    ]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "cert1"
  key_vault_id = azurerm_key_vault.test.id

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

resource "azurerm_key_vault" "test2" {
  name                = "acct2%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Get",
      "Set",
      "Delete",
      "Purge",
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Import",
    ]
  }

  access_policy {
    tenant_id = azurerm_api_management.test.identity.0.tenant_id
    object_id = azurerm_api_management.test.identity.0.principal_id

    secret_permissions = [
      "Get",
      "Set",
      "Delete",
      "Purge",
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Import",
    ]
  }
}

resource "azurerm_key_vault_certificate" "test2" {
  name         = "cert2"
  key_vault_id = azurerm_key_vault.test2.id

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
`, data.RandomInteger)
}

func (r ApiManagementWorkspaceCertificateResource) keyVaultTemplateWithUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {
}

resource "azurerm_key_vault" "test" {
  name                = "acct%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Get",
      "Set",
      "Delete",
      "Purge",
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Import",
    ]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.test.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    secret_permissions = [
      "Get",
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Import",
    ]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "cert1"
  key_vault_id = azurerm_key_vault.test.id

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

resource "azurerm_key_vault" "test2" {
  name                = "acct2%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Get",
      "Set",
      "Delete",
      "Purge",
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Import",
    ]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.test2.tenant_id
    object_id = azurerm_user_assigned_identity.test2.principal_id

    secret_permissions = [
      "Get",
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Import",
    ]
  }
}

resource "azurerm_key_vault_certificate" "test2" {
  name         = "cert2"
  key_vault_id = azurerm_key_vault.test2.id

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
`, data.RandomInteger)
}

func (r ApiManagementWorkspaceCertificateResource) templateWithIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apim-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestapim-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_api_management_workspace" "test" {
  name              = "acctestapimws-%d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test Workspace"
  description       = "Test workspace for certificate"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementWorkspaceCertificateResource) templateWithUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apim-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctestUAI2-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_api_management" "test" {
  name                = "acctestapim-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
      azurerm_user_assigned_identity.test2.id,
    ]
  }
}

resource "azurerm_api_management_workspace" "test" {
  name              = "acctestapimws-%[1]d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test Workspace"
  description       = "Test workspace description"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (ApiManagementWorkspaceCertificateResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apim-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestapim-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"
}

resource "azurerm_api_management_workspace" "test" {
  name              = "acctestapimws-%d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test Workspace"
  description       = "Test workspace for certificate"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
