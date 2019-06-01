package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMBatchAccount_basic(t *testing.T) {
	dataSourceName := "data.azurerm_batch_account.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	config := testAccDataSourceAzureRMBatchAccount_basic(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("testaccbatch%s", rs)),
					resource.TestCheckResourceAttr(dataSourceName, "location", azureRMNormalizeLocation(location)),
					resource.TestCheckResourceAttr(dataSourceName, "pool_allocation_mode", "BatchService"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMBatchAccount_complete(t *testing.T) {
	dataSourceName := "data.azurerm_batch_account.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	config := testAccDataSourceAzureRMBatchAccount_complete(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("testaccbatch%s", rs)),
					resource.TestCheckResourceAttr(dataSourceName, "location", azureRMNormalizeLocation(location)),
					resource.TestCheckResourceAttr(dataSourceName, "pool_allocation_mode", "BatchService"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.env", "test"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMBatchAccount_basic(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batch"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  pool_allocation_mode = "BatchService"
}

data "azurerm_batch_account" "test" {
  name                = "${azurerm_batch_account.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rString)
}

func testAccDataSourceAzureRMBatchAccount_complete(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batch"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "testaccsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  pool_allocation_mode = "BatchService"
  storage_account_id   = "${azurerm_storage_account.test.id}"

  tags = {
    env = "test"
  }
}

data "azurerm_batch_account" "test" {
  name                = "${azurerm_batch_account.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rString, rString)
}
