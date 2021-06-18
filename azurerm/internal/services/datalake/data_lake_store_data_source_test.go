package datalake_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type DataLakeStoreDataSource struct {
}

func TestAccDataLakeStoreDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_lake_store", "test")
	r := DataLakeStoreDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tier").HasValue("Consumption"),
			),
		},
	})
}

func TestAccDataLakeStoreDataSource_tier(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_lake_store", "test")
	r := DataLakeStoreDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.tier(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tier").HasValue("Commitment_1TB"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.hello").HasValue("world"),
			),
		},
	})
}

func (DataLakeStoreDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datalake-%d"
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

func (DataLakeStoreDataSource) tier(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datalake-%d"
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
