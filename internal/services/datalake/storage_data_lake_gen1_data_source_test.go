package datalake_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type StorageDataLakeGen1FilesystemDataSource struct {
}

func TestAccStorageDataLakeGen1FilesystemDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_data_lake_gen1_filesystem", "test")
	r := StorageDataLakeGen1FilesystemDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tier").HasValue("Consumption"),
			),
		},
	})
}

func TestAccStorageDataLakeGen1FilesystemDataSource_tier(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_data_lake_gen1_filesystem", "test")
	r := StorageDataLakeGen1FilesystemDataSource{}

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

func (StorageDataLakeGen1FilesystemDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datalake-%d"
  location = "%s"
}

resource "azurerm_storage_data_lake_gen1_filesystem" "test" {
  name                = "unlikely23exst2acct%s"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_storage_data_lake_gen1_filesystem" "test" {
  name                = azurerm_storage_data_lake_gen1_filesystem.test.name
  resource_group_name = azurerm_storage_data_lake_gen1_filesystem.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.Locations.Primary)
}

func (StorageDataLakeGen1FilesystemDataSource) tier(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datalake-%d"
  location = "%s"
}

resource "azurerm_storage_data_lake_gen1_filesystem" "test" {
  name                = "unlikely23exst2acct%s"
  location            = "%s"
  tier                = "Commitment_1TB"
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    hello = "world"
  }
}

data "azurerm_storage_data_lake_gen1_filesystem" "test" {
  name                = azurerm_storage_data_lake_gen1_filesystem.test.name
  resource_group_name = azurerm_storage_data_lake_gen1_filesystem.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.Locations.Primary)
}
