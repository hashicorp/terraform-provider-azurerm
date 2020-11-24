package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMDigitalTwins_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_digital_twins_instance", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDigitalTwinsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDigitalTwinsInstance_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDigitalTwinsInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
		},
	})
}

func testAccDataSourceDigitalTwinsInstance_basic(data acceptance.TestData) string {
	config := testAccAzureRMDigitalTwinsInstance_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_digital_twins_instance" "test" {
  name                = azurerm_digital_twins_instance.test.name
  resource_group_name = azurerm_digital_twins_instance.test.resource_group_name
}
`, config)
}
