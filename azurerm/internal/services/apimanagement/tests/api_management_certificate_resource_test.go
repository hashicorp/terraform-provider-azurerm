package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementCertificate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementCertificateExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "expiration"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subject"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "thumbprint"),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"data",
					"password",
				},
			},
		},
	})
}

func TestAccAzureRMApiManagementCertificate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementCertificate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementCertificateExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementCertificate_requiresImport),
		},
	})
}

func testCheckAzureRMAPIManagementCertificateDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.CertificatesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_certificate" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		resp, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}
	return nil
}

func testCheckAzureRMAPIManagementCertificateExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.CertificatesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		resp, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: API Management Certificate %q (Resource Group %q / API Management Service %q) does not exist", name, resourceGroup, serviceName)
			}
			return fmt.Errorf("Bad: Get on apiManagementCertificatesClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementCertificate_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_certificate" "test" {
  name                = "example-cert"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  data                = filebase64("testdata/keyvaultcert.pfx")
  password            = ""
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMApiManagementCertificate_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementCertificate_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_certificate" "import" {
  name                = azurerm_api_management_certificate.test.name
  api_management_name = azurerm_api_management_certificate.test.api_management_name
  resource_group_name = azurerm_api_management_certificate.test.resource_group_name
  data                = azurerm_api_management_certificate.test.data
  password            = azurerm_api_management_certificate.test.password
}
`, template)
}
