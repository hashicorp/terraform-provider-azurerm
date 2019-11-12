package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMResources_ByName(t *testing.T) {
	dataSourceName := "data.azurerm_resources.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMStorageAccount_basic(ri, rs, location),
			},
			{
				Config: testAccDataSourceAzureRMResources_ByName(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "resources.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMResources_ByResourceGroup(t *testing.T) {
	dataSourceName := "data.azurerm_resources.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMStorageAccount_basic(ri, rs, location),
			},
			{
				Config: testAccDataSourceAzureRMResources_ByResourceGroup(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "resources.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMResources_ByResourceType(t *testing.T) {
	dataSourceName := "data.azurerm_resources.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMStorageAccount_basic(ri, rs, location),
			},
			{
				Config: testAccDataSourceAzureRMResources_ByResourceType(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "resources.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMResources_FilteredByTags(t *testing.T) {
	dataSourceName := "data.azurerm_resources.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMStorageAccount_basic(ri, rs, location),
			},
			{
				Config: testAccDataSourceAzureRMResources_FilteredByTags(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "resources.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMResources_ByName(rInt int, rString string, location string) string {
	r := testAccDataSourceAzureRMStorageAccount_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_resources" "test" {
  name = "${azurerm_storage_account.test.name}"
}
`, r)
}

func testAccDataSourceAzureRMResources_ByResourceGroup(rInt int, rString string, location string) string {
	r := testAccDataSourceAzureRMStorageAccount_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_resources" "test" {
  resource_group_name = "${azurerm_storage_account.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMResources_ByResourceType(rInt int, rString string, location string) string {
	r := testAccDataSourceAzureRMStorageAccount_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_resources" "test" {
  resource_group_name = "${azurerm_storage_account.test.resource_group_name}"
  type                = "Microsoft.Storage/storageAccounts"
}
`, r)
}

func testAccDataSourceAzureRMResources_FilteredByTags(rInt int, rString string, location string) string {
	r := testAccDataSourceAzureRMStorageAccount_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_resources" "test" {
  name                = "${azurerm_storage_account.test.name}"
  resource_group_name = "${azurerm_storage_account.test.resource_group_name}"

  required_tags = {
    environment = "production"
  }
}
`, r)
}
