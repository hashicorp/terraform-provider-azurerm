package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func preCheck(t *testing.T) {
	variables := []string{
		"ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME",
		"ARM_TEST_DNS_ZONE_NAME",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Skipf("`%s` must be set for acceptance tests!", variable)
		}
	}
}

func preCheckUserManagedCertificate(t *testing.T) {
	preCheck(t)
	variables := []string{
		"ARM_TEST_DNS_CERTIFICATE",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Skipf("`%s` must be set for acceptance tests!", variable)
		}
	}
}

func TestAccAzureRMCdnEndpointCustomDomain_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint_custom_domain", "test")
	dnsZoneRg, dnsZoneName, subdomain := os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME"), os.Getenv("ARM_TEST_DNS_ZONE_NAME"), acctest.RandStringFromCharSet(3, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t); preCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointCustomDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpointCustomDomain_basic(
					data, dnsZoneRg, dnsZoneName, subdomain,
				),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointCustomDomainExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCdnEndpointCustomDomain_httpsCdn(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint_custom_domain", "test")
	dnsZoneRg, dnsZoneName, subdomain := os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME"), os.Getenv("ARM_TEST_DNS_ZONE_NAME"), acctest.RandStringFromCharSet(3, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t); preCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointCustomDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpointCustomDomain_httpsCdn(
					data, dnsZoneRg, dnsZoneName, subdomain,
				),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointCustomDomainExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCdnEndpointCustomDomain_httpsCdnUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint_custom_domain", "test")
	dnsZoneRg, dnsZoneName, subdomain := os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME"), os.Getenv("ARM_TEST_DNS_ZONE_NAME"), acctest.RandStringFromCharSet(3, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t); preCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointCustomDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpointCustomDomain_basic(
					data, dnsZoneRg, dnsZoneName, subdomain,
				),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointCustomDomainExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCdnEndpointCustomDomain_httpsCdn(
					data, dnsZoneRg, dnsZoneName, subdomain,
				),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointCustomDomainExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCdnEndpointCustomDomain_basic(
					data, dnsZoneRg, dnsZoneName, subdomain,
				),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointCustomDomainExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCdnEndpointCustomDomain_httpsUserManaged(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint_custom_domain", "test")
	dnsZoneRg, dnsZoneName, subdomain, certificate := os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME"), os.Getenv("ARM_TEST_DNS_ZONE_NAME"), acctest.RandStringFromCharSet(3, acctest.CharSetAlpha), os.Getenv("ARM_TEST_DNS_CERTIFICATE")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t); preCheckUserManagedCertificate(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointCustomDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpointCustomDomain_httpsUserManaged(
					data, dnsZoneRg, dnsZoneName, subdomain, certificate,
				),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointCustomDomainExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCdnEndpointCustomDomain_httpsUserManagedUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint_custom_domain", "test")
	dnsZoneRg, dnsZoneName, subdomain, certificate := os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME"), os.Getenv("ARM_TEST_DNS_ZONE_NAME"), acctest.RandStringFromCharSet(3, acctest.CharSetAlpha), os.Getenv("ARM_TEST_DNS_CERTIFICATE")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t); preCheckUserManagedCertificate(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointCustomDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpointCustomDomain_basic(
					data, dnsZoneRg, dnsZoneName, subdomain,
				),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointCustomDomainExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCdnEndpointCustomDomain_httpsUserManaged(
					data, dnsZoneRg, dnsZoneName, subdomain, certificate,
				),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointCustomDomainExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCdnEndpointCustomDomain_basic(
					data, dnsZoneRg, dnsZoneName, subdomain,
				),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointCustomDomainExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCdnEndpointCustomDomain_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint_custom_domain", "test")
	dnsZoneRg, dnsZoneName, subdomain := os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME"), os.Getenv("ARM_TEST_DNS_ZONE_NAME"), acctest.RandStringFromCharSet(3, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointCustomDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpointCustomDomain_basic(
					data, dnsZoneRg, dnsZoneName, subdomain,
				),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointCustomDomainExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(
				func(data acceptance.TestData) string {
					template := testAccAzureRMCdnEndpointCustomDomain_basic(data, dnsZoneRg, dnsZoneName, subdomain)
					return fmt.Sprintf(`
%s

resource "azurerm_cdn_endpoint_custom_domain" "import" {
  name            = azurerm_cdn_endpoint_custom_domain.test.id
  cdn_endpoint_id = azurerm_cdn_endpoint_custom_domain.test.cdn_endpoint_id
  host_name       = azurerm_cdn_endpoint_custom_domain.test.host_name
}
`, template)
				}),
		},
	})
}

