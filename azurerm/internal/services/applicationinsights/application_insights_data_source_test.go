package applicationinsights_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceApplicationInsights_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_application_insights", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApplicationInsights_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "instrumentation_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "app_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttr(data.ResourceName, "application_type", "other"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.foo", "bar"),
				),
			},
		},
	})
}

func testAccResourceApplicationInsights_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "other"

  tags = {
    "foo" = "bar"
  }
}

data "azurerm_application_insights" "test" {
  resource_group_name = azurerm_resource_group.test.name
  name                = azurerm_application_insights.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}
