package web_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ServiceCustomHostnameBindingResource struct {
}

func TestAccAppServiceCustomHostnameBinding(t *testing.T) {
	appServiceEnvVariable := "ARM_TEST_APP_SERVICE"
	appServiceEnv := os.Getenv(appServiceEnvVariable)
	if appServiceEnv == "" {
		t.Skipf("Skipping as %q is not specified", appServiceEnvVariable)
	}

	domainEnvVariable := "ARM_TEST_DOMAIN"
	domainEnv := os.Getenv(domainEnvVariable)
	if domainEnv == "" {
		t.Skipf("Skipping as %q is not specified", domainEnvVariable)
	}

	// NOTE: this is a combined test rather than separate split out tests due to
	// the app service name being shared (so the tests don't conflict with each other)
	testCases := map[string]map[string]func(t *testing.T, appServiceEnv, domainEnv string){
		"basic": {
			"basic":          testAccAppServiceCustomHostnameBinding_basic,
			"multiple":       testAccAppServiceCustomHostnameBinding_multiple,
			"requiresImport": testAccAppServiceCustomHostnameBinding_requiresImport,
			"ssl":            testAccAppServiceCustomHostnameBinding_ssl,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t, appServiceEnv, domainEnv)
				})
			}
		})
	}
}

func testAccAppServiceCustomHostnameBinding_basic(t *testing.T, appServiceEnv, domainEnv string) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_custom_hostname_binding", "test")
	r := ServiceCustomHostnameBindingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data, appServiceEnv, domainEnv),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccAppServiceCustomHostnameBinding_requiresImport(t *testing.T, appServiceEnv, domainEnv string) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_custom_hostname_binding", "test")
	r := ServiceCustomHostnameBindingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data, appServiceEnv, domainEnv),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImport(data, appServiceEnv, domainEnv)
		}),
	})
}

func testAccAppServiceCustomHostnameBinding_multiple(t *testing.T, appServiceEnv, domainEnv string) {
	altDomainEnvVariable := "ARM_ALT_TEST_DOMAIN"
	altDomainEnv := os.Getenv(altDomainEnvVariable)
	if altDomainEnv == "" {
		t.Skipf("Skipping as %q is not specified", altDomainEnvVariable)
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_custom_hostname_binding", "test")
	r := ServiceCustomHostnameBindingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleConfig(data, appServiceEnv, domainEnv, altDomainEnv),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func testAccAppServiceCustomHostnameBinding_ssl(t *testing.T, appServiceEnv, domainEnv string) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_custom_hostname_binding", "test")
	r := ServiceCustomHostnameBindingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sslConfig(data, appServiceEnv, domainEnv),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ServiceCustomHostnameBindingResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AppServiceCustomHostnameBindingID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Web.AppServicesClient.GetHostNameBinding(ctx, id.ResourceGroup, id.AppServiceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving App Service Custom Hostname Binding %q (App Service %q / Resource Group %q): %+v", id.Name, id.AppServiceName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.HostNameBindingProperties != nil), nil
}

func (ServiceCustomHostnameBindingResource) basicConfig(data acceptance.TestData, appServiceName string, domain string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}

resource "azurerm_app_service_custom_hostname_binding" "test" {
  hostname            = "%s"
  app_service_name    = azurerm_app_service.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, appServiceName, domain)
}

func (r ServiceCustomHostnameBindingResource) requiresImport(data acceptance.TestData, appServiceName string, domain string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_custom_hostname_binding" "import" {
  hostname            = azurerm_app_service_custom_hostname_binding.test.name
  app_service_name    = azurerm_app_service_custom_hostname_binding.test.app_service_name
  resource_group_name = azurerm_app_service_custom_hostname_binding.test.resource_group_name
}
`, r.basicConfig(data, appServiceName, domain))
}

func (r ServiceCustomHostnameBindingResource) multipleConfig(data acceptance.TestData, appServiceName, domain, altDomain string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_custom_hostname_binding" "test2" {
  hostname            = "%s"
  app_service_name    = azurerm_app_service.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, r.basicConfig(data, appServiceName, domain), altDomain)
}

func (r ServiceCustomHostnameBindingResource) sslConfig(data acceptance.TestData, appServiceName, domain string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}

data "azurerm_client_config" "test" {
}

resource "azurerm_key_vault" "test" {
  name                = "acct-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.test.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id               = data.azurerm_client_config.test.tenant_id
    object_id               = data.azurerm_client_config.test.object_id
    secret_permissions      = ["delete", "get", "set"]
    certificate_permissions = ["create", "delete", "get", "import"]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acct-%d"
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

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      extended_key_usage = ["1.3.6.1.5.5.7.3.1"]

      key_usage = [
        "digitalSignature",
        "keyEncipherment",
      ]

      subject            = "CN=%s"
      validity_in_months = 12
    }
  }
}

data "azurerm_key_vault_secret" "test" {
  name         = azurerm_key_vault_certificate.test.name
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_app_service_certificate" "test" {
  name                = "acctestCert-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  pfx_blob            = data.azurerm_key_vault_secret.test.value
}

resource "azurerm_app_service_custom_hostname_binding" "test" {
  hostname            = "%s"
  app_service_name    = azurerm_app_service.test.name
  resource_group_name = azurerm_resource_group.test.name
  ssl_state           = "SniEnabled"
  thumbprint          = azurerm_app_service_certificate.test.thumbprint
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, appServiceName, data.RandomInteger, data.RandomInteger, domain, data.RandomInteger, domain)
}
