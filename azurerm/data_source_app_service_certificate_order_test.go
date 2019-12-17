package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAppServiceCertificateOrder_basic(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_APP_SERVICE_CERTIFICATE") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_APP_SERVICE_CERTIFICATE is not specified")
		return
	}

	dataSourceName := "data.azurerm_app_service_certificate_order.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppServiceCertificateOrder_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "csr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "domain_verification_token"),
					resource.TestCheckResourceAttr(dataSourceName, "distinguished_name", "CN=example.com"),
					resource.TestCheckResourceAttr(dataSourceName, "product_type", "Standard"),
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

	dataSourceName := "data.azurerm_app_service_certificate_order.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppServiceCertificateOrder_wildcard(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "csr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "domain_verification_token"),
					resource.TestCheckResourceAttr(dataSourceName, "distinguished_name", "CN=*.example.com"),
					resource.TestCheckResourceAttr(dataSourceName, "product_type", "WildCard"),
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

	dataSourceName := "data.azurerm_app_service_certificate_order.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppServiceCertificateOrder_complete(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "csr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "domain_verification_token"),
					resource.TestCheckResourceAttr(dataSourceName, "distinguished_name", "CN=example.com"),
					resource.TestCheckResourceAttr(dataSourceName, "product_type", "Standard"),
					resource.TestCheckResourceAttr(dataSourceName, "validity_in_years", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "key_size", "4096"),
				),
			},
		},
	})
}

func testAccDataSourceAppServiceCertificateOrder_basic(rInt int, location string) string {
	config := testAccAzureRMAppServiceCertificateOrder_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_app_service_certificate_order" "test" {
  name                = "${azurerm_app_service_certificate_order.test.name}"
  resource_group_name = "${azurerm_app_service_certificate_order.test.resource_group_name}"
}
`, config)
}

func testAccDataSourceAppServiceCertificateOrder_wildcard(rInt int, location string) string {
	config := testAccAzureRMAppServiceCertificateOrder_wildcard(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_app_service_certificate_order" "test" {
  name                = "${azurerm_app_service_certificate_order.test.name}"
  resource_group_name = "${azurerm_app_service_certificate_order.test.resource_group_name}"
}
`, config)
}

func testAccDataSourceAppServiceCertificateOrder_complete(rInt int, location string) string {
	config := testAccAzureRMAppServiceCertificateOrder_complete(rInt, location, 4096)
	return fmt.Sprintf(`
%s

data "azurerm_app_service_certificate_order" "test" {
  name                = "${azurerm_app_service_certificate_order.test.name}"
  resource_group_name = "${azurerm_app_service_certificate_order.test.resource_group_name}"
}
`, config)
}
