package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMManagementPartner_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_management_partner", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceManagementPartner_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "partner_id"),
				),
			},
		},
	})
}

func testAccDataSourceManagementPartner_basic() string {
	config := testAccAzureRMManagementPartner_basic()
	return fmt.Sprintf(`
%s

data "azurerm_management_partner" "test" {
  partner_id = "${azurerm_management_partner.test.partner_id}"
}
`, config)
}
