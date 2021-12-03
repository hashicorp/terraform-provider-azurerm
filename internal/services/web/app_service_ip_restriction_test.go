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
				check.That(data.ResourceName).Key("ip_restriction.0.name").HasValue("basic-restriction"),
				check.That(data.ResourceName).Key("ip_restriction.0.ip_address").HasValue("10.10.10.10/32"),
				check.That(data.ResourceName).Key("ip_restriction.0.action").HasValue("Allow"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServiceIpRestriction_basicUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_ip_restriction", "test")
	r := AppServiceIpRestrictionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_restriction.0.name").HasValue("basic-restriction"),
				check.That(data.ResourceName).Key("ip_restriction.0.ip_address").HasValue("10.10.10.10/32"),
				check.That(data.ResourceName).Key("ip_restriction.0.action").HasValue("Allow"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_restriction.0.name").HasValue("basic-restriction"),
				check.That(data.ResourceName).Key("ip_restriction.0.ip_address").HasValue("20.20.20.20/32"),
				check.That(data.ResourceName).Key("ip_restriction.0.action").HasValue("Allow"),
				check.That(data.ResourceName).Key("ip_restriction.0.headers.0.x_azure_fdid.#").HasValue("1"),
				check.That(data.ResourceName).Key("ip_restriction.0.headers.0.x_azure_fdid.0").HasValue("55ce4ed1-4b06-4bf1-b40e-4638452104da"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServiceIpRestriction_headers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_ip_restriction", "test")
	r := AppServiceIpRestrictionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.headers(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_restriction.0.ip_address").HasValue("10.10.10.10/32"),
				check.That(data.ResourceName).Key("ip_restriction.0.name").HasValue("headers-restriction"),
				check.That(data.ResourceName).Key("ip_restriction.0.priority").HasValue("123"),
				check.That(data.ResourceName).Key("ip_restriction.0.action").HasValue("Allow"),
				check.That(data.ResourceName).Key("ip_restriction.0.headers.0.x_forwarded_for.#").HasValue("2"),
				check.That(data.ResourceName).Key("ip_restriction.0.headers.0.x_forwarded_for.0").HasValue("2002::1234:abcd:ffff:c0a8:101/64"),
				check.That(data.ResourceName).Key("ip_restriction.0.headers.0.x_forwarded_for.1").HasValue("9.9.9.9/32"),
				check.That(data.ResourceName).Key("ip_restriction.0.headers.0.x_azure_fdid.#").HasValue("1"),
				check.That(data.ResourceName).Key("ip_restriction.0.headers.0.x_azure_fdid.0").HasValue("55ce4ed1-4b06-4bf1-b40e-4638452104da"),
				check.That(data.ResourceName).Key("ip_restriction.0.headers.0.x_fd_health_probe.#").HasValue("1"),
				check.That(data.ResourceName).Key("ip_restriction.0.headers.0.x_fd_health_probe.0").HasValue("1"),
				check.That(data.ResourceName).Key("ip_restriction.0.headers.0.x_forwarded_host.#").HasValue("1"),
				check.That(data.ResourceName).Key("ip_restriction.0.headers.0.x_forwarded_host.0").HasValue("example.com"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServiceIpRestriction_multipleResources(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_ip_restriction", "test")
	r := AppServiceIpRestrictionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleResources(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_app_service_ip_restriction.test").Key("ip_restriction.0.name").HasValue("basic-restriction"),
				check.That("azurerm_app_service_ip_restriction.test").Key("ip_restriction.0.ip_address").HasValue("10.10.10.10/32"),
				check.That("azurerm_app_service_ip_restriction.test").Key("ip_restriction.0.action").HasValue("Allow"),
				check.That("azurerm_app_service_ip_restriction.test-1").Key("ip_restriction.0.ip_address").HasValue("20.20.20.20/32"),
				check.That("azurerm_app_service_ip_restriction.test-1").Key("ip_restriction.0.name").HasValue("headers-restriction"),
				check.That("azurerm_app_service_ip_restriction.test-1").Key("ip_restriction.0.priority").HasValue("123"),
				check.That("azurerm_app_service_ip_restriction.test-1").Key("ip_restriction.0.action").HasValue("Allow"),
				check.That("azurerm_app_service_ip_restriction.test-1").Key("ip_restriction.0.headers.0.x_forwarded_for.#").HasValue("2"),
				check.That("azurerm_app_service_ip_restriction.test-1").Key("ip_restriction.0.headers.0.x_forwarded_for.0").HasValue("2002::1234:abcd:ffff:c0a8:101/64"),
				check.That("azurerm_app_service_ip_restriction.test-1").Key("ip_restriction.0.headers.0.x_forwarded_for.1").HasValue("9.9.9.9/32"),
				check.That("azurerm_app_service_ip_restriction.test-1").Key("ip_restriction.0.headers.0.x_azure_fdid.#").HasValue("1"),
				check.That("azurerm_app_service_ip_restriction.test-1").Key("ip_restriction.0.headers.0.x_azure_fdid.0").HasValue("55ce4ed1-4b06-4bf1-b40e-4638452104da"),
				check.That("azurerm_app_service_ip_restriction.test-1").Key("ip_restriction.0.headers.0.x_fd_health_probe.#").HasValue("1"),
				check.That("azurerm_app_service_ip_restriction.test-1").Key("ip_restriction.0.headers.0.x_fd_health_probe.0").HasValue("1"),
				check.That("azurerm_app_service_ip_restriction.test-1").Key("ip_restriction.0.headers.0.x_forwarded_host.#").HasValue("1"),
				check.That("azurerm_app_service_ip_restriction.test-1").Key("ip_restriction.0.headers.0.x_forwarded_host.0").HasValue("example.com"),
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

	if restrictionName == "" {
		return nil, fmt.Errorf("ID was missing the 'ipRestriction' element")
	}

	resp, err := clients.Web.AppServicesClient.GetConfiguration(ctx, id.ResourceGroup, siteName)
	if err != nil {
		return nil, fmt.Errorf("reading App Service %q (Resource Group %q): %s", siteName, id.ResourceGroup, err)
	}
	if resp.SiteConfig == nil || resp.SiteConfig.IPSecurityRestrictions == nil {
		return utils.Bool(false), fmt.Errorf("failed reading IP Restrictions for %q (resource group %q)", siteName, id.ResourceGroup)
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
		name       = "basic-restriction"
		ip_address = "10.10.10.10/32"
		action     = "Allow"
  }
}
`, template)
}

func (r AppServiceIpRestrictionResource) basicUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_ip_restriction" "test" {
  app_service_id = azurerm_app_service.test.id

  ip_restriction {
		name       = "basic-restriction"
		ip_address = "20.20.20.20/32"
		action     = "Allow"
		headers {
			x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
		}
  }
}
`, template)
}

func (r AppServiceIpRestrictionResource) headers(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_ip_restriction" "test" {
	app_service_id = azurerm_app_service.test.id

	ip_restriction {
		ip_address = "10.10.10.10/32"
		name       = "headers-restriction"
		priority   = 123
		action     = "Allow"
		headers {
			x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
			x_fd_health_probe = ["1"]
			x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
			x_forwarded_host  = ["example.com"]
		}
	}
}
`, template)
}

func (r AppServiceIpRestrictionResource) multipleResources(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_ip_restriction" "test" {
  app_service_id = azurerm_app_service.test.id

  ip_restriction {
		name       = "basic-restriction"
		ip_address = "10.10.10.10/32"
		action     = "Allow"
  }
}

resource "azurerm_app_service_ip_restriction" "test-1" {
	app_service_id = azurerm_app_service.test.id

	ip_restriction {
		ip_address = "20.20.20.20/32"
		name       = "headers-restriction"
		priority   = 123
		action     = "Allow"
		headers {
			x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
			x_fd_health_probe = ["1"]
			x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
			x_forwarded_host  = ["example.com"]
		}
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
