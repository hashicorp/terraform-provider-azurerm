package datalake_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceDataLakeStore_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_lake_store", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataLakeStore_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataLakeStoreExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tier", "Consumption"),
				),
			},
		},
	})
}

func TestAccDataSourceDataLakeStore_tier(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_lake_store", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataLakeStore_tier(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "tier", "Commitment_1TB"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.hello", "world"),
				),
			},
		},
	})
}

func testAccDataSourceDataLakeStore_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_lake_store" "test" {
  name                = "unlikely23exst2acct%s"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_data_lake_store" "test" {
  name                = azurerm_data_lake_store.test.name
  resource_group_name = azurerm_data_lake_store.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.Locations.Primary)
}

func testAccDataSourceDataLakeStore_tier(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_lake_store" "test" {
  name                = "unlikely23exst2acct%s"
  location            = "%s"
  tier                = "Commitment_1TB"
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    hello = "world"
  }
}

data "azurerm_data_lake_store" "test" {
  name                = azurerm_data_lake_store.test.name
  resource_group_name = azurerm_data_lake_store.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.Locations.Primary)
}
