package managedapplications_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMManagedApplicationDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_application_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedApplicationDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceManagedApplicationDefinition_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceManagedApplicationDefinition_basic(data acceptance.TestData) string {
	config := testAccAzureRMManagedApplicationDefinition_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_managed_application_definition" "test" {
  name                = azurerm_managed_application_definition.test.name
  resource_group_name = azurerm_managed_application_definition.test.resource_group_name
}
`, config)
}
