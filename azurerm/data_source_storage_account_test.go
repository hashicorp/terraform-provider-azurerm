package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMStorageAccount_basic(t *testing.T) {
	dataSourceName := "data.azurerm_storage_account.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := acceptance.Location()
	preConfig := testAccDataSourceAzureRMStorageAccount_basic(ri, rs, location)
	config := testAccDataSourceAzureRMStorageAccount_basicWithDataSource(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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

func TestAccDataSourceAzureRMStorageAccount_withWriteLock(t *testing.T) {
	dataSourceName := "data.azurerm_storage_account.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMStorageAccount_basicWriteLock(ri, rs, location),
			},
			{
				Config: testAccDataSourceAzureRMStorageAccount_basicWriteLockWithDataSource(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "account_tier", "Standard"),
					resource.TestCheckResourceAttr(dataSourceName, "account_replication_type", "LRS"),
					resource.TestCheckResourceAttr(dataSourceName, "primary_connection_string", ""),
					resource.TestCheckResourceAttr(dataSourceName, "secondary_connection_string", ""),
					resource.TestCheckResourceAttr(dataSourceName, "primary_blob_connection_string", ""),
					resource.TestCheckResourceAttr(dataSourceName, "secondary_blob_connection_string", ""),
					resource.TestCheckResourceAttr(dataSourceName, "primary_access_key", ""),
					resource.TestCheckResourceAttr(dataSourceName, "secondary_access_key", ""),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMStorageAccount_basic(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
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

func testAccDataSourceAzureRMStorageAccount_basicWriteLock(rInt int, rString string, location string) string {
	template := testAccDataSourceAzureRMStorageAccount_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = "${azurerm_storage_account.test.id}"
  lock_level = "ReadOnly"
}
`, template, rInt)
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

func testAccDataSourceAzureRMStorageAccount_basicWriteLockWithDataSource(rInt int, rString string, location string) string {
	config := testAccDataSourceAzureRMStorageAccount_basicWriteLock(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_storage_account" "test" {
  name                = "${azurerm_storage_account.test.name}"
  resource_group_name = "${azurerm_storage_account.test.resource_group_name}"
}
`, config)
}
