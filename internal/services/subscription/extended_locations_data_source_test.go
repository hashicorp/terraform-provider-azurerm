package subscription_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type ExtendedLocationsDataSource struct{}

func TestAccDataSourceExtendedLocations_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_extended_locations", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: ExtendedLocationsDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckOutput("extended_location", "microsoftlosangeles1"),
			),
		},
	})
}

func (d ExtendedLocationsDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_extended_locations" "test" {
  location = "West US"
}

output "extended_location" {
  value = data.azurerm_extended_locations.test.extended_locations[0]
}
`)
}
