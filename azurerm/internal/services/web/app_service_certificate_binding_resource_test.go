package web_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AppServiceCertificateBindingResource struct{}

func TestAccAzureRMAppServiceCertificateBinding_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_certificate_binding", "test")
	r := AppServiceCertificateBindingResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("thumbprint").Exists(),
				check.That(data.ResourceName).Key("ssl_state").HasValue("IpBasedEnabled"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMAppServiceCertificateBinding_basicSniEnabled(t *testing.T) {
       data := acceptance.BuildTestData(t, "azurerm_app_service_certificate_binding", "test")
       r := AppServiceCertificateBindingResource{}

	       data.ResourceTest(t, r, []resource.TestStep{
		               {
		                       Config: r.basicSniEnabled(data),
		                       Check: resource.ComposeTestCheckFunc(
		                               check.That(data.ResourceName).Key("thumbprint").Exists(),
		                               check.That(data.ResourceName).Key("ssl_state").HasValue("SniEnabled"),
		                       ),
		               },
		               data.ImportStep(),
		       })
}

func TestAccAzureRMAppServiceCertificateBinding_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_certificate_binding", "test")
	r := AppServiceCertificateBindingResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("thumbprint").Exists(),
				check.That(data.ResourceName).Key("ssl_state").HasValue("IpBasedEnabled"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t AppServiceCertificateBindingResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.CertificateBindingID(state.ID)
	if err != nil {
		return nil, err
	}

	binding, err := client.Web.AppServicesClient.GetHostNameBinding(ctx, id.HostnameBindingId.ResourceGroup, id.HostnameBindingId.SiteName, id.HostnameBindingId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(binding.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving App Service Hostname Binding %q (resource group %q) to check for Certificate Binding %q: %+v", id.HostnameBindingId.Name, id.HostnameBindingId.ResourceGroup, id.HostnameBindingId.Name, err)
	}
	certificate, err := client.Web.CertificatesClient.Get(ctx, id.CertificateId.ResourceGroup, id.CertificateId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(certificate.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving App Service Certificate %q (resource group %q) to check for Certificate Binding: %+v", id.CertificateId.Name, id.CertificateId.ResourceGroup, err)
	}
	bindingProps := binding.HostNameBindingProperties
	if bindingProps == nil || bindingProps.Thumbprint == nil {
		return utils.Bool(false), nil
	}
	certProps := certificate.CertificateProperties
	if certProps == nil || certProps.Thumbprint == nil {
		return nil, fmt.Errorf("reading Certificate thumbprint for verification on binding")
	}
	if *certProps.Thumbprint != *bindingProps.Thumbprint {
		return utils.Bool(false), nil
	}
	return utils.Bool(true), nil
}

func (t AppServiceCertificateBindingResource) basic(data acceptance.TestData) string {
	template := t.testAccCertificateBinding_template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_app_service_certificate_binding" "test" {
  hostname_binding_id = azurerm_app_service_custom_hostname_binding.test.id
  certificate_id      = azurerm_app_service_managed_certificate.test.id
  ssl_state           = "IpBasedEnabled"
}

%s
`, template)
}

func (t AppServiceCertificateBindingResource) basicSniEnabled(data acceptance.TestData) string {
	template := t.testAccCertificateBinding_template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_app_service_certificate_binding" "test" {
  hostname_binding_id = azurerm_app_service_custom_hostname_binding.test.id
  certificate_id      = azurerm_app_service_managed_certificate.test.id
  ssl_state           = "SniEnabled"
}

%s
`, template)
}

func (t AppServiceCertificateBindingResource) requiresImport(data acceptance.TestData) string {
	template := t.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_certificate_binding" "import" {
  hostname_binding_id = azurerm_app_service_certificate_binding.test.hostname_binding_id
  certificate_id      = azurerm_app_service_certificate_binding.test.certificate_id
  ssl_state           = azurerm_app_service_certificate_binding.test.ssl_state
}
`, template)
}

func (AppServiceCertificateBindingResource) testAccCertificateBinding_template(data acceptance.TestData) string {
	dnsZone := os.Getenv("ARM_TEST_DNS_ZONE")
	dataResourceGroup := os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP")
	return fmt.Sprintf(`

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

resource "azurerm_app_service_managed_certificate" "test" {
  custom_hostname_binding_id = azurerm_app_service_custom_hostname_binding.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, dnsZone, dataResourceGroup, data.RandomString, data.RandomString)
}
