package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAppServiceCertificateOrder_basic(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_APP_SERVICE_CERTIFICATE") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_APP_SERVICE_CERTIFICATE is not specified")
		return
	}

	data := acceptance.BuildTestData(t, "data.azurerm_app_service_certificate_order", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppServiceCertificateOrder_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "csr"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "domain_verification_token"),
					resource.TestCheckResourceAttr(data.ResourceName, "distinguished_name", "CN=example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "product_type", "Standard"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppServiceCertificateOrder_wildcard(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_APP_SERVICE_CERTIFICATE") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_APP_SERVICE_CERTIFICATE is not specified")
		return
	}

	data := acceptance.BuildTestData(t, "data.azurerm_app_service_certificate_order", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppServiceCertificateOrder_wildcard(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "csr"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "domain_verification_token"),
					resource.TestCheckResourceAttr(data.ResourceName, "distinguished_name", "CN=*.example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "product_type", "WildCard"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppServiceCertificateOrder_complete(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_APP_SERVICE_CERTIFICATE") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_APP_SERVICE_CERTIFICATE is not specified")
		return
	}

	data := acceptance.BuildTestData(t, "data.azurerm_app_service_certificate_order", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppServiceCertificateOrder_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "csr"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "domain_verification_token"),
					resource.TestCheckResourceAttr(data.ResourceName, "distinguished_name", "CN=example.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "product_type", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "validity_in_years", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "key_size", "4096"),
				),
			},
		},
	})
}

func testAccDataSourceAppServiceCertificateOrder_basic(data acceptance.TestData) string {
	config := testAccAzureRMAppServiceCertificateOrder_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service_certificate_order" "test" {
  name                = azurerm_app_service_certificate_order.test.name
  resource_group_name = azurerm_app_service_certificate_order.test.resource_group_name
}
`, config)
}

func testAccDataSourceAppServiceCertificateOrder_wildcard(data acceptance.TestData) string {
	config := testAccAzureRMAppServiceCertificateOrder_wildcard(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service_certificate_order" "test" {
  name                = azurerm_app_service_certificate_order.test.name
  resource_group_name = azurerm_app_service_certificate_order.test.resource_group_name
}
`, config)
}

func testAccDataSourceAppServiceCertificateOrder_complete(data acceptance.TestData) string {
	config := testAccAzureRMAppServiceCertificateOrder_complete(data, 4096)
	return fmt.Sprintf(`
%s

data "azurerm_app_service_certificate_order" "test" {
  name                = azurerm_app_service_certificate_order.test.name
  resource_group_name = azurerm_app_service_certificate_order.test.resource_group_name
}
`, config)
}
