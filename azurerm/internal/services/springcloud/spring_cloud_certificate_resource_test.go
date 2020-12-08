package springcloud_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSpringCloudCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSpringCloudCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudCertificate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudCertificateExists(data.ResourceName),
				),
			},
			data.ImportStep("key_vault_certificate_id"),
		},
	})
}

func TestAccAzureRMSpringCloudCertificate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSpringCloudCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudCertificate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudCertificateExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSpringCloudCertificate_requiresImport),
		},
	})
}

func testCheckAzureRMSpringCloudCertificateExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Spring Cloud Certificate not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["service_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).AppPlatform.CertificatesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, resourceGroup, serviceName, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Spring Cloud Certificate %q (Spring Cloud Name %q / Resource Group %q) does not exist", name, serviceName, resourceGroup)
			}
			return fmt.Errorf("bad: Get on AppPlatform.CertificatesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSpringCloudCertificateDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).AppPlatform.CertificatesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_spring_cloud_certificate" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["service_name"]
		resp, err := client.Get(ctx, resGroup, serviceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on AppPlatform.CertificatesClient: %+v", err)
			}
			return nil
		}
		return fmt.Errorf("expected no spring cloud certificate but found %+v", resp)
	}

	return nil
}

func testAccAzureRMSpringCloudCertificate_basic(data acceptance.TestData) string {
	template := testAccAzureRMSpringCloudCertificate_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_certificate" "test" {
  name                     = "acctest-scc-%d"
  resource_group_name      = azurerm_spring_cloud_service.test.resource_group_name
  service_name             = azurerm_spring_cloud_service.test.name
  key_vault_certificate_id = azurerm_key_vault_certificate.test.id
}
`, template, data.RandomInteger)
}

func testAccAzureRMSpringCloudCertificate_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMSpringCloudCertificate_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_certificate" "import" {
  name                     = azurerm_spring_cloud_certificate.test.name
  resource_group_name      = azurerm_spring_cloud_certificate.test.resource_group_name
  service_name             = azurerm_spring_cloud_certificate.test.service_name
  key_vault_certificate_id = azurerm_spring_cloud_certificate.test.key_vault_certificate_id
}
`, template)
}

func testAccAzureRMSpringCloudCertificate_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%d"
  location = "%s"
}

data "azurerm_client_config" "current" {
}

data "azuread_service_principal" "test" {
  display_name = "Azure Spring Cloud Domain-Management"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id               = data.azurerm_client_config.current.tenant_id
    object_id               = data.azurerm_client_config.current.object_id
    secret_permissions      = ["set"]
    certificate_permissions = ["create", "delete", "get", "update"]
  }

  access_policy {
    tenant_id               = data.azurerm_client_config.current.tenant_id
    object_id               = data.azuread_service_principal.test.object_id
    secret_permissions      = ["get", "list"]
    certificate_permissions = ["get", "list"]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=contoso.com"
      validity_in_months = 12
    }
  }
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger)
}
