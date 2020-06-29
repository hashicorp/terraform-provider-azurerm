package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMDataSourceTrafficManagerGeographicalLocation_europe(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_traffic_manager_geographical_location", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceTrafficManagerGeographicalLocation_template("Europe"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "id", "GEO-EU"),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "Europe"),
				),
			},
		},
	})
}

func TestAccAzureRMDataSourceTrafficManagerGeographicalLocation_germany(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_traffic_manager_geographical_location", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceTrafficManagerGeographicalLocation_template("Germany"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "id", "DE"),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "Germany"),
				),
			},
		},
	})
}

func TestAccAzureRMDataSourceTrafficManagerGeographicalLocation_unitedKingdom(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_traffic_manager_geographical_location", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceTrafficManagerGeographicalLocation_template("United Kingdom"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "id", "GB"),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "United Kingdom"),
				),
			},
		},
	})
}

func TestAccAzureRMDataSourceTrafficManagerGeographicalLocation_world(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_traffic_manager_geographical_location", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceTrafficManagerGeographicalLocation_template("World"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "id", "WORLD"),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "World"),
				),
			},
		},
	})
}

func testAccAzureRMDataSourceTrafficManagerGeographicalLocation_template(name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_traffic_manager_geographical_location" "test" {
  name = "%s"
}
`, name)
}
