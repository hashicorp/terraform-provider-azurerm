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

type AppServiceManagedCertificateResource struct{}

func TestAccAppServiceManagedCertificate_basicLinux(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE and/or ARM_TEST_DATA_RESOURCE_GROUP are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_managed_certificate", "test")
	r := AppServiceManagedCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicLinux(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccAppServiceManagedCertificate_requiresImport(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE and/or ARM_TEST_DATA_RESOURCE_GROUP are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_managed_certificate", "test")
	r := AppServiceManagedCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicLinux(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAppServiceManagedCertificate_basicWindows(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE and/or ARM_TEST_DATA_RESOURCE_GROUP are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_managed_certificate", "test")
	r := AppServiceManagedCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWindows(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (t AppServiceManagedCertificateResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ManagedCertificateID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Web.CertificatesClient.Get(ctx, id.ResourceGroup, id.CertificateName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("App Service Managed Certificate %q (resource group %q) does not exist", id.CertificateName, id.ResourceGroup)
	}

	return utils.Bool(resp.CertificateProperties != nil), nil
}

func (t AppServiceManagedCertificateResource) basicLinux(data acceptance.TestData) string {
	template := t.linuxTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_managed_certificate" "test" {
  custom_hostname_binding_id = azurerm_app_service_custom_hostname_binding.test.id
}

`, template)
}

func (t AppServiceManagedCertificateResource) requiresImport(data acceptance.TestData) string {
	template := t.basicLinux(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_managed_certificate" "import" {
  custom_hostname_binding_id = azurerm_app_service_managed_certificate.test.custom_hostname_binding_id
}

`, template)
}

func (AppServiceManagedCertificateResource) linuxTemplate(data acceptance.TestData) string {
	dnsZone := os.Getenv("ARM_TEST_DNS_ZONE")
	dataResourceGroup := os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP")
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

data "azurerm_dns_zone" "test" {
  name                = "%s"
  resource_group_name = "%s"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "%s"
  zone_name           = data.azurerm_dns_zone.test.name
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  ttl                 = 300
  record              = azurerm_app_service.test.default_site_hostname
}

resource "azurerm_dns_txt_record" "test" {
  name                = join(".", ["asuid", "%s"])
  zone_name           = data.azurerm_dns_zone.test.name
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  ttl                 = 300

  record {
    value = azurerm_app_service.test.custom_domain_verification_id
  }
}

resource "azurerm_app_service_custom_hostname_binding" "test" {
  hostname            = join(".", [azurerm_dns_cname_record.test.name, azurerm_dns_cname_record.test.zone_name])
  app_service_name    = azurerm_app_service.test.name
  resource_group_name = azurerm_resource_group.test.name
}

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, dnsZone, dataResourceGroup, data.RandomString, data.RandomString)
}

func (t AppServiceManagedCertificateResource) basicWindows(data acceptance.TestData) string {
	template := t.windowsTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_managed_certificate" "test" {
  custom_hostname_binding_id = azurerm_app_service_custom_hostname_binding.test.id
}

`, template)
}

func (AppServiceManagedCertificateResource) windowsTemplate(data acceptance.TestData) string {
	dnsZone := os.Getenv("ARM_TEST_DNS_ZONE")
	dataResourceGroup := os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP")
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

  sku {
    tier = "Basic"
    size = "B1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctest%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}

data "azurerm_dns_zone" "test" {
  name                = "%s"
  resource_group_name = "%s"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "%s"
  zone_name           = data.azurerm_dns_zone.test.name
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  ttl                 = 300
  record              = azurerm_app_service.test.default_site_hostname
}

resource "azurerm_dns_txt_record" "test" {
  name                = join(".", ["asuid", "%s"])
  zone_name           = data.azurerm_dns_zone.test.name
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  ttl                 = 300

  record {
    value = azurerm_app_service.test.custom_domain_verification_id
  }
}

resource "azurerm_app_service_custom_hostname_binding" "test" {
  hostname            = join(".", [azurerm_dns_cname_record.test.name, azurerm_dns_cname_record.test.zone_name])
  app_service_name    = azurerm_app_service.test.name
  resource_group_name = azurerm_resource_group.test.name
}

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, dnsZone, dataResourceGroup, data.RandomString, data.RandomString)
}
