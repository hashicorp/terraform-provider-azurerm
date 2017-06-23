package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMManagedDisk_basic(t *testing.T) {
	ri := acctest.RandInt()

	name := fmt.Sprintf("acctestmanageddisk-%d", ri)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", ri)

	config := testAccDatSourceAzureRMManagedDiskBasic(name, resourceGroupName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.azurerm_managed_disk.test", "name", name),
					resource.TestCheckResourceAttr("data.azurerm_managed_disk.test", "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr("data.azurerm_managed_disk.test", "storage_account_type", "Premium_LRS"),
					resource.TestCheckResourceAttr("data.azurerm_managed_disk.test", "disk_size_gb", "10"),
					resource.TestCheckResourceAttr("data.azurerm_managed_disk.test", "tags.%", "1"),
					resource.TestCheckResourceAttr("data.azurerm_managed_disk.test", "tags.environment", "acctest"),
				),
			},
		},
	})
}

func testAccDatSourceAzureRMManagedDiskBasic(name string, resourceGroupName string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "%s"
    location = "West US"
}

resource "azurerm_managed_disk" "test" {
    name = "%s"
    location = "West US"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_type = "Premium_LRS"
    create_option = "Empty"
    disk_size_gb = "10"

    tags {
        environment = "acctest"
    }
}

data "azurerm_managed_disk" "test" {
    name = "${azurerm_managed_disk.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}
`, resourceGroupName, name)
}
