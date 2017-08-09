package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMManagedDisk_basic(t *testing.T) {
	dataSourceName := "data.azurerm_managed_disk.test"
	ri := acctest.RandInt()

	name := fmt.Sprintf("acctestmanageddisk-%d", ri)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", ri)

	config := testAccDatSourceAzureRMManagedDiskBasic(name, resourceGroupName, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dataSourceName, "storage_account_type", "Premium_LRS"),
					resource.TestCheckResourceAttr(dataSourceName, "disk_size_gb", "10"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.environment", "acctest"),
				),
			},
		},
	})
}

func testAccDatSourceAzureRMManagedDiskBasic(name string, resourceGroupName string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "%s"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = "10"

  tags {
    environment = "acctest"
  }
}

data "azurerm_managed_disk" "test" {
  name                = "${azurerm_managed_disk.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, resourceGroupName, location, name)
}
