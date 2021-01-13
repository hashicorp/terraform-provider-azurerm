package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ManagedDiskDataSource struct {
}

func TestAccDataSourceManagedDisk_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_disk", "test")
	r := ManagedDiskDataSource{}

	name := fmt.Sprintf("acctestmanageddisk-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data, name, resourceGroupName),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(resourceGroupName),
				check.That(data.ResourceName).Key("storage_account_type").HasValue("Premium_LRS"),
				check.That(data.ResourceName).Key("disk_size_gb").HasValue("10"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("acctest"),
				check.That(data.ResourceName).Key("zones.#").HasValue("1"),
				check.That(data.ResourceName).Key("zones.0").HasValue("2"),
			),
		},
	})
}

func TestAccDataSourceManagedDisk_basic_withUltraSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_disk", "test")
	r := ManagedDiskDataSource{}

	name := fmt.Sprintf("acctestmanageddisk-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic_withUltraSSD(data, name, resourceGroupName),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("disk_iops_read_write").HasValue("101"),
				check.That(data.ResourceName).Key("disk_mbps_read_write").HasValue("10"),
			),
		},
	})
}

func (ManagedDiskDataSource) basic(data acceptance.TestData, name string, resourceGroupName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "%s"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = "10"
  zones                = ["2"]

  tags = {
    environment = "acctest"
  }
}

data "azurerm_managed_disk" "test" {
  name                = azurerm_managed_disk.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, resourceGroupName, data.Locations.Primary, name)
}

func (ManagedDiskDataSource) basic_withUltraSSD(data acceptance.TestData, name string, resourceGroupName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "%s"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "UltraSSD_LRS"
  create_option        = "Empty"
  disk_size_gb         = "4"
  disk_iops_read_write = "101"
  disk_mbps_read_write = "10"
  zones                = ["2"]

  tags = {
    environment = "acctest"
  }
}

data "azurerm_managed_disk" "test" {
  name                = azurerm_managed_disk.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, resourceGroupName, data.Locations.Primary, name)
}
