// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ServiceCustomHostnameBindingResource struct{}

func TestAccAppServiceCustomHostnameBinding_basic(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE and/or ARM_TEST_DATA_RESOURCE_GROUP are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_custom_hostname_binding", "test")
	r := ServiceCustomHostnameBindingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServiceCustomHostnameBinding_requiresImport(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE and/or ARM_TEST_DATA_RESOURCE_GROUP are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_custom_hostname_binding", "test")
	r := ServiceCustomHostnameBindingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImport(data)
		}),
	})
}

func TestAccAppServiceCustomHostnameBinding_multiple(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE and/or ARM_TEST_DATA_RESOURCE_GROUP are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_custom_hostname_binding", "test")
	r := ServiceCustomHostnameBindingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccAppServiceCustomHostnameBinding_ssl(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE and/or ARM_TEST_DATA_RESOURCE_GROUP are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_custom_hostname_binding", "test")
	r := ServiceCustomHostnameBindingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sslConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ServiceCustomHostnameBindingResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := webapps.ParseHostNameBindingID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Web.WebAppsClient.GetHostNameBinding(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ServiceCustomHostnameBindingResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_app_service_custom_hostname_binding" "test" {
  hostname            = trimsuffix(azurerm_dns_cname_record.test.fqdn, ".")
  app_service_name    = azurerm_app_service.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data))
}

func (r ServiceCustomHostnameBindingResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_custom_hostname_binding" "import" {
  hostname            = azurerm_app_service_custom_hostname_binding.test.hostname
  app_service_name    = azurerm_app_service_custom_hostname_binding.test.app_service_name
  resource_group_name = azurerm_app_service_custom_hostname_binding.test.resource_group_name
}
`, r.basicConfig(data))
}

func (r ServiceCustomHostnameBindingResource) multipleConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_dns_cname_record" "test2" {
  name                = "%[2]s"
  zone_name           = data.azurerm_dns_zone.test.name
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  ttl                 = 300
  record              = azurerm_app_service.test.default_site_hostname
}

resource "azurerm_dns_txt_record" "test2" {
  name                = join(".", ["asuid", "%[2]s"])
  zone_name           = data.azurerm_dns_zone.test.name
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  ttl                 = 300

  record {
    value = azurerm_app_service.test.custom_domain_verification_id
  }
}

resource "azurerm_app_service_custom_hostname_binding" "test2" {
  hostname            = trimsuffix(azurerm_dns_cname_record.test2.fqdn, ".")
  app_service_name    = azurerm_app_service.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, r.basicConfig(data), data.RandomStringOfLength(7))
}

func (r ServiceCustomHostnameBindingResource) sslConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

data "azurerm_client_config" "test" {}

resource "azurerm_key_vault" "test" {
  name                = "acct-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.test.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id               = data.azurerm_client_config.test.tenant_id
    object_id               = data.azurerm_client_config.test.object_id
    secret_permissions      = ["Delete", "Get", "Set"]
    certificate_permissions = ["Create", "Delete", "Get", "Import", "Purge"]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctest-%[2]d"
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

      subject            = "CN=${trimsuffix(azurerm_dns_cname_record.test.fqdn, ".")}"
      validity_in_months = 12
    }
  }
}

data "azurerm_key_vault_secret" "test" {
  name         = azurerm_key_vault_certificate.test.name
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_app_service_certificate" "test" {
  name                = "acctestCert-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  pfx_blob            = data.azurerm_key_vault_secret.test.value
}

resource "azurerm_app_service_custom_hostname_binding" "test" {
  hostname            = trimsuffix(azurerm_dns_cname_record.test.fqdn, ".")
  app_service_name    = azurerm_app_service.test.name
  resource_group_name = azurerm_resource_group.test.name
  ssl_state           = "SniEnabled"
  thumbprint          = azurerm_app_service_certificate.test.thumbprint
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r ServiceCustomHostnameBindingResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}

data "azurerm_dns_zone" "test" {
  name                = "%[3]s"
  resource_group_name = "%[4]s"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "%[5]s"
  zone_name           = data.azurerm_dns_zone.test.name
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  ttl                 = 300
  record              = azurerm_app_service.test.default_site_hostname
}

resource "azurerm_dns_txt_record" "test" {
  name                = join(".", ["asuid", "%[5]s"])
  zone_name           = data.azurerm_dns_zone.test.name
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  ttl                 = 300

  record {
    value = azurerm_app_service.test.custom_domain_verification_id
  }
}
`, data.RandomInteger, data.Locations.Primary, os.Getenv("ARM_TEST_DNS_ZONE"), os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP"), data.RandomString)
}
