package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMDigitalTwins_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_iothub_digital_twins", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDigitalTwinsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDigitalTwins_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDigitalTwinsExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
		},
	})
}

func testAccDataSourceDigitalTwins_basic(data acceptance.TestData) string {
	config := testAccAzureRMDigitalTwins_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_iothub_digital_twins" "test" {
  name                = azurerm_iothub_digital_twins.test.name
  resource_group_name = azurerm_iothub_digital_twins.test.resource_group_name
}
`, config)
}
