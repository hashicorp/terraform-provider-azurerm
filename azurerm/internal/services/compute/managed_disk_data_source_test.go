package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMManagedDisk_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_disk", "test")

	name := fmt.Sprintf("acctestmanageddisk-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMManagedDiskBasic(data, name, resourceGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", name),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_account_type", "Premium_LRS"),
					resource.TestCheckResourceAttr(data.ResourceName, "disk_size_gb", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "acctest"),
					resource.TestCheckResourceAttr(data.ResourceName, "zones.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "zones.0", "2"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMManagedDisk_basic_withUltraSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_disk", "test")

	name := fmt.Sprintf("acctestmanageddisk-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMManagedDisk_basic_withUltraSSD(data, name, resourceGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "disk_iops_read_write", "101"),
					resource.TestCheckResourceAttr(data.ResourceName, "disk_mbps_read_write", "10"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMManagedDiskBasic(data acceptance.TestData, name string, resourceGroupName string) string {
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

func testAccDataSourceAzureRMManagedDisk_basic_withUltraSSD(data acceptance.TestData, name string, resourceGroupName string) string {
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
