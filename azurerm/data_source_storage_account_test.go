package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMStorageAccount_basic(t *testing.T) {
	dataSourceName := "data.azurerm_storage_account.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccDataSourceAzureRMStorageAccount_basic(ri, rs, location)
	config := testAccDataSourceAzureRMStorageAccount_basicWithDataSource(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "account_tier", "Standard"),
					resource.TestCheckResourceAttr(dataSourceName, "account_replication_type", "LRS"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.environment", "production"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMStorageAccount_basic(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestsa-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsads%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccDataSourceAzureRMStorageAccount_basicWithDataSource(rInt int, rString string, location string) string {
	config := testAccDataSourceAzureRMStorageAccount_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_storage_account" "test" {
  name                = "${azurerm_storage_account.test.name}"
  resource_group_name = "${azurerm_storage_account.test.resource_group_name}"
}
`, config)
}
