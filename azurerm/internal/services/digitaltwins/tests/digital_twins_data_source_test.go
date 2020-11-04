package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMdigitaltwinsDigitalTwin_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_digital_twins", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMdigitaltwinsDigitalTwinDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcedigitaltwinsDigitalTwin_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMdigitaltwinsDigitalTwinExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
		},
	})
}

func testAccDataSourcedigitaltwinsDigitalTwin_basic(data acceptance.TestData) string {
	config := testAccAzureRMdigitaltwinsDigitalTwin_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_digital_twins" "test" {
  name                = azurerm_digital_twins.test.name
  resource_group_name = azurerm_digital_twins.test.resource_group_name
}
`, config)
}