func testCheckAzureRMCdnEndpointCustomDomainExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Cdn.CustomDomainsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Cdn Endpoint Custom Domain not found: %s", resourceName)
		}

		id, err := parse.CdnEndpointCustomDomainID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Cdn Endpoint Custom Domain %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Getting on CDN.CustomDomains: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMCdnEndpointCustomDomainDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Cdn.CustomDomainsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cdn_endpoint_custom_domain" {
			continue
		}

		id, err := parse.CdnEndpointCustomDomainID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name)
		if err == nil {
			return fmt.Errorf("CDN.CustomDomains still exists")
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Getting on CDN.CustomDomains: %+v", err)
		}
		return nil
	}

	return nil
}

func testAccAzureRMCdnEndpointCustomDomain_basic(data acceptance.TestData, dnsZoneRg, dnsZoneName, subdomain string) string {
	template := testAccAzureRMCdnEndpointCustomDomain_template(data, dnsZoneRg, dnsZoneName, subdomain)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_endpoint_custom_domain" "test" {
  name            = "acctest-customdomain"
  cdn_endpoint_id = azurerm_cdn_endpoint.test.id
  host_name       = "${azurerm_dns_cname_record.test.name}.${data.azurerm_dns_zone.test.name}"
}
`, template)
}

func testAccAzureRMCdnEndpointCustomDomain_httpsCdn(data acceptance.TestData, dnsZoneRg, dnsZoneName, subdomain string) string {
	template := testAccAzureRMCdnEndpointCustomDomain_template(data, dnsZoneRg, dnsZoneName, subdomain)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_endpoint_custom_domain" "test" {
  name            = "acctest-customdomain"
  cdn_endpoint_id = azurerm_cdn_endpoint.test.id
  host_name       = "${azurerm_dns_cname_record.test.name}.${data.azurerm_dns_zone.test.name}"
  cdn_managed_https_settings {
    certificate_type = "Shared"
  }
}
`, template)
}

func testAccAzureRMCdnEndpointCustomDomain_httpsUserManaged(data acceptance.TestData, dnsZoneRg, dnsZoneName, subdomain, certificate string) string {
	template := testAccAzureRMCdnEndpointCustomDomain_template(data, dnsZoneRg, dnsZoneName, subdomain)
	return fmt.Sprintf(`
%[1]s

data "azurerm_client_config" "test" {
}

resource "azurerm_key_vault" "test" {
  name                = "testkeyvault-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.test.tenant_id

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    object_id = data.azurerm_client_config.test.object_id

    certificate_permissions = [
      "get",
      "delete",
      "import",
    ]

    key_permissions = [
      "get",
      "create",
    ]

    secret_permissions = [
      "get",
      "set",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    object_id = "4dbab725-22a4-44d5-ad44-c267ca38a954" # The Microsoft Azure.Cdn application object ID

    certificate_permissions = [
      "list",
      "get",
    ]

    secret_permissions = [
      "list",
      "get",
    ]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "testkeyvaultcert-%[2]d"
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("%[3]s")
    password = ""
  }

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }
  }
}

resource "azurerm_cdn_endpoint_custom_domain" "test" {
  name            = "testcustomdomain-%[2]d"
  cdn_endpoint_id = azurerm_cdn_endpoint.test.id
  host_name       = "${azurerm_dns_cname_record.test.name}.${data.azurerm_dns_zone.test.name}"
  user_managed_https_settings {
    subscription_id     = data.azurerm_client_config.test.subscription_id
    resource_group_name = azurerm_key_vault.test.resource_group_name
    vault_name          = azurerm_key_vault.test.name
    secret_name         = azurerm_key_vault_certificate.test.name
    secret_version      = azurerm_key_vault_certificate.test.version
  }
}
`, template, data.RandomInteger, certificate)
}

func testAccAzureRMCdnEndpointCustomDomain_template(data acceptance.TestData, dnsZoneRg, dnsZoneName, subdomain string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%[1]d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin {
    name      = "test"
    host_name = azurerm_storage_account.test.primary_blob_host
  }
}

data "azurerm_dns_zone" "test" {
  name                = "%[3]s"
  resource_group_name = "%[4]s"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "%[5]s"
  zone_name           = data.azurerm_dns_zone.test.name
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  ttl                 = 3600
  target_resource_id  = azurerm_cdn_endpoint.test.id
}
`, data.RandomIntOfLength(8), data.Locations.Primary, dnsZoneName, dnsZoneRg, subdomain)
}
