package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMDataLakeStore_payasyougo(t *testing.T) {
	dataSourceName := "data.azurerm_data_lake_store.test"
	rInt := acctest.RandIntRange(1, 999999)
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataLakeStore_payasyougo(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(dataSourceName),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMDataLakeStore_monthlycommitment(t *testing.T) {
	dataSourceName := "data.azurerm_data_lake_store.test"
	rInt := acctest.RandIntRange(1, 999999)
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataLakeStore_monthlycommitment(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "tier", "Commitment_1TB"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.hello", "world"),
				),
			},
		},
	})
}

func testAccDataSourceDataLakeStore_payasyougo(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name = "acctestRG_%d"
	location = "%s"
}

resource "azurerm_data_lake_store" "test" {
	name = "acctest%d"
	location = "%s"
	resource_group_name = "${azurerm_resource_group.test.name}"
	tags {
		hello = "world"
	}
}

data "azurerm_data_lake_store" "test" {
	name                = "${azurerm_data_lake_store.test.name}"
	resource_group_name = "${azurerm_data_lake_store.test.resource_group_name}"
}
`, rInt, location, rInt, location)
}

func testAccDataSourceDataLakeStore_monthlycommitment(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name = "acctestRG_%d"
	location = "%s"
}

resource "azurerm_data_lake_store" "test" {
	name = "acctest%d"
	location = "%s"
	tier = "Commitment_1TB"
	resource_group_name = "${azurerm_resource_group.test.name}"
	tags {
		hello = "world"
	}
}

data "azurerm_data_lake_store" "test" {
	name                = "${azurerm_data_lake_store.test.name}"
	resource_group_name = "${azurerm_data_lake_store.test.resource_group_name}"
}
`, rInt, location, rInt, location)
}
