package web_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type StaticSiteResource struct{}

func TestAccAzureStaticSite_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_site", "test")
	r := StaticSiteResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_host_name").Exists(),
				check.That(data.ResourceName).Key("api_key").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureStaticSite_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_static_site", "test")
	r := StaticSiteResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r StaticSiteResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.StaticSiteID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Web.StaticSitesClient.GetStaticSite(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Static Site %q: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r StaticSiteResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_static_site" "test" {
  name                = "acctestSS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r StaticSiteResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_static_site" "import" {
  name                = azurerm_static_site.test.name
  location            = azurerm_static_site.test.location
  resource_group_name = azurerm_static_site.test.resource_group_name
}
`, template)
}
