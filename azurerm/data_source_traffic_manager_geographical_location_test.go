package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMDataSourceTrafficManagerGeographicalLocation_europe(t *testing.T) {
	dataSourceName := "data.azurerm_traffic_manager_geographical_location.test"
	config := testAccAzureRMDataSourceTrafficManagerGeographicalLocation_europe()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(*terraform.State) error {
			// nothing to do
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "GEO-EU"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "Europe"),
				),
			},
		},
	})
}

func TestAccAzureRMDataSourceTrafficManagerGeographicalLocation_germany(t *testing.T) {
	dataSourceName := "data.azurerm_traffic_manager_geographical_location.test"
	config := testAccAzureRMDataSourceTrafficManagerGeographicalLocation_germany()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(*terraform.State) error {
			// nothing to do
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "DE"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "Germany"),
				),
			},
		},
	})
}

func TestAccAzureRMDataSourceTrafficManagerGeographicalLocation_unitedKingdom(t *testing.T) {
	dataSourceName := "data.azurerm_traffic_manager_geographical_location.test"
	config := testAccAzureRMDataSourceTrafficManagerGeographicalLocation_unitedKingdom()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(*terraform.State) error {
			// nothing to do
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "GB"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "United Kingdom"),
				),
			},
		},
	})
}

func TestAccAzureRMDataSourceTrafficManagerGeographicalLocation_world(t *testing.T) {
	dataSourceName := "data.azurerm_traffic_manager_geographical_location.test"
	config := testAccAzureRMDataSourceTrafficManagerGeographicalLocation_world()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(*terraform.State) error {
			// nothing to do
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "WORLD"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "World"),
				),
			},
		},
	})
}

func testAccAzureRMDataSourceTrafficManagerGeographicalLocation_europe() string {
	return fmt.Sprintf(`
data "azurerm_traffic_manager_geographical_location" "test" {
  name = "Europe"
}
`)
}

func testAccAzureRMDataSourceTrafficManagerGeographicalLocation_germany() string {
	return fmt.Sprintf(`
data "azurerm_traffic_manager_geographical_location" "test" {
  name = "Germany"
}
`)
}

func testAccAzureRMDataSourceTrafficManagerGeographicalLocation_unitedKingdom() string {
	return fmt.Sprintf(`
data "azurerm_traffic_manager_geographical_location" "test" {
  name = "United Kingdom"
}
`)
}

func testAccAzureRMDataSourceTrafficManagerGeographicalLocation_world() string {
	return fmt.Sprintf(`
data "azurerm_traffic_manager_geographical_location" "test" {
  name = "World"
}
`)
}
