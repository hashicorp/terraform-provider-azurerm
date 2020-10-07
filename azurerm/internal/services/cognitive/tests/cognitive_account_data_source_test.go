package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMCognitiveAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppCognitiveAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCognitiveAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCognitiveAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "Face"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
				),
			},
		},
	})
}

func testAccDataSourceCognitiveAccount_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
	name                = "acctestcogacc-%d"
	location            = azurerm_resource_group.test.location
	resource_group_name = azurerm_resource_group.test.name
	kind                = "Face"
	sku_name            = "S0"

  tags = {
    Acceptance = "Test"
  }
}
`, rInt, location, rString)
}
