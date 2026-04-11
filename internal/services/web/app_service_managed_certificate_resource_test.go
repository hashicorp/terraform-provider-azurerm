// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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
		data.ImportStep("custom_hostname_binding_id"),
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

func TestAccAppServiceManagedCertificate_completeLinux(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE and/or ARM_TEST_DATA_RESOURCE_GROUP are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_managed_certificate", "test")
	r := AppServiceManagedCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeLinux(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_hostname_binding_id"),
	})
}

func TestAccAppServiceManagedCertificate_updateLinux(t *testing.T) {
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
		data.ImportStep("custom_hostname_binding_id"),
		{
			Config: r.completeLinux(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_hostname_binding_id"),
		{
			Config: r.updateLinux(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_hostname_binding_id"),
		{
			Config: r.basicLinux(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_hostname_binding_id"),
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
		data.ImportStep("custom_hostname_binding_id"),
	})
}

func (t AppServiceManagedCertificateResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := certificates.ParseCertificateID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Web.CertificatesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %w", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (t AppServiceManagedCertificateResource) basicLinux(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_managed_certificate" "test" {
  custom_hostname_binding_id = azurerm_app_service_custom_hostname_binding.test.id
}
`, t.linuxTemplate(data))
}

func (t AppServiceManagedCertificateResource) completeLinux(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_managed_certificate" "test" {
  custom_hostname_binding_id = azurerm_app_service_custom_hostname_binding.test.id

  tags = {
    hello = "world"
  }
}
`, t.linuxTemplate(data))
}

func (t AppServiceManagedCertificateResource) updateLinux(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_managed_certificate" "test" {
  custom_hostname_binding_id = azurerm_app_service_custom_hostname_binding.test.id

  tags = {
    hello = "world"
    foo   = "bar"
  }
}
`, t.linuxTemplate(data))
}

func (t AppServiceManagedCertificateResource) basicWindows(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_managed_certificate" "test" {
  custom_hostname_binding_id = azurerm_app_service_custom_hostname_binding.test.id
}
`, t.windowsTemplate(data))
}

func (t AppServiceManagedCertificateResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_managed_certificate" "import" {
  custom_hostname_binding_id = azurerm_app_service_managed_certificate.test.custom_hostname_binding_id
}
`, t.basicLinux(data))
}

func (AppServiceManagedCertificateResource) linuxTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-asmc-%[1]d"
  location = "%[2]s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
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
  name                = "acctest%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}

data "azurerm_dns_zone" "test" {
  name                = "%s"
  resource_group_name = "%s"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "%[3]s"
  zone_name           = data.azurerm_dns_zone.test.name
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  ttl                 = 300
  record              = azurerm_app_service.test.default_site_hostname
}

resource "azurerm_dns_txt_record" "test" {
  name                = join(".", ["asuid", "%[3]s"])
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, os.Getenv("ARM_TEST_DNS_ZONE"), os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP"))
}

func (AppServiceManagedCertificateResource) windowsTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-asmc-%[1]d"
  location = "%[2]s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Basic"
    size = "B1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctest%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}

data "azurerm_dns_zone" "test" {
  name                = "%[4]s"
  resource_group_name = "%[5]s"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "%[3]s"
  zone_name           = data.azurerm_dns_zone.test.name
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  ttl                 = 300
  record              = azurerm_app_service.test.default_site_hostname
}

resource "azurerm_dns_txt_record" "test" {
  name                = join(".", ["asuid", "%[3]s"])
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, os.Getenv("ARM_TEST_DNS_ZONE"), os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP"))
}
