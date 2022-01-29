package subscription_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type SubscriptionLocationsDataSource struct{}

func TestAccDataSourceSubscriptionLocations_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subscription_locations", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SubscriptionLocationsDataSource{}.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckOutput("extended_location", "microsoftlosangeles1"),
			),
		},
	})
}

func (d SubscriptionLocationsDataSource) basic() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_subscription_locations" "test" {
  subscription_id = data.azurerm_client_config.current.subscription_id
}

output "extended_location" {
  value = tolist([for location in data.azurerm_subscription_locations.test.locations : location if location.type == "EdgeZone"])[0].extended_location
}
`
}
