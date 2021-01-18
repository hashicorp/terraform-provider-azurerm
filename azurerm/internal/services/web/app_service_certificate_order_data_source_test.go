package web_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type AppServiceCertificateOrderDataSource struct{}

func TestAccDataSourceAppServiceCertificateOrder_basic(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_APP_SERVICE_CERTIFICATE") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_APP_SERVICE_CERTIFICATE is not specified")
		return
	}

	data := acceptance.BuildTestData(t, "data.azurerm_app_service_certificate_order", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceCertificateOrderDataSource{}.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("csr").Exists(),
				check.That(data.ResourceName).Key("domain_verification_token").Exists(),
				check.That(data.ResourceName).Key("distinguished_name").HasValue("CN=example.com"),
				check.That(data.ResourceName).Key("product_type").HasValue("Standard"),
			),
		},
	})
}

func TestAccDataSourceAppServiceCertificateOrder_wildcard(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_APP_SERVICE_CERTIFICATE") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_APP_SERVICE_CERTIFICATE is not specified")
		return
	}

	data := acceptance.BuildTestData(t, "data.azurerm_app_service_certificate_order", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceCertificateOrderDataSource{}.wildcard(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("csr").Exists(),
				check.That(data.ResourceName).Key("domain_verification_token").Exists(),
				check.That(data.ResourceName).Key("distinguished_name").HasValue("CN=*.example.com"),
				check.That(data.ResourceName).Key("product_type").HasValue("WildCard"),
			),
		},
	})
}

func TestAccDataSourceAppServiceCertificateOrder_complete(t *testing.T) {
	if os.Getenv("ARM_RUN_TEST_APP_SERVICE_CERTIFICATE") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_APP_SERVICE_CERTIFICATE is not specified")
		return
	}

	data := acceptance.BuildTestData(t, "data.azurerm_app_service_certificate_order", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceCertificateOrderDataSource{}.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("csr").Exists(),
				check.That(data.ResourceName).Key("domain_verification_token").Exists(),
				check.That(data.ResourceName).Key("distinguished_name").HasValue("CN=example.com"),
				check.That(data.ResourceName).Key("product_type").HasValue("Standard"),
				check.That(data.ResourceName).Key("validity_in_years").HasValue("1"),
				check.That(data.ResourceName).Key("auto_renew").HasValue("false"),
				check.That(data.ResourceName).Key("key_size").HasValue("4096"),
			),
		},
	})
}

func (d AppServiceCertificateOrderDataSource) basic(data acceptance.TestData) string {
	config := AppServiceCertificateOrderResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service_certificate_order" "test" {
  name                = azurerm_app_service_certificate_order.test.name
  resource_group_name = azurerm_app_service_certificate_order.test.resource_group_name
}
`, config)
}

func (d AppServiceCertificateOrderDataSource) wildcard(data acceptance.TestData) string {
	config := AppServiceCertificateOrderResource{}.wildcard(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service_certificate_order" "test" {
  name                = azurerm_app_service_certificate_order.test.name
  resource_group_name = azurerm_app_service_certificate_order.test.resource_group_name
}
`, config)
}

func (d AppServiceCertificateOrderDataSource) complete(data acceptance.TestData) string {
	config := AppServiceCertificateOrderResource{}.complete(data, 4096)
	return fmt.Sprintf(`
%s

data "azurerm_app_service_certificate_order" "test" {
  name                = azurerm_app_service_certificate_order.test.name
  resource_group_name = azurerm_app_service_certificate_order.test.resource_group_name
}
`, config)
}
