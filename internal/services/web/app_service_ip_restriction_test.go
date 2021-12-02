package web_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AppServiceIpRestrictionResource struct {
}

func TestAccAppServiceIpRestriction_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_ip_restriction", "test")
	r := AppServiceIpRestrictionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_restriction.0.ip_address").HasValue("10.10.10.10/32"),
				check.That(data.ResourceName).Key("ip_restriction.0.action").HasValue("Allow"),
			),
		},
		data.ImportStep(),
	})
}

func (t AppServiceIpRestrictionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	siteName := id.Path["sites"]
	restrictionName := id.Path["ipRestriction"]

	resp, err := clients.Web.AppServicesClient.Get(ctx, id.ResourceGroup, siteName)
	if err != nil {
		return nil, fmt.Errorf("reading App Service %q (Resource Group %q): %s", siteName, id.ResourceGroup, err)
	}
	if resp.SiteConfig == nil || resp.SiteConfig.IPSecurityRestrictions == nil {
		return nil, fmt.Errorf("failed reading IP Restrictions for %q (resource group %q)", siteName, id.ResourceGroup)
	}

	idx, _ := web.FindIPRestriction(resp.SiteConfig.IPSecurityRestrictions, restrictionName)
	return utils.Bool(idx >= 0), nil
}

func (r AppServiceIpRestrictionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_ip_restriction" "test" {
  app_service_id = azurerm_app_service.test.id

  ip_restriction {
	name       = "basic"
	ip_address = "10.10.10.10/32"
	action     = "Allow"
  }
}
`, template)
}

func (r AppServiceIpRestrictionResource) template(data acceptance.TestData) string {
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
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
