package logic_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccLogicAppIntegrationAccountDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_logic_app_integration_account", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogicAppIntegrationAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLogicAppIntegrationAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogicAppIntegrationAccountExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "sku_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "tags.%"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "tags.ENV"),
				),
			},
		},
	})
}

func testAccDataSourceLogicAppIntegrationAccount_basic(data acceptance.TestData) string {
	config := testAccAzureRMLogicAppIntegrationAccount_complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_logic_app_integration_account" "test" {
  name                = azurerm_logic_app_integration_account.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, config)
}
