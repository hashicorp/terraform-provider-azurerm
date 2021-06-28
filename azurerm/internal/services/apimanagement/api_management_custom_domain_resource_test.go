package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ApiManagementCustomDomainResource struct {
}

func TestAccApiManagementCustomDomain_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_custom_domain", "test")
	r := ApiManagementCustomDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementCustomDomain_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_custom_domain", "test")
	r := ApiManagementCustomDomainResource{}

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

func TestAccApiManagementCustomDomain_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_custom_domain", "test")
	r := ApiManagementCustomDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.proxyOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.developerPortalOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (ApiManagementCustomDomainResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName

	resp, err := clients.ApiManagement.ServiceClient.Get(ctx, resourceGroup, serviceName)
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagement Custom Domain (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r ApiManagementCustomDomainResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_custom_domain" "test" {
  api_management_id = azurerm_api_management.test.id

  proxy {
    host_name    = "api.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }

  developer_portal {
    host_name    = "portal.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }
}
`, r.template(data))
}

func (r ApiManagementCustomDomainResource) proxyOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_custom_domain" "test" {
  api_management_id = azurerm_api_management.test.id

  proxy {
    host_name    = "api.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }
}
`, r.template(data))
}

func (r ApiManagementCustomDomainResource) developerPortalOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_custom_domain" "test" {
  api_management_id = azurerm_api_management.test.id

  developer_portal {
    host_name    = "portal.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }
}
`, r.template(data))
}

func (r ApiManagementCustomDomainResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_custom_domain" "import" {
  api_management_id = azurerm_api_management_custom_domain.test.api_management_id

  proxy {
    host_name    = "api.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }

  developer_portal {
    host_name    = "portal.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }
}
`, r.basic(data))
}

func (ApiManagementCustomDomainResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault" "test" {
  name                = "apimkv%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "create",
      "delete",
      "get",
      "update",
      "purge",
    ]

    key_permissions = [
      "create",
      "get",
    ]

    secret_permissions = [
      "get",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_api_management.test.identity.0.principal_id

    certificate_permissions = [
      "get",
    ]

    secret_permissions = [
      "get",
    ]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%[3]s"
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
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=api.example.com"
      validity_in_months = 12

      subject_alternative_names {
        dns_names = [
          "api.example.com",
          "portal.example.com",
        ]
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
