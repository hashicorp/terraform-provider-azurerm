// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StaticWebAppCustomDomainResource struct{}

func TestAccAzureStaticSiteCustomDomain_basic(t *testing.T) {
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

func TestAccAzureStaticSiteCustomDomain_requiresImport(t *testing.T) {
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
	id, err := parse.StaticSiteCustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Web.StaticSitesClient.GetStaticSiteCustomDomain(ctx, id.ResourceGroup, id.StaticSiteName, id.CustomDomainName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Static Site %q: %+v", id, err)
	}

	return utils.Bool(true), nil
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

func (r StaticWebAppCustomDomainResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_static_web_app_custom_domain" "import" {
  static_web_app_id = azurerm_static_web_app.test.id
  domain_name       = "acctestSS-%d.contoso.com"
  validation_type   = "dns-txt-token"
}
`, template, data.RandomInteger)
}
