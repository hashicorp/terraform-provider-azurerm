package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMsubscriptionAlias_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subscription_alias", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMsubscriptionAliasDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcesubscriptionAlias_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMsubscriptionAliasExists(data.ResourceName),
				),
			},
		},
	})
}

func testAccDataSourcesubscriptionAlias_basic(data acceptance.TestData) string {
	config := testAccAzureRMsubscriptionAlias_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_subscription_alias" "test" {
  name = azurerm_subscription_alias.test.name
}
`, config)
}
