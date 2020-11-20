package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AppServiceManagedCertificate struct{}

func TestAccAzureRMAppServiceManagedCertificate_basicLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_managed_certificate", "test")

	r := AppServiceManagedCertificate{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicLinux(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (t AppServiceManagedCertificate) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.AppServiceManagedCertificateID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Web.CertificatesClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("App Service Managed Certificate %q (resource group %q) does not exist", id.Name, id.ResourceGroup)
	}

	return utils.Bool(resp.CertificateProperties != nil), nil
}

func (t AppServiceManagedCertificate) basicLinux(data acceptance.TestData) string {
	template := t.linuxTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_managed_certificate" "test" {
  name                = "acctest-appMS-%d"
  canonical_name      = azurerm_app_service_custom_hostname_binding.test.hostname
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id


}

`, template, data.RandomInteger, data.Locations.Primary)
}

func (AppServiceManagedCertificate) linuxTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-asmc-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Linux"

  sku {
    tier = "Basic"
    size = "B1"
  }

  reserved = true
}

resource "azurerm_app_service" "test" {
  name                = "acctest%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dns_cname_record" "test" {
  name                = "%s"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  record              = azurerm_app_service.test.default_site_hostname
}

resource "azurerm_app_service_custom_hostname_binding" "test" {
  hostname            = join(".", [azurerm_dns_cname_record.test.name, azurerm_dns_cname_record.test.zone_name])
  app_service_name    = azurerm_app_service.test.name
  resource_group_name = azurerm_resource_group.test.name
}

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomString)
}
