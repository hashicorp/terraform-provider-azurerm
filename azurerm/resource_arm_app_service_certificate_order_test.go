package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAppServiceCertificateOrder_basic(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_APP_SERVICE_CERTIFICATE") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_APP_SERVICE_CERTIFICATE is not specified")
		return
	}

	resourceName := "azurerm_app_service_certificate_order.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppServiceCertificateOrder_basic(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceCertificateOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceCertificateOrderExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "csr"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_verification_token"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name", "CN=example.com"),
					resource.TestCheckResourceAttr(resourceName, "product_type", "Standard"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppServiceCertificateOrder_wildcard(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_APP_SERVICE_CERTIFICATE") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_APP_SERVICE_CERTIFICATE is not specified")
		return
	}

	resourceName := "azurerm_app_service_certificate_order.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppServiceCertificateOrder_wildcard(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceCertificateOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceCertificateOrderExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "csr"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_verification_token"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name", "CN=*.example.com"),
					resource.TestCheckResourceAttr(resourceName, "product_type", "WildCard"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppServiceCertificateOrder_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	if os.Getenv("ARM_RUN_TEST_APP_SERVICE_CERTIFICATE") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_APP_SERVICE_CERTIFICATE is not specified")
		return
	}

	resourceName := "azurerm_app_service_certificate_order.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceCertificateOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceCertificateOrder_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceCertificateOrderExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMAppServiceCertificateOrder_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_app_service_certificate_order"),
			},
		},
	})
}

func TestAccAzureRMAppServiceCertificateOrder_complete(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_APP_SERVICE_CERTIFICATE") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_APP_SERVICE_CERTIFICATE is not specified")
		return
	}

	resourceName := "azurerm_app_service_certificate_order.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAppServiceCertificateOrder_complete(ri, acceptance.Location(), 4096)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceCertificateOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceCertificateOrderExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "csr"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_verification_token"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name", "CN=example.com"),
					resource.TestCheckResourceAttr(resourceName, "product_type", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "validity_in_years", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(resourceName, "key_size", "4096"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppServiceCertificateOrder_update(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_APP_SERVICE_CERTIFICATE") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_APP_SERVICE_CERTIFICATE is not specified")
		return
	}

	resourceName := "azurerm_app_service_certificate_order.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceCertificateOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceCertificateOrder_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceCertificateOrderExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "csr"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_verification_token"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name", "CN=example.com"),
					resource.TestCheckResourceAttr(resourceName, "product_type", "Standard"),
				),
			},
			{
				Config: testAccAzureRMAppServiceCertificateOrder_complete(ri, acceptance.Location(), 2048), // keySize cannot be updated
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceCertificateOrderExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "domain_verification_token"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name", "CN=example.com"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccAzureRMAppServiceCertificateOrder_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceCertificateOrderExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "domain_verification_token"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name", "CN=example.com"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
					resource.TestCheckResourceAttr(resourceName, "key_size", "2048"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMAppServiceCertificateOrderDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Web.CertificatesOrderClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
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

		client := acceptance.AzureProvider.Meta().(*clients.Client).Web.CertificatesOrderClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
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

func testAccAzureRMAppServiceCertificateOrder_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_certificate_order" "test" {
  name                = "acctestASCO-%d"
  location            = "global"
  resource_group_name = "${azurerm_resource_group.test.name}"
  distinguished_name  = "CN=example.com"
  product_type        = "Standard"
}
`, rInt, location, rInt)
}

func testAccAzureRMAppServiceCertificateOrder_wildcard(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_certificate_order" "test" {
  name                = "acctestASCO-%d"
  location            = "global"
  resource_group_name = "${azurerm_resource_group.test.name}"
  distinguished_name  = "CN=*.example.com"
  product_type        = "WildCard"
}
`, rInt, location, rInt)
}

func testAccAzureRMAppServiceCertificateOrder_requiresImport(rInt int, location string) string {
	template := testAccAzureRMAppServiceCertificateOrder_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_certificate_order" "import" {
  name                = "${azurerm_app_service_certificate_order.test.name}"
  location            = "${azurerm_app_service_certificate_order.test.location}"
  resource_group_name = "${azurerm_app_service_certificate_order.test.resource_group_name}"
  distinguished_name  = "${azurerm_app_service_certificate_order.test.distinguished_name}"
  product_type        = "${azurerm_app_service_certificate_order.test.product_type}"
}
`, template)
}

func testAccAzureRMAppServiceCertificateOrder_complete(rInt int, location string, keySize int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_certificate_order" "test" {
  name                = "acctestASCO-%d"
  location            = "global"
  resource_group_name = "${azurerm_resource_group.test.name}"
  distinguished_name  = "CN=example.com"
  product_type        = "Standard"
  auto_renew          = false
  validity_in_years   = 1
  key_size            = %d
}
`, rInt, location, rInt, keySize)
}
