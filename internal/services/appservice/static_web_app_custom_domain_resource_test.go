// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StaticWebAppCustomDomainResource struct{}

func TestAccAzureStaticWebAppCustomDomain_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app_custom_domain", "test")
	r := StaticWebAppCustomDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("validation_token").Exists(),
			),
		},
		data.ImportStep("validation_type", "validation_token"),
	})
}

func TestAccAzureStaticWebAppCustomDomain_cnameValidation(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" {
		t.Skipf("Skipping Static Web App Test CNAME test as ARM_TEST_DNS_ZONE is not set.")
	}
	data := acceptance.BuildTestData(t, "azurerm_static_web_app_custom_domain", "test")
	r := StaticWebAppCustomDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cnameValidation(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("validation_type", "validation_token"),
	})
}

func TestAccAzureStaticWebAppCustomDomain_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_web_app_custom_domain", "test")
	r := StaticWebAppCustomDomainResource{}

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

func (r StaticWebAppCustomDomainResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := staticsites.ParseCustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppService.StaticSitesClient.GetStaticSiteCustomDomain(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r StaticWebAppCustomDomainResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_static_web_app" "test" {
  name                = "acctestSS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_size            = "Standard"
  sku_tier            = "Standard"
}

resource "azurerm_static_web_app_custom_domain" "test" {
  static_web_app_id = azurerm_static_web_app.test.id
  domain_name       = "acctestSS-%d.contoso.com"
  validation_type   = "dns-txt-token"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r StaticWebAppCustomDomainResource) cnameValidation(data acceptance.TestData) string {
	dnsZone := os.Getenv("ARM_TEST_DNS_ZONE")
	dnsZoneRG := os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

data "azurerm_dns_zone" "test" {
  name                = "%[3]s"
  resource_group_name = "%[4]s"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "swa%[1]d"
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  zone_name           = data.azurerm_dns_zone.test.name
  ttl                 = 300
  record              = azurerm_static_web_app.test.default_host_name
}

resource "azurerm_static_web_app" "test" {
  name                = "acctestSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_size            = "Standard"
  sku_tier            = "Standard"
}

resource "azurerm_static_web_app_custom_domain" "test" {
  static_web_app_id = azurerm_static_web_app.test.id
  domain_name       = trimsuffix(azurerm_dns_cname_record.test.fqdn, ".")
  validation_type   = "cname-delegation"
}
`, data.RandomInteger, data.Locations.Primary, dnsZone, dnsZoneRG)
}

func (r StaticWebAppCustomDomainResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_static_web_app_custom_domain" "import" {
  static_web_app_id = azurerm_static_web_app_custom_domain.test.static_web_app_id
  domain_name       = azurerm_static_web_app_custom_domain.test.domain_name
  validation_type   = azurerm_static_web_app_custom_domain.test.validation_type
}
`, r.basic(data))
}
