package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAppServiceCertificateOrder_basic(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_APP_SERVICE_CERTIFICATE") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_APP_SERVICE_CERTIFICATE is not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_app_service_certificate_order", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceCertificateOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceCertificateOrder_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceCertificateOrderExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "csr"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "domain_verification_token"),
					resource.TestCheckResourceAttr(data.ResourceName, "distinguished_name", "CN=example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "product_type", "Standard"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceCertificateOrder_wildcard(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_APP_SERVICE_CERTIFICATE") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_APP_SERVICE_CERTIFICATE is not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_certificate_order", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceCertificateOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceCertificateOrder_wildcard(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceCertificateOrderExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "csr"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "domain_verification_token"),
					resource.TestCheckResourceAttr(data.ResourceName, "distinguished_name", "CN=*.example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "product_type", "WildCard"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceCertificateOrder_requiresImport(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_APP_SERVICE_CERTIFICATE") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_APP_SERVICE_CERTIFICATE is not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_certificate_order", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceCertificateOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceCertificateOrder_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceCertificateOrderExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAppServiceCertificateOrder_requiresImport),
		},
	})
}

func TestAccAzureRMAppServiceCertificateOrder_complete(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_APP_SERVICE_CERTIFICATE") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_APP_SERVICE_CERTIFICATE is not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_app_service_certificate_order", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceCertificateOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceCertificateOrder_complete(data, 4096),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceCertificateOrderExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "csr"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "domain_verification_token"),
					resource.TestCheckResourceAttr(data.ResourceName, "distinguished_name", "CN=example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "product_type", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "validity_in_years", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "key_size", "4096"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceCertificateOrder_update(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_APP_SERVICE_CERTIFICATE") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_APP_SERVICE_CERTIFICATE is not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_certificate_order", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceCertificateOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceCertificateOrder_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceCertificateOrderExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "csr"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "domain_verification_token"),
					resource.TestCheckResourceAttr(data.ResourceName, "distinguished_name", "CN=example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "product_type", "Standard"),
				),
			},
			{
				Config: testAccAzureRMAppServiceCertificateOrder_complete(data, 2048), // keySize cannot be updated
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceCertificateOrderExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "domain_verification_token"),
					resource.TestCheckResourceAttr(data.ResourceName, "distinguished_name", "CN=example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccAzureRMAppServiceCertificateOrder_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceCertificateOrderExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "domain_verification_token"),
					resource.TestCheckResourceAttr(data.ResourceName, "distinguished_name", "CN=example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_renew", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "key_size", "2048"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAppServiceCertificateOrderDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Web.CertificatesOrderClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_certificate_order" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMAppServiceCertificateOrderExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Web.CertificatesOrderClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		appServiceName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Service Certificate Order: %s", appServiceName)
		}

		resp, err := client.Get(ctx, resourceGroup, appServiceName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: App Service Certificate Order %q (resource group: %q) does not exist", appServiceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on appServicesCertificateOrderClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMAppServiceCertificateOrder_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_certificate_order" "test" {
  name                = "acctestASCO-%d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  distinguished_name  = "CN=example.com"
  product_type        = "Standard"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAppServiceCertificateOrder_wildcard(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_certificate_order" "test" {
  name                = "acctestASCO-%d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  distinguished_name  = "CN=*.example.com"
  product_type        = "WildCard"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAppServiceCertificateOrder_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceCertificateOrder_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_certificate_order" "import" {
  name                = azurerm_app_service_certificate_order.test.name
  location            = azurerm_app_service_certificate_order.test.location
  resource_group_name = azurerm_app_service_certificate_order.test.resource_group_name
  distinguished_name  = azurerm_app_service_certificate_order.test.distinguished_name
  product_type        = azurerm_app_service_certificate_order.test.product_type
}
`, template)
}

func testAccAzureRMAppServiceCertificateOrder_complete(data acceptance.TestData, keySize int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_certificate_order" "test" {
  name                = "acctestASCO-%d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  distinguished_name  = "CN=example.com"
  product_type        = "Standard"
  auto_renew          = false
  validity_in_years   = 1
  key_size            = %d
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, keySize)
}
