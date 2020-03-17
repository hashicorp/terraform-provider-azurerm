package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAppServiceCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMAppServiceCertificate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subject_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "issue_date"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "expiration_date"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "thumbprint"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMAppServiceCertificate_basic(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceCertificatePfxNoPassword(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service_certificate" "test" {
  name                = azurerm_app_service_certificate.test.name
  resource_group_name = azurerm_app_service_certificate.test.resource_group_name
}
`, template)
}
