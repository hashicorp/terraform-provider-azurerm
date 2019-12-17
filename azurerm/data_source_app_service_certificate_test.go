package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAppServiceCertificate_basic(t *testing.T) {
	dataSourceName := "data.azurerm_app_service_certificate.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMAppServiceCertificate_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subject_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "issue_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "expiration_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "thumbprint"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMAppServiceCertificate_basic(rInt int, location string) string {
	template := testAccAzureRMAppServiceCertificatePfxNoPassword(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_app_service_certificate" "test" {
  name                = "${azurerm_app_service_certificate.test.name}"
  resource_group_name = "${azurerm_app_service_certificate.test.resource_group_name}"
}
`, template)
}
