package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMDiskAccess_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_disk_access", "test")

	name := fmt.Sprintf("acctestdiskaccess-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMDiskAccessBasic(data, name, resourceGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", name),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "acctest"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMDiskAccessBasic(data acceptance.TestData, name string, resourceGroupName string) string {
	return fmt.Sprintf(`
	provider "azurerm" {
		features{}
	}

	resource "azurerm_resource_group" "test" {
		name 	 = "%s"
		location = "%s"
	}

	resource "azurerm_disk_access"  "test" {
		name			    = "%s"
		location 			= azurerm_resource_group.test.location
		resource_group_name = azurerm_resource_group.test.name

		tags = {
			environment = "acctest"
		}
	}

	data "azurerm_disk_access" "test" {
		name 				= azurerm_disk_access.test.name
		resource_group_name = azurerm_resource_group.test.name
	}
	`, resourceGroupName, data.Locations.Primary, name)
}
