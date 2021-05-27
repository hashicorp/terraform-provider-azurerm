package trafficmanager_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type TrafficManagerGeographicalLocationDataSource struct{}

func TestAccAzureRMDataSourceTrafficManagerGeographicalLocation_europe(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_traffic_manager_geographical_location", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: TrafficManagerGeographicalLocationDataSource{}.template("Europe"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue("GEO-EU"),
				check.That(data.ResourceName).Key("name").HasValue("Europe"),
			),
		},
	})
}

func TestAccAzureRMDataSourceTrafficManagerGeographicalLocation_germany(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_traffic_manager_geographical_location", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: TrafficManagerGeographicalLocationDataSource{}.template("Germany"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue("DE"),
				check.That(data.ResourceName).Key("name").HasValue("Germany"),
			),
		},
	})
}

func TestAccAzureRMDataSourceTrafficManagerGeographicalLocation_unitedKingdom(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_traffic_manager_geographical_location", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: TrafficManagerGeographicalLocationDataSource{}.template("United Kingdom"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue("GB"),
				check.That(data.ResourceName).Key("name").HasValue("United Kingdom"),
			),
		},
	})
}

func TestAccAzureRMDataSourceTrafficManagerGeographicalLocation_world(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_traffic_manager_geographical_location", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: TrafficManagerGeographicalLocationDataSource{}.template("World"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue("WORLD"),
				check.That(data.ResourceName).Key("name").HasValue("World"),
			),
		},
	})
}

func (d TrafficManagerGeographicalLocationDataSource) template(name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_traffic_manager_geographical_location" "test" {
  name = "%s"
}
`, name)
}
